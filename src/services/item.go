package services

//MailItem 项
type MailItem struct {
	Rowid    string
	Module   string
	Sender   string
	Receiver string
	Data     string
}

//CenterItem 项
type CenterItem struct {
	Name   string
	Dbase  string
	Remark string
	Data   string
}
