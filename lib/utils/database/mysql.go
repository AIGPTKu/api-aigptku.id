package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	Host            string
	Port            int  
	Username        string
	Password        string 
	DatabaseName    string 
	MaxIdleConns    int  
	MaxOpenConns    int  
	ConnMaxLifetime int 
	ConnMaxIdleTime int  
}

func NewMysqlConn(cfg MysqlConfig) *sql.DB {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Minute)

	return db
}