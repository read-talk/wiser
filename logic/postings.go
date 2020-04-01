package logic

import (
	"fmt"
	"github.com/read-talk/wiser/dao"
	"github.com/read-talk/wiser/util"
)

// 将内存上（小倒排索引中）的倒排列表与存储器上的倒排列表合并后存储到数据库中
// env 存储着应用程序运行环境的结构体
// p 含有倒排列表的倒排索引中的索引项
func (env *WiserEnv) UpdatePostingsAndFree() error {
	for tokenID, IIEntry := range env.IIBuffer {
		// 从数据库中取出作为合并源的倒排列表
		oldPostings, err := FetchPostings(tokenID)
		if err != nil {
			return err
		}
		// 如果数据库中存在作为合并源的倒排列表
		if oldPostings != nil {
			// 就将该倒排列表和要合并进来的倒排列表合并在一起
			IIEntry.PostingsMap = util.MergePostings(oldPostings, IIEntry.PostingsMap)
		}
		// 将内存上的倒排列表转换成了字节序列
		buf, err := util.EncodePostings(IIEntry.PostingsMap)
		if err != nil {
			return err
		}
		// 将转换后的字节序列存储到了数据库中
		dao.UpdatePostings(tokenID, len(IIEntry.PostingsMap), buf)
	}
	env.IIBuffer = map[int]InvertedIndex{}
	fmt.Println("Index flushed")
	return nil
}

// 从数据库中获取关联到指定词元上的倒排列表
// env 存储着应用程序运行环境的结构体
// token id 词元编号
// 返回 postings 获取到的倒排列表
func FetchPostings(tokenID int) (map[int][]int, error) {
	_, buf, err := dao.GetPostings(tokenID)
	if buf == "" {
		return nil, nil
	}

	postings, err := util.DecodePostings(buf)
	if err != nil {
		return nil, err
	}
	return postings, nil
}
