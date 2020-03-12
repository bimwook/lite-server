package woo

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

//数据库地址
const maindb = "./dbase/main.db"

//IsDone 是否已接收到退出信号
var IsDone = false

//Now 当前日期及时间
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//NewSerial 返回一个随机字符串
func NewSerial() string {
	ret := time.Now().Format("20060102")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		ret = ret + strconv.Itoa(r.Intn(10))
	}
	return ret
}

//ResetMain 初始化
func ResetMain() {
	os.MkdirAll("./dbase", os.ModePerm)
	_, err := os.Stat(maindb)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", maindb)
		defer db.Close()
		uuid := NewSerial()
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
		_, e := db.Exec(cmd, uuid, Now(), uuid)
		if e != nil {
			fmt.Println(e)
		}
	}
}

//GetServerSerial 获取服务器标识
func GetServerSerial() string {
	ResetMain()
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
	ResetMain()
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
	db, _ := sql.Open("sqlite3", maindb)
	defer db.Close()
	cmd := `REPLACE INTO [main] ([rowid], [content]) VALUES(?, ?);`
	_, e := db.Exec(cmd, name, value)
	if e != nil {
		fmt.Println(e)
	}
}
