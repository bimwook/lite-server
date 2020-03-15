package center

import (
	"database/sql"
	"fmt"
	"os"

	"../../server"
	"../../woo"

	_ "github.com/mattn/go-sqlite3" // Sqlite3
)

type oCenter struct {
}

//GetCache 获取缓存文件
func (o *oCenter) GetCache(name, dbase string) (string, bool) {
	ret := "./dbase/cache/" + name + "/" + dbase + ".db"
	_, err := os.Stat(ret)
	return ret, err == nil || os.IsExist(err)
}

//ResetCache 初始化缓存文件
func (o *oCenter) ResetCache(name, dbase string) bool {
	dir := "./dbase/cache/" + name
	fn := dir + "/" + dbase + ".db"
	os.MkdirAll(dir, os.ModePerm)
	db, error := sql.Open("sqlite3", fn)
	defer db.Close()
	if error != nil {
		fmt.Println(error)
		return false
	}
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
func (o *oCenter) Save(item *Item) bool {
	fn, exists := o.GetCache(item.Name, item.Dbase)
	if !exists {
		o.ResetCache(item.Name, item.Dbase)
	}
	db, error := sql.Open("sqlite3", fn)
	defer db.Close()
	if error != nil {
		return false
	}
	cmd := `INSERT INTO [main] ([rowid], [hash], [remark], [data], [created]) VALUES(?,?,?,?,?);`
	_, err := db.Exec(cmd, item.Rowid, woo.Sha256(item.Data), item.Remark, item.Data, woo.Now())
	return err != nil
}

//Hash 哈希
func (o *oCenter) Hash(name, dbase, data string) string {
	fn, exists := o.GetCache(name, dbase)
	if !exists {
		return "-1"
	}
	hash := woo.Sha256(data)
	db, error := sql.Open("sqlite3", fn)
	defer db.Close()
	if error != nil {
		fmt.Println(error)
		return "-1"
	}
	cnt := "0"
	cmd := `SELECT COUNT(*) AS [cnt] FROM [main] WHERE [hash]=?;`
	rows, err := db.Query(cmd, hash)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	if rows.Next() {
		err := rows.Scan(&cnt)
		if err != nil {
			fmt.Println(err)
			cnt = "-1"
		}
	}
	return cnt
}

//Load 加载
func (o *oCenter) Load(name, dbase, rowid string) Item {
	var item Item
	fn, exists := o.GetCache(name, dbase)
	if !exists {
		item.Rowid = ""
		return item
	}

	db, error := sql.Open("sqlite3", fn)
	if error == nil {
		defer db.Close()
		cmd := `SELECT [remark],[data] FROM [main] WHERE [rowid]=?;`
		rows, err := db.Query(cmd, rowid)
		defer rows.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			if rows.Next() {
				err := rows.Scan(&item.Remark, &item.Data)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return item
}
