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
	var bufferPostings = &InvertedIndexHash{}
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

	if env.IIBuffer != nil {
		MergeInvertedIndex(env.IIBuffer, bufferPostings)
	} else {
		env.IIBuffer = bufferPostings
	}
	return nil
}

// 为传入的词元创建倒排列表
// document id 文档编号
// token 词元
// start 词元出现的位置
func (env *WiserEnv) TokenToPostingsList(id int, token string, start int) error {
	// 获取词元对应的编号
	tokenID, docsCount := DBGetTokenID(token, id)
	var pl = &PostingsList{}
	IIEntry, ok := env.IIBuffer.HashMap[tokenID]
	if ok {
		pl = IIEntry.PostingsList
		pl.PositionsCount++
	} else {
		IIEntry = &InvertedIndexValue{
			TokenID:       tokenID,   // 词元编号（Token ID）
			PostingsList:  nil,       // 指向包含该词元的倒排列表的指针
			DocsCount:     docsCount, // 出现过该词元的文档数
			PostingsCount: 0,         // 该词元在所有文档中的出现次数之和
		}
		env.IIBuffer.HashMap[tokenID] = IIEntry

		pl = &PostingsList{
			DocumentID:     id,
			Positions:      nil,
			PositionsCount: 1,
			Next:           nil,
		}
		IIEntry.PostingsList = pl
	}
	// 存储位置信息
	pl.Positions = append(pl.Positions, start)
	IIEntry.PostingsCount++
	return nil
}

// 如果之前已将编号分配给了该词元，那么获取的正是这个编号
// 如果之前没有分配编号，则为该词元分配一个新的编号
// id 传入的是文档id，如果不为空，就存储该词元
// 返回该词元的id和出现过指定词元的文档数
func DBGetTokenID(token string, id int) (int, int) {
	if id != 0 {
		dao.StoreToken(token, "")
	}
	tokenID, count, err := dao.GetTokenId(token)
	if err != nil {
		fmt.Println("为传入的词元创建倒排列表时出错, err: ", err)
		return 0, 0
	}
	return tokenID, count
}
