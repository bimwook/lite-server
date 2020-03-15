package member

import (
	"database/sql"
	"os"

	"../server"
	"../woo"
)

type oMember struct {
	maindb string
}

//Start 初始化
func (o *oMember) Start() {
	os.MkdirAll("./dbase", os.ModePerm)
	_, err := os.Stat(o.maindb)
	if !((err == nil) || os.IsExist(err)) {
		db, _ := sql.Open("sqlite3", o.maindb)
		defer db.Close()
		cmd := `
			CREATE TABLE IF NOT EXISTS [meta] ([rowid] PRIMARY KEY, [name], [content]);
			CREATE TABLE IF NOT EXISTS [main] ([rowid] PRIMARY KEY, [name], [secret], [level], [created]);
		`
		db.Exec(cmd)
		cmd = `
			INSERT INTO [meta] ([rowid], [name], [content]) VALUES('server.uuid', 'SYSTEM', ?);
			INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.uuid', 'SYSTEM', ?);
			INSERT INTO [meta] ([rowid], [name], [content]) VALUES('db.created', 'SYSTEM', ?);
		`
		db.Exec(cmd, server.GetServerSerial(), woo.NewSerial(), woo.Now())
		cmd = `
			INSERT INTO [main] ([rowid], [name], [secret], [level], [created]) VALUES(?,?,?,?,?);
		`
		rowid := woo.NewSerial()
		secret := woo.Sha256(rowid + ":" + "110629")
		db.Exec(cmd, rowid, "root", secret, 999, woo.Now())
	}
}

//New 新成员
func (o *oMember) New(item Item) (string, bool) {
	db, error := sql.Open("sqlite3", o.maindb)
	defer db.Close()
	if error != nil {
		return "ERROR", false
	}
	cmd := `SELECT [name] FROM [main] WHERE [name]=?;`
	rows, e := db.Query(cmd, item.Name)
	defer rows.Close()
	if e != nil {
		return "ERROR", false
	}
	if rows.Next() {
		return "Exsits", false
	}
	cmd = `INSERT INTO [main] ([rowid], [name], [secret], [level], [created]) VALUES(?,?,?,?,?);`
	rowid := woo.NewSerial()
	secret := woo.Sha256(rowid + ":" + item.Secret)
	_, err := db.Exec(cmd, rowid, item.Name, secret, item.Level, woo.Now())
	if err != nil {
		return "ERROR", false
	}
	return "OK", true

}

//Renew 重置密钥
func (o *oMember) Renew(name string, secret string) (string, bool) {
	db, error := sql.Open("sqlite3", o.maindb)
	defer db.Close()
	if error != nil {
		return "Error", false
	}
	cmd := `SELECT [rowid],[secret] FROM [main] WHERE [name]=?;`
	var item Item
	rows, e := db.Query(cmd, name)
	if e == nil {
		if rows.Next() {
			rows.Scan(&item.Rowid, &item.Secret)
		} else {
			item.Rowid = ""
		}
		rows.Close()
	}
	if item.Rowid == "" {
		return "Bad-Member", false
	}
	hash := woo.Sha256(item.Rowid + ":" + secret)
	if hash != item.Secret {
		return "Access-Denied", false
	}
	item.Secret = woo.NewSerial()
	cmd = `UPDATE [main] SET [secret]=? WHERE [rowid]=?;`
	_, err := db.Exec(cmd, woo.Sha256(item.Rowid+":"+item.Secret), item.Rowid)
	if err != nil {
		item.Secret = "Error"
	}
	return item.Secret, err == nil
}

//Remove 删除成员
func (o *oMember) Remove(name string, secret string) (string, bool) {
	db, error := sql.Open("sqlite3", o.maindb)
	defer db.Close()
	if error != nil {
		return "Error", false
	}
	cmd := `SELECT [rowid],[secret] FROM [main] WHERE [name]=?;`
	rows, e := db.Query(cmd, name)
	defer rows.Close()
	if e != nil {
		return "Error", false
	}

	var item Item
	if rows.Next() {
		rows.Scan(&item.Rowid, &item.Secret)
	} else {
		return "Bad-Member", false
	}

	hash := woo.Sha256(item.Rowid + ":" + secret)
	if hash != item.Secret {
		return "Access-Denied", false
	}
	cmd = `DELETE FROM [main] WHERE [rowid]=?;`
	_, err := db.Exec(cmd, item.Rowid)
	return "", err == nil
}

//Check 重置密钥
func (o *oMember) Check(name string, secret string) bool {
	db, error := sql.Open("sqlite3", o.maindb)
	defer db.Close()
	if error != nil {
		return false
	}
	cmd := `SELECT [rowid],[secret] FROM [main] WHERE [name]=?;`
	var item Item
	rows, e := db.Query(cmd, name)
	if e == nil {
		if rows.Next() {
			rows.Scan(&item.Rowid, &item.Secret)
		} else {
			item.Rowid = ""
		}
		rows.Close()
	}
	if item.Rowid == "" {
		return false
	}
	hash := woo.Sha256(item.Rowid + ":" + secret)
	return hash == item.Secret
}

//GetActions 获取接口
func (o *oMember) GetActions() IMember {
	return o
}
