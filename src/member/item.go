package member

//Item 项
type Item struct {
	Rowid   string
	Name    string
	Secret  string
	Level   int32
	Created string
}

//IMember 接口
type IMember interface {
	New(item Item) (string, bool)
	Renew(name string, secret string) (string, bool)
	Remove(name string, secret string) (string, bool)
	Check(name string, secret string) bool
}
