package session

import (
	"database/sql"
	"os"

	"../server"
	"../woo"
	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

//数据库地址
const maindb = "./dbase/session.db"

var sessions, _ = sql.Open("sqlite3", ":memory:")

//Reset 初始化
func Reset() bool {
	os.MkdirAll("./dbase", os.ModePerm)
	_, err := os.Stat(maindb)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", maindb)
		defer db.Close()
		cmd := `
			CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
			CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [name], [secret], [token], [level], [created]);
		`
		db.Exec(cmd)
		cmd = `
			INSERT INTO [meta] ([rowid], [name], [content]) VALUES('server.uuid', 'SYSTEM', ?);
			INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.uuid', 'SYSTEM', ?);
			INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.created', 'SYSTEM', ?);
		`
		db.Exec(cmd, server.GetServerSerial(), woo.NewSerial(), woo.Now())
		cmd = `
			INSERT INTO [main] ([rowid], [name], [secret], [token], [level], [created]) VALUES(?,?,?,?,?,?);
		`
		db.Exec(cmd, "root", "root", "110629", woo.NewSerial(), 999, woo.Now())
		return true
	}
	cmd := `
		CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
		CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [name], [token], [created], [modified]);
	`
	sessions.Exec(cmd)
	return false
}

//CheckIn 开启新会话
func CheckIn() {

}

//CheckOut 关闭会话
func CheckOut() {

}
