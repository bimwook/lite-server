package mail

import (
	"database/sql"
	"fmt"
	"os"

	"../../woo"

	_ "github.com/mattn/go-sqlite3" //Sqlite3
)

const mailcache = "./dbase/mail.db"

//ResetMailCache 初始化
func ResetMailCache() bool {
	_, err := os.Stat(mailcache)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", mailcache)
		defer db.Close()
		cmd := `
			CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [data]);
			CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [module], [sender], [receiver], [data], [created], [status]);
			INSERT INTO [meta] ([rowid], [name], [data]) VALUES('db.uuid', 'SYSTEM', ?);
			INSERT INTO [meta] ([rowid], [name], [data]) VALUES('db.created', 'SYSTEM', ?);
	  `
		_, e := db.Exec(cmd, woo.NewSerial(), woo.Now())
		if e != nil {
			fmt.Println(e)
		}
		return true
	}
	return false
}

//Save 保存
func Save(rowid string, module string, sender string, receiver string, data string) string {
	db, _ := sql.Open("sqlite3", mailcache)
	defer db.Close()
	cmd := `
		INSERT INTO [main] ([rowid], [module], [sender], [receiver], [data], [created], [status]) VALUES(?,?,?,?,?,?,?);
	`
	_, err := db.Exec(cmd, rowid, module, sender, receiver, data, woo.Now(), 0)
	if err != nil {
		fmt.Println(err)
	}
	return rowid
}

//Peek 窥
func Peek(module string, receiver string) string {
	ret := ""
	db, _ := sql.Open("sqlite3", mailcache)
	defer db.Close()
	cmd := `SELECT [sender], [data], [created] FROM [main] WHERE [module]=? AND [receiver]=? ORDER BY [created] ASC LIMIT 1;`
	rows, err := db.Query(cmd, module, receiver)
	defer rows.Close()
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
	return ret
}

//Receive 接收
func Receive(module string, receiver string) (string, string) {
	rowid := ""
	ret := ""
	db, _ := sql.Open("sqlite3", mailcache)
	defer db.Close()
	cmd := `SELECT [rowid], [sender], [data], [created] FROM [main] WHERE [module]=? AND [receiver]=? ORDER BY [created] ASC LIMIT 1;`
	rows, err := db.Query(cmd, module, receiver)
	defer rows.Close()
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
	return rowid, ret
}

//Remove 删除
func Remove(rowid string) bool {
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
