package woo

import (
	"database/sql"
	"fmt"
	"muen"
	"os"

	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

const maindb = "./dbase/main.db"

//ResetMain 初始化
func ResetMain() {
	os.MkdirAll("./dbase", os.ModePerm)
	_, err := os.Stat(maindb)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", maindb)
		uuid := muen.NewKey()
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
		_, e := db.Exec(cmd, uuid, muen.Now(), uuid)
		if e != nil {
			fmt.Println(e)
		}
		defer db.Close()
	}
}

//GetServerSerial 获取服务器标识
func GetServerSerial() string {
	ResetMain()
	db, _ := sql.Open("sqlite3", maindb)
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
	defer db.Close()
	return uuid
}
