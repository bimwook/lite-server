package session

//Item 项
type Item struct {
	Rowid    string
	Name     string
	Created  string
	Modified string
}

//ISession 会话接口
type ISession interface {
	CheckIn(name string, secret string) (string, bool)
	Check(name string, token string) bool
	CheckOut(name string, token string) bool
}
