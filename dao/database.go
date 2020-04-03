package dao

// 根据指定的文档标题获取文档编号
// title 文档标题
// 返回文档编号
func DBGetDocumentID(title string) int {
	id, _ := GetDocumentId(title)
	return id
}

// 将文档添加到 documents 表中
// title 文档标题
// body 文档正文
func DBAddDocument(title, body string) {
	id := DBGetDocumentID(title)
	if id != 0 {
		UpdateDocument(id, body)
	} else {
		InsertDocument(title, body)
	}
}
