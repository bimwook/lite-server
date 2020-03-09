package woo

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"muen"
	"os"

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
	cmd := `
	  CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
	  CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [hash], [remark], [data], [created]);
	  INSERT INTO [meta] ([rowid], [name], [content]) VALUES('server.uuid', 'SYSTEM', ?);
	  INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.uuid', 'SYSTEM', ?);
	  INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.created', 'SYSTEM', ?);
	`
	_, err := db.Exec(cmd, GetServerSerial(), muen.NewKey(), muen.Now())
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//Sha256 加密
func Sha256(data string) string {
	hm := sha256.New()
	hm.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}

//CenterSave 保存
func CenterSave(name, dbase, remark, data string) string {
	rowid := muen.NewKey()
	fn, exists := getCache(name, dbase)
	if !exists {
		resetCache(name, dbase)
	}
	db, _ := sql.Open("sqlite3", fn)
	cmd := `INSERT INTO [main] ([rowid], [hash], [remark], [data], [created]) VALUES(?,?,?,?,?);`
	_, err := db.Exec(cmd, rowid, Sha256(data), remark, data, muen.Now())
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return rowid
}

//CenterHash 哈希
func CenterHash(name, dbase, data string) string {
	fn, exists := getCache(name, dbase)
	cnt := "0"
	if exists {
		hash := Sha256(data)
		db, _ := sql.Open("sqlite3", fn)
		cmd := `SELECT COUNT(*) AS [cnt] FROM [main] WHERE [hash]=?;`
		rows, err := db.Query(cmd, hash)
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
		defer rows.Close()
		defer db.Close()
	}
	return cnt
}

//CenterLoad 加载
func CenterLoad(name, dbase, rowid string) (string, string) {
	fn, exists := getCache(name, dbase)
	if exists {
		remark := ""
		data := ""
		db, _ := sql.Open("sqlite3", fn)
		cmd := `SELECT [remark],[data] FROM [main] WHERE [rowid]=?;`
		rows, err := db.Query(cmd, rowid)
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
		defer rows.Close()
		defer db.Close()
		return remark, data
	}
	return "BAD-ROWID", ""
}
