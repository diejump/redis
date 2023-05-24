package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type User struct {
	Username string
	Account  string
	Password string
}

var db *sql.DB

var Rds redis.Conn

func RedisPOllInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		Wait:        true,
		IdleTimeout: time.Duration(1) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			println("连接成功")
			redis.DialDatabase(0)
			return c, err
		},
	}
}

func RedisClose() {
	Rds.Close()
}

func InitDB() {
	var err error

	dsn := "root:123@tcp(127.0.0.1:3306)/school?charset=utf8mb4&parseTime=True&loc=Local"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB connect success")
	return
}

func AddUser(Account, Username, Password string) {
	sqlstr := "insert into user (username,account,password) values (?,?,?)"
	_, err := db.Exec(sqlstr, Username, Account, Password)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	log.Println("insert success")
}

func SelectUser(Account string) bool {
	sqlStr := "select password from user where account =?"
	var password string
	db.QueryRow(sqlStr, Account).Scan(&password)
	if password != "" {
		return true
	}
	return false
}

func SelectPasswordFromAccount(Account string) string {
	sqlstr := "select password from user where account=?"
	var password string
	db.QueryRow(sqlstr, Account).Scan(&password)
	return password
}

func Username() [20]string {
	sqlstr := "select username from user"
	var username [20]string
	rows, err := db.Query(sqlstr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(&username[i])
		i++
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
	}
	return username
}
