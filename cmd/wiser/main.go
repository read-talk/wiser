package main

import (
	"flag"
	"fmt"
	"github.com/read-talk/wiser/dao"
	"github.com/read-talk/wiser/logic"
)

var (
	x                              string
	q                              string
	m                              int
	DefaultIiBufferUpdateThreshold = 2048
)

func init() {
	flag.StringVar(&x, "x", "", "wikipedia dump xml path for indexing")
	flag.StringVar(&q, "q", "", "query for search")
	flag.IntVar(&m, "m", 10, "max count for indexing document")
}

func main() {
	flag.Parse()
	// 初始化全局环境
	env := logic.NewEnv(DefaultIiBufferUpdateThreshold)

	var err error
	// 加载wiki的词条数据
	if x != "" {
		fmt.Println("需要构建索引的文件: ", x)
		// 加载 wiki 的词条数据
		err = env.LoadWikiDump(x, m)
		if err != nil {
			fmt.Println("failed to load wiki, err: ", err)
			return
		}
	}

	// 进行检索
	if q != "" {
		fmt.Println("查询: ", q)
		env.IndexedCount, err = dao.GetDocumentCount()
		if err != nil {
			fmt.Println("failed to query, err: ", err)
			return
		}
		env.Search(q)
	}
}
