package dao

import "fmt"

// 根据指定的文档标题获取文档编号
// env 存储着应用程序运行环境的结构体
// title 文档标题
// 返回文档编号
func DBGetDocumentID(title string) int {
	id, _ := GetDocumentId(title)
	return id
}

// 将文档添加到 documents 表中
// env 存储这应用程序运行环境的结构体
// title 文档标题
// body 文档正文
func DBAddDocument(title, body string) {
	id := DBGetDocumentID(title)
	fmt.Println("=== title: ", title)
	if id != 0 {
		UpdateDocument(id, body)
	} else {
		InsertDocument(title, body)
	}
}
