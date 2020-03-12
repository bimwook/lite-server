package center

import (
	"database/sql"
	"fmt"
	"os"

	"../../server"
	"../../woo"

	_ "github.com/mattn/go-sqlite3" // Sqlite3
)

//GetCache 获取缓存文件
func getCache(name, dbase string) (string, bool) {
	ret := "./dbase/cache/" + name + "/" + dbase + ".db"
	_, err := os.Stat(ret)
	return ret, err == nil || os.IsExist(err)
}

//ResetCache 初始化缓存文件
func resetCache(name, dbase string) bool {
	dir := "./dbase/cache/" + name
	fn := dir + "/" + dbase + ".db"
	os.MkdirAll(dir, os.ModePerm)
	db, _ := sql.Open("sqlite3", fn)
	defer db.Close()
	cmd := `
		CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
		CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [hash], [remark], [data], [created]);
		INSERT INTO [meta] ([rowid], [name], [content]) VALUES('server.uuid', 'SYSTEM', ?);
		INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.uuid', 'SYSTEM', ?);
		INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.created', 'SYSTEM', ?);
	`
	_, err := db.Exec(cmd, server.GetServerSerial(), woo.NewSerial(), woo.Now())
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//Save 保存
func Save(name, dbase, remark, data string) string {
	rowid := woo.NewSerial()
	fn, exists := getCache(name, dbase)
	if !exists {
		resetCache(name, dbase)
	}
	db, _ := sql.Open("sqlite3", fn)
	defer db.Close()
	cmd := `INSERT INTO [main] ([rowid], [hash], [remark], [data], [created]) VALUES(?,?,?,?,?);`
	_, err := db.Exec(cmd, rowid, woo.Sha256(data), remark, data, woo.Now())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return rowid
}

//Hash 哈希
func Hash(name, dbase, data string) string {
	fn, exists := getCache(name, dbase)
	cnt := "0"
	if exists {
		hash := woo.Sha256(data)
		db, _ := sql.Open("sqlite3", fn)
		defer db.Close()
		cmd := `SELECT COUNT(*) AS [cnt] FROM [main] WHERE [hash]=?;`
		rows, err := db.Query(cmd, hash)
		defer rows.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			if rows.Next() {
				err := rows.Scan(&cnt)
				if err != nil {
					cnt = "0"
				}
			}
		}

	}
	return cnt
}

//Load 加载
func Load(name, dbase, rowid string) (string, string) {
	fn, exists := getCache(name, dbase)
	if exists {
		remark := ""
		data := ""
		db, _ := sql.Open("sqlite3", fn)
		defer db.Close()
		cmd := `SELECT [remark],[data] FROM [main] WHERE [rowid]=?;`
		rows, err := db.Query(cmd, rowid)
		defer rows.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			if rows.Next() {
				err := rows.Scan(&remark, &data)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		return remark, data
	}
	return "BAD-ROWID", ""
}
