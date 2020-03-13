package server

import (
	"database/sql"
	"fmt"
	"os"

	"../woo"
	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

//ServerName 服务器标识
const ServerName = "Lite-Server/1.1"

//数据库地址
const maindb = "./dbase/main.db"

//IsTerminated 是否已接收到退出信号
var IsTerminated = false

//Start 初始化
func Start() {
	reset()
}

//reset 初始化
func reset() {
	os.MkdirAll("./dbase", os.ModePerm)
	_, err := os.Stat(maindb)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", maindb)
		defer db.Close()
		uuid := woo.NewSerial()
		cmd := `
      CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [content]);
      CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
      INSERT INTO [meta] ([rowid], [name], [content]) VALUES('about.author', 'SYSTEM', '杨波 ( Bamboo Young )');
      INSERT INTO [meta] ([rowid], [name], [content]) VALUES('about.email', 'SYSTEM', 'lokme@foxmail.com');
      INSERT INTO [meta] ([rowid], [name], [content]) VALUES('about.website', 'SYSTEM', 'https://me.bimwook.com');
      INSERT INTO [meta] ([rowid], [name], [content]) VALUES('server.uuid', 'SYSTEM', ?);
      INSERT INTO [meta] ([rowid], [name], [content]) VALUES('server.created', 'SYSTEM', ?);
      INSERT INTO [main] ([rowid], [content]) VALUES('server.uuid', ?);
    `
		_, e := db.Exec(cmd, uuid, woo.Now(), uuid)
		if e != nil {
			fmt.Println(e)
		}
	}
}

//GetServerSerial 获取服务器标识
func GetServerSerial() string {
	reset()
	db, _ := sql.Open("sqlite3", maindb)
	defer db.Close()
	uuid := "BAD-KEY"
	cmd := `SELECT [content] FROM [main] WHERE [rowid]='server.uuid';`
	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			var content string
			e := rows.Scan(&content)
			if e == nil {
				uuid = content
			}
		}
	}
	return uuid
}

//GetValue 获取配置
func GetValue(name string) string {
	reset()
	db, _ := sql.Open("sqlite3", maindb)
	defer db.Close()
	v := ""
	cmd := `SELECT [content] FROM [main] WHERE [rowid]=?;`
	rows, err := db.Query(cmd, name)
	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			var content string
			e := rows.Scan(&content)
			if e == nil {
				v = content
			}
		}
	}
	return v
}

//SetValue 设置配置
func SetValue(name string, value string) {
	reset()
	db, _ := sql.Open("sqlite3", maindb)
	defer db.Close()
	cmd := `REPLACE INTO [main] ([rowid], [content]) VALUES(?, ?);`
	_, e := db.Exec(cmd, name, value)
	if e != nil {
		fmt.Println(e)
	}
}
