package logic

import (
	"fmt"
	"github.com/read-talk/wiser/dao"
	"github.com/read-talk/wiser/util"
)

// 为构成文档内容的字符串建立倒排列表的集合(倒排文件)
// document id 文档编号。为0时表示要把查询的关键词作为处理对象
// text 输入的字符串
func (env *WiserEnv) TextToPostingsLists(documentId int, text string) error {
	// 分隔 N-gram 词元
	runeBody := []rune(text)
	start := 0
	for {
		// 每次从字符串中取出长度为 N-gram 的词元
		tokenLen, position := util.NgramNext(runeBody, &start, env.TokenLen)
		if tokenLen == 0 {
			break
		}
		if tokenLen < env.TokenLen {
			continue
		}
		// 将该词元添加到倒排列表中
		token := string(runeBody[position : position+env.TokenLen])
		err := env.TokenToPostingsList(documentId, token, start)
		if err != nil {
			return err
		}
	}
	// 当循环结束后，传入的 text 构成的倒排索引就构建好了。
	return nil
}

// 为传入的词元创建倒排列表
// document id 文档编号
// token 词元
// start 词元出现的位置
func (env *WiserEnv) TokenToPostingsList(id int, token string, start int) error {
	// 获取词元对应的编号
	tokenID, _ := DBGetTokenID(token, id)

	IIEntry, ok := env.IIBuffer[tokenID]
	if !ok {
		IIEntry = InvertedIndex{
			PostingsMap:   map[int][]int{},
			PostingsCount: 1,
		}
		env.IIBuffer[tokenID] = IIEntry
	}
	_, ok = IIEntry.PostingsMap[id]
	if !ok {
		IIEntry.PostingsMap[id] = []int{}
	}
	// 存储位置信息
	IIEntry.PostingsMap[id] = append(IIEntry.PostingsMap[id], start)
	IIEntry.PostingsCount++
	return nil
}

// 如果之前已将编号分配给了该词元，那么获取的正是这个编号
// 如果之前没有分配编号，则为该词元分配一个新的编号
func DBGetTokenID(token string, id int) (int, int) {
	if id != 0 {
		dao.StoreToken(token, "")
	}
	id, count, err := dao.GetTokenId(token)
	if err != nil {
		fmt.Println("为传入的词元创建倒排列表时出错, err: ", err)
		return 0, 0
	}
	return id, count
}
