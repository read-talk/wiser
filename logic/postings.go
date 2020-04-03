package logic

import (
	"encoding/json"
	"fmt"
	"github.com/read-talk/wiser/dao"
)

// 将内存上（小倒排索引中）的倒排列表与存储器上的倒排列表合并后存储到数据库中
// env 存储着应用程序运行环境的结构体
// p 含有倒排列表的倒排索引中的索引项
func (env *WiserEnv) UpdatePostings(p *InvertedIndexValue) error {

	// 从数据库中取出作为合并源的倒排列表
	oldPostings, oldPostingsLen, err := FetchPostings(p.TokenID)
	if err != nil {
		fmt.Printf("cannot fetch old postings list of token(%d) for update.", p.TokenID)
		return err
	}
	// 如果数据库中存在作为合并源的倒排列表
	if oldPostings != nil {
		// 就将该倒排列表和要合并进来的倒排列表合并在一起
		p.PostingsList = MergePostings(oldPostings, p.PostingsList)
		p.DocsCount += oldPostingsLen
	}
	// 将内存上的倒排列表转换成了字节序列
	buf, err := EncodePostings(p.PostingsList)
	if err != nil {
		return err
	}
	// 将转换后的字节序列存储到了数据库中
	dao.UpdatePostings(p.TokenID, p.DocsCount, buf)
	return nil
}

// 从数据库中获取关联到指定词元上的倒排列表
// env 存储着应用程序运行环境的结构体
// token id 词元编号
// 返回 postings 获取到的倒排列表
func FetchPostings(tokenID int) (*PostingsList, int, error) {
	docsCount, buf, err := dao.GetPostings(tokenID)
	if buf == "" {
		return nil, 0, nil
	}
	postings, err := DecodePostings(buf)
	if err != nil {
		return nil, 0, err
	}
	return postings, docsCount, nil
}

// 将倒排列表转换成字节序列
func EncodePostings(postings *PostingsList) (string, error) {
	buf, err := json.Marshal(postings)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// 对倒排列表进行还原或解码
func DecodePostings(buf string) (*PostingsList, error) {
	postings := &PostingsList{}

	err := json.Unmarshal([]byte(buf), postings)
	if err != nil {
		return nil, err
	}
	return postings, nil
}

// 获取将两个倒排列表合并后得到的倒排列表
func MergePostings(pa, pb *PostingsList) *PostingsList {
	var ret = &PostingsList{}
	var p = &PostingsList{}
	// 用pa和pb分别遍历base和to_be_added（参见函数merge_inverted_index）中的倒排列表中的元素，
	// 将二者连接成按文档编号升序排列的链表
	for pa != nil || pb != nil {
		var e = &PostingsList{}
		if pb == nil || (pa != nil && (pa.DocumentID <= pb.DocumentID)) {
			e = pa
			pa = pa.Next
		} else if pa == nil || pa.DocumentID >= pb.DocumentID {
			e = pb
			pb = pb.Next
		}
		e.Next = nil
		if ret == nil {
			ret = e
		} else {
			p.Next = e
		}
		p = e
	}
	return ret
}

// 合并两个倒排索引
// base 合并后其中的元素会增多的倒排索引(合并目标)
// to_be_added 合并后就被释放的倒排索引(合并源)
func MergeInvertedIndex(base, toBeAdded *InvertedIndexHash) {
	for _, p := range toBeAdded.HashMap {
		var t = &InvertedIndexValue{}
		delete(toBeAdded.HashMap, p.TokenID)
		t, ok := base.HashMap[p.TokenID]
		if ok {
			t.PostingsList = MergePostings(t.PostingsList, p.PostingsList)
			t.DocsCount += p.DocsCount
		} else {
			base.HashMap[p.TokenID] = p
		}
	}
}

// 打印倒排列表中的内容，用于调试
// postings 待打印的倒排列表
func dumpPostingsList(postings *PostingsList) {
	for e := postings; e != nil; e = e.Next {
		fmt.Printf(" doc_id: %d (", e.DocumentID)
		if e.Positions != nil {
			for _, p := range e.Positions {
				fmt.Printf(" (位置: %d ) ", p)
			}
		}
		fmt.Printf(") ")
	}
}

// 输出倒排索引的内容
func DumpInvertedIndex(c *InvertedIndexHash) {
	for _, it := range c.HashMap {
		if it.TokenID != 0 {
			token, _ := dao.GetToken(it.TokenID)
			fmt.Printf("TOKEN %d.%s(%d):\n", it.TokenID, token, it.DocsCount)
		} else {
			fmt.Println("TOKEN NONE:")
		}
		if it.PostingsList != nil {
			fmt.Printf("POSTINGS: [")
			dumpPostingsList(it.PostingsList)
			fmt.Printf("]\n")
		}
	}
}
