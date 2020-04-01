package logic

import (
	"fmt"
	"github.com/read-talk/wiser/dao"
)

type Results struct {
	DocumentID int     // 检索出的文档编号
	Score      float64 // 检索得分
}

type SearchResults struct {
	Results []*Results
}

// 进行全文检索
func (env *WiserEnv) Search(query string) {
	var results = SearchResults{}
	// 判断查询字符串的长度是否大于 N-gram 中 N 的取值
	if len(query) < env.TokenLen {
		fmt.Println("too short query.")
	} else {
		// 如果长度大于N，就调用函数将词元从查询字符串中提取出来
		err := env.splitQueryToTokens(query)
		if err != nil {
			fmt.Println("将词元从查询字符串中提取出来失败, err: ", err)
		}
		results = *env.searchDocs()
	}
	env.printSearchResults(results)
}

// 从查询字符串中提取出词元的信息
// text 查询字符串
// 返回 queryTokens 按词元编号存储位置信息序列的关联数组
//				    若传入的是空，则新建一个关联数组
func (env *WiserEnv) splitQueryToTokens(query string) error {
	return env.TextToPostingsLists(0, query)
}

// 检索文档
func (env *WiserEnv) searchDocs() *SearchResults {
	if len(env.IIBuffer) < 1 {
		return &SearchResults{}
	}

	return &SearchResults{}
}

// 打印检索结果
func (env *WiserEnv) printSearchResults(result SearchResults) {
	for _, res := range result.Results {
		title, _ := dao.GetDocumentTitle(res.DocumentID)
		fmt.Printf("document id: %d title: %s score: %f\n",
			res.DocumentID, title, res.Score)
	}
	fmt.Printf("Total %d documents are found!\n", len(result.Results))
}
