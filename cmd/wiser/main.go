package main

import (
	"flag"
	"fmt"
	"github.com/read-talk/wiser/logic"
	"github.com/read-talk/wiser/util"
)

var (
	x                              string
	DefaultIiBufferUpdateThreshold = 2048
)

func init() {
	flag.StringVar(&x, "x", "wiki.xml", "wiki dump file")
}

func main() {
	flag.Parse()
	fmt.Println("需要构建索引的文件: ", x)
	// 初始化全局环境
	env := logic.NewEnv(DefaultIiBufferUpdateThreshold)
	util.PrintTimeDiff()
	if x != "" {
		// 加载 wiki 的词条数据
		err := env.LoadWikiDump(x)
		if err != nil {
			fmt.Println("failed to load wiki, err: ", err)
			return
		}
	}
	util.PrintTimeDiff()
}
