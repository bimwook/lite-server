package session

import (
	"database/sql"
	"fmt"
	"time"

	"../member"
	"../server"
	"../woo"
	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

//数据库地址

var sessions, _ = sql.Open("sqlite3", ":memory:")

type oSession struct {
}

//Start 初始化
func (s *oSession) Start() bool {
	cmd := `
		CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
		CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [name], [created], [modified]);
	`
	sessions.Exec(cmd)
	go func() {
		for {
			if server.IsTerminated {
				sessions.Close()
				break
			}
			time.Sleep(10)
		}
		fmt.Println("[Session] Closed.")
	}()

	return true
}

//CheckIn 开启新会话
func (s *oSession) CheckIn(name string, secret string) (string, bool) {
	ret := member.Actions.Check(name, secret)
	token := ""
	if ret {
		token = woo.NewSerial()
		cmd := `INSERT INTO [main] ([rowid], [name], [created], [modified]) VALUES(?,?,?,?);`
		sessions.Exec(cmd, token, name, woo.Now(), woo.Now())
	}
	return token, ret
}

// Check 检查会话状态
func (s *oSession) Check(name string, token string) bool {
	cmd := `SELECT [rowid] FROM [main] WHERE [rowid]=? AND [name]=?;`
	rowid := ""
	rows, e := sessions.Query(cmd, token, name)
	if e == nil {
		if rows.Next() {
			rows.Scan(&rowid)
		} else {
			rowid = ""
		}
		rows.Close()
	}
	if rowid == "" {
		return false
	}
	sessions.Exec("UPDATE [main] SET [modified]=? WHERE [rowid]=?;", woo.Now(), token)
	return true
}

//CheckOut 关闭会话
func (s *oSession) CheckOut(name string, token string) bool {
	cmd := `SELECT [rowid] FROM [main] WHERE [rowid]=? AND [name]=?;`
	rowid := ""
	rows, e := sessions.Query(cmd, token, name)
	if e == nil {
		if rows.Next() {
			rows.Scan(&rowid)
		} else {
			rowid = ""
		}
		rows.Close()
	}
	if rowid == "" {
		return false
	}
	sessions.Exec("DELETE FROM [main] WHERE [rowid]=?;", token)
	return true
}

//GetActions 获取接口
func (s *oSession) GetActions() ISession {
	return s
}
