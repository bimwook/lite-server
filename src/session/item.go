package session

//Item 项
type Item struct {
	Rowid  string
	Name   string
	Secret string
	Token  string
	Level  int32
}

//Session 会话接口
type Session interface {
	CheckIn(name string, secret string) (bool, string)
	CheckOut() bool
	New(item Item)
}
