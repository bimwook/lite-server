package mail

//Item 项
type Item struct {
	Rowid    string
	Module   string
	Sender   string
	Receiver string
	Data     string
}

//IMail 接口
type IMail interface {
	Save(item *Item) string
	Peek(module string, receiver string) string
	Receive(module string, receiver string) (string, string)
	Remove(rowid string) bool
}
