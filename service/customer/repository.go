package customer

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
    "time"
)


type Repository interface {
	AddCustomer(id string)(error)
	DeactiveCustomer(id string)(error)
	IncreaseUsedToken(id string,token int)(error)
    GetUserToken(id string)(int,int,error)
}

type DefatultRepository struct {
	DB *sql.DB
}

func (repo *DefatultRepository)Begin()(*sql.Tx, error){
	return repo.DB.Begin()
}

func (repo *DefatultRepository)AddCustomer(id string)(error){
    log.Printf("add customer:%s\n",id)
    now:=time.Now().Format("2006-01-02 15:04:05")
    //开启事务
	tx,err:= repo.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	_,err=tx.Exec("INSERT INTO gpt_customer (id,update_time,create_time,update_user,create_user) VALUES (?,?,?,'sys','sys') ON DUPLICATE KEY UPDATE subscribe =1,version=version+1,update_time=?,update_user='sys'", id,now,now,now)
    if err!=nil {
        tx.Rollback()
        return err
    }

    _,err=tx.Exec("INSERT INTO gpt_customer_token (id,update_time,create_time,update_user,create_user) VALUES (?,?,?,'sys','sys') ON DUPLICATE KEY UPDATE version=version+1,update_time=?,update_user='sys'", id,now,now,now)
    if err!=nil {
        tx.Rollback()
        return err
    }
    
    //提交事务
    err = tx.Commit(); 
    return err
}

func (repo *DefatultRepository)DeactiveCustomer(id string)(error){
    log.Printf("DeactiveCustomer customer:%s\n",id)
    now:=time.Now().Format("2006-01-02 15:04:05")
	_,err:=repo.DB.Exec("update gpt_customer set subscribe =0,version=version+1,update_user='sys',update_time=? where id=?", now,id)
    return err
}

func (repo *DefatultRepository)IncreaseUsedToken(id string,token int)(error){
    now:=time.Now().Format("2006-01-02 15:04:05")
	_,err:=repo.DB.Exec("update gpt_customer_token set  used=used+?,version=version+1,update_user='sys',update_time=?  where id=?",token,now,id)
    return err
}

func (repo *DefatultRepository)GetUserToken(id string)(int,int,error){
    row:= repo.DB.QueryRow("select total,used from gpt_customer_token where id=?",id)

    var total int
    var used int
    if err := row.Scan(&total,&used); err != nil {
        return 0,0,err
    }
    return total,used,nil
}

func (repo *DefatultRepository)GetUserAccessToken(id string)(string,string,error){
    row:= repo.DB.QueryRow("select weichat_access_token as accessToken,weichat_refresh_token as refreshToken from gpt_customer where id=?",id)

    var accessToken string
    var refreshToken string
    if err := row.Scan(&accessToken,&refreshToken); err != nil {
        return "","",err
    }
    return total,used,nil
}

func (repo *DefatultRepository)UpdateUserAccessToken(id string,accessToken string,refreshToken string)(error){
    now:=time.Now().Format("2006-01-02 15:04:05")
    _,err:=repo.DB.Exec("update gpt_customer set  weichat_access_token=?,weichat_refresh_token=?,version=version+1,update_user='sys',update_time=?  where id=?",accessToken,refreshToken,now,id)
    return err
}

func (repo *DefatultRepository)Connect(
	server,user,password,dbName string,
	connMaxLifetime,maxOpenConns,maxIdleConns int){ 
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
		
    repo.DB.SetConnMaxLifetime(time.Minute * time.Duration(connMaxLifetime))
	repo.DB.SetMaxOpenConns(maxOpenConns)
	repo.DB.SetMaxIdleConns(maxIdleConns)
    log.Println("connect to mysql server "+server)
}