package woo

import (
	"database/sql"
	"fmt"
	"muen"
	"os"

	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

const mailcache = "./dbase/mail.db"

//ResetMailCache 初始化
func ResetMailCache() bool {
	_, err := os.Stat(mailcache)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", mailcache)
		cmd := `
		CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [data]);
		CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [module], [sender], [receiver], [data], [created], [status]);
		INSERT INTO [meta] ([rowid], [name], [data]) VALUES('db.uuid', 'SYSTEM', ?);
		INSERT INTO [meta] ([rowid], [name], [data]) VALUES('db.created', 'SYSTEM', ?);
	  `
		_, e := db.Exec(cmd, muen.NewKey(), muen.Now())
		if e != nil {
			fmt.Println(e)
		}
		defer db.Close()
		return true
	}
	return false
}

//MailSave 保存
func MailSave(rowid string, module string, sender string, receiver string, data string) string {
	db, _ := sql.Open("sqlite3", mailcache)
	cmd := `
	  INSERT INTO [main] ([rowid], [module], [sender], [receiver], [data], [created], [status]) VALUES(?,?,?,?,?,?,?);
	`
	_, err := db.Exec(cmd, rowid, module, sender, receiver, data, muen.Now(), 0)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	return rowid
}

//MailPeek 窥
func MailPeek(module string, receiver string) string {
	ret := ""
	db, _ := sql.Open("sqlite3", mailcache)
	cmd := `SELECT [sender], [data], [created] FROM [main] WHERE [module]=? AND [receiver]=? ORDER BY [created] ASC LIMIT 1;`
	rows, err := db.Query(cmd, module, receiver)
	if err != nil {
		fmt.Println(err)
	} else {
		if rows.Next() {
			var sender string
			var data string
			var created string
			e := rows.Scan(&sender, &data, &created)
			if e == nil {
				ret = created + " | " + sender + "\r\n" + data
			}
		}
	}
	defer rows.Close()
	defer db.Close()
	return ret
}

//MailReceive 接收
func MailReceive(module string, receiver string) (string, string) {
	rowid := ""
	ret := ""
	db, _ := sql.Open("sqlite3", mailcache)
	cmd := `SELECT [rowid], [sender], [data], [created] FROM [main] WHERE [module]=? AND [receiver]=? ORDER BY [created] ASC LIMIT 1;`
	rows, err := db.Query(cmd, module, receiver)
	if err != nil {
		fmt.Println(err)
	} else {
		if rows.Next() {
			var sender string
			var data string
			var created string
			e := rows.Scan(&rowid, &sender, &data, &created)
			if e == nil {
				ret = created + " | " + sender + "\r\n" + data
			}
		}
	}
	defer rows.Close()
	defer db.Close()
	return rowid, ret
}

//MailRemove 删除
func MailRemove(rowid string) bool {
	db, _ := sql.Open("sqlite3", mailcache)
	cmd := `DELETE FROM [main] WHERE [rowid]=?;`
	_, err := db.Exec(cmd, rowid)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
