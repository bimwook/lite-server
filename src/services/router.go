package services

import (
	"./center"
	"./mail"
)

//Start 启动服务
func Start() {
	mail.Start()
	center.Start()
}
