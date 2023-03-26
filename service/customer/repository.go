package customer

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)


type Repository interface {
	AddCustomer(id string)(error)
	UpdateCustomer()(error)
	getFaultCountByType()(*FaultTypeCount,error)
	getFaultCountByStatus()(*FaultStatusCount,error)
	getFaultList()([]map[string]interface{},error)
	query(sql string)([]map[string]interface{},error)
	closeFault(diagReport string,remark string)
}

type DefatultRepository struct {
	DB *sql.DB
}

func (repo *DefatultRepository)query(sql string)([]map[string]interface{},error){
	rows, err := repo.DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows)
}

func (repo *DefatultRepository)getCarCount()(int,error){
	row := repo.DB.QueryRow("select count(*) as count from vehiclemanagement")
    var count int = 0
	if err := row.Scan(&count); err != nil {
        log.Println("getCarCount error")
		log.Println(err)
        return 0,nil
    }
	return count, nil
}

func (repo *DefatultRepository)toMap(rows *sql.Rows)([]map[string]interface{},error){
	cols,_:=rows.Columns()
	columns:=make([]interface{},len(cols))
	colPointers:=make([]interface{},len(cols))
	for i,_:=range columns {
		colPointers[i] = &columns[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err:= rows.Scan(colPointers...)
		if err != nil {
			log.Println(err)
			return nil,nil
		}
		row:=make(map[string]interface{})
		for i,colName :=range cols {
			val:=colPointers[i].(*interface{})
			switch (*val).(type) {
			case []byte:
				row[colName]=string((*val).([]byte))
			default:
				row[colName]=*val
			} 
		}
		list=append(list,row)
	}
	return list,nil
}

func (repo *DefatultRepository)getFaultList()([]map[string]interface{},error){
	rows, err := repo.DB.Query("select * from  diag_result order by status asc,time desc limit 0,500")
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows)
}

func (repo *DefatultRepository)getCarCountByProject()([]map[string]interface{},error){
	rows, err := repo.DB.Query("select ProjectNum, count(*) as count from vehiclemanagement group by ProjectNum order by ProjectNum")
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows) 
}

func (repo *DefatultRepository)getFaultCountByType()(*FaultTypeCount,error){
	var typeCount FaultTypeCount
	row:= repo.DB.QueryRow("select sum(eps) as epsCount,sum(ibs) as ibsCount,sum(esc) as escCount from diag_result")
	if err := row.Scan(&typeCount.EpsCount, &typeCount.IbsCount, &typeCount.EscCount); err != nil {
        log.Println("getFaultCountByType error")
		log.Println(err)
    } 
	return &typeCount, nil
}

func (repo *DefatultRepository)getFaultCountByStatus()(*FaultStatusCount,error){
	var statusCount FaultStatusCount
	row:= repo.DB.QueryRow("SELECT count(if(status=0,true,null)) as openCount,count(if(status=1,true,null)) as closedCount FROM diag_result")
	if err := row.Scan(&statusCount.OpenCount, &statusCount.ClosedCount); err != nil {
        log.Println("getFaultCountByType error")
		log.Println(err)
    } 
	return &statusCount, nil
}

func (repo *DefatultRepository)closeFault(diagReport string,remark string){
	repo.DB.Exec("update DiagResult set Status='1',Remark=? where DiagReport = ?", remark,diagReport)
}

func (repo *DefatultRepository)Connect(server string,user string,password string,dbName string){
	// Capture connection properties.
    cfg := mysql.Config{
        User:   user,
        Passwd: password,
        Net:    "tcp",
        Addr:   server,
        DBName: dbName,
		AllowNativePasswords:true,
    }
    // Get a database handle.
    var err error
    repo.DB, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := repo.DB.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    log.Println("connect to mysql server "+server)
}