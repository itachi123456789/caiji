package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var Db *sql.DB //全局数据库连接
const ApiKey = "66CAB6BAFAA0E405640B175814887B01920CB880836C6BD2B1FFF23136EE
var gMysql = make(chan string, 3000) //使用channel 一条一条入库 避免对数据库造成压力

func openMysql() {
	var err error
	Db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/xsweb?charset=utf8")
	if err != nil {
		log.Fatalf("连接数据库错误: %s\n", err)
	}
	//defer Db.Close()
	Db.SetMaxOpenConns(1)
	Db.SetMaxIdleConns(1)
	err = Db.Ping()
	if err != nil {
		log.Fatalln("Db.Ping():", err)
	}
}

func init() {
	openMysql()
}

func main() {
	go dataIn() //入库协程
	go http.HandleFunc("/dataIn", rec)
	http.ListenAndServe(":7274", nil)

}

func rec(w http.ResponseWriter, r *http.Request) {
	salt := r.FormValue("salt")
	code := r.FormValue("code")
	body := r.FormValue("body")

	// 验证权限
	m := md5.Sum([]byte(ApiKey + body + salt))
	mcode := salt + hex.EncodeToString(m[:])
	if code == mcode {
		gMysql <- body

	} else {
		w.Write([]byte("验证错误！"))
	}

}

func dataIn() {
	defer fmt.Println("dataIn QUIT")
	var err error
	for {
		select {
		case s := <-gMysql:
			_, err = Db.Query(s)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
