package logic

import (
	"encoding/xml"
	"fmt"
	"github.com/read-talk/wiser/dao"
	"github.com/read-talk/wiser/util"
	"io"
	"os"
)

func NewEnv(v int) *WiserEnv {
	return &WiserEnv{
		TokenLen:           NGram,            // 词元的长度。NGram中N的取值
		Compress:           CompressMethod{}, // 压缩倒排列表等数据的方法
		EnablePharseSearch: 0,                // 是否进行短语检索
		IIBuffer: &InvertedIndexHash{ // 用于更新倒排索引的缓冲区（Buffer）
			HashMap: make(map[int]*InvertedIndexValue),
			Items:   make([]*InvertedIndexValue, 0),
		},
		IIBufferCount:           0, // 用户更新倒排索引的缓冲区中的文档数
		IIBufferUpdateThreshold: v, // 缓冲区中文档数的阈值
		IndexedCount:            0, // 建立了索引的文档数
	}
}

// 将文档添加到数据库中，建立倒排索引
// title 文档标题，为 Nil 时将会清空缓冲区
// body 文档正文
func (env *WiserEnv) AddDocument(title, body string) error {
	if len(title) > 0 && len(body) > 0 {
		// 将文档标题和正文存储到数据库中
		dao.DBAddDocument(title, body)
		// 并获取该文档对应的文档编号
		documentID := dao.DBGetDocumentID(title)

		// 为文档创建倒排列表
		// 根据文档编号和文档内容更新存储在变量 env.IIBuffer 中的小倒排索引
		err := env.TextToPostingsLists(documentID, body)
		if err != nil {
			return err
		}
		env.IIBufferCount++ // 用户更新在缓冲区中已建立倒排索引的文档数
		env.IndexedCount++  // 建立了索引的文档数
		fmt.Printf("count: %d title: %s\n", env.IndexedCount, title)
	}

	// 存储在缓冲区中的文档数量达到了指定的阈值时，更新存储器上的倒排索引
	// 当 title 为空时，或者当已构建出小倒排索引的文档数量达到了阈值时，就合并索引
	// 另外，title 为空，还标志着所有的文档都已经处理完了。
	// 阈值设定得越小，内存的使用量也就越小，但会增加堆数据库的访问次数。
	// 反过来，阅知设定得越大，内存的使用量就越大，也减少了对数据库的访问次数。
	if len(env.IIBuffer.HashMap) > env.IIBufferUpdateThreshold || title == "" {
		fmt.Println("开始合并倒排索引")
		util.PrintTimeDiff()
		// 更新所有词元对应的倒排项，合并倒排索引，
		// 并将合并后的结果写入数据库(存储器)中。
		for _, p := range env.IIBuffer.HashMap {
			err := env.UpdatePostings(p)
			if err != nil {
				return err
			}
		}
		env.IIBuffer = nil
		env.IIBufferCount = 0
		util.PrintTimeDiff()
		fmt.Println("Index flushed合并倒排索引结束")
	}
	return nil
}

// 导入 wiki 数据
func (env *WiserEnv) LoadWikiDump(wikiDumpFile string,m int) error {
	xmlFile, err := os.Open(wikiDumpFile)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	var cnt int
	decoder := xml.NewDecoder(xmlFile)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				err = decoder.DecodeElement(&p, &se)
				err = env.AddDocument(p.Title, p.Text)
				if err != nil {
					fmt.Println("add document failed: ", err)
					return err
				}
				cnt++
				if cnt > m {
					break
				}
			}
		}
	}
	return nil
}
