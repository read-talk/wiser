package logic

import (
	"fmt"
	"github.com/read-talk/wiser/dao"
	"sort"
)

// 将类型 InvertedIndexHash InvertedIndexValue 也用于检索
type QueryTokenHash = InvertedIndexHash
type QueryTokenValue = InvertedIndexValue
type TokenPositionsList = PostingsList

type DocSearchCursor struct {
	Documents *TokenPositionsList // 文档编号的序列
	Current   *TokenPositionsList // 当前的文档编号
}

type PhraseSearchCursor struct {
	Positions []int // 位置信息
	Base      int   // 词元在查询中的位置
	Current   int   // 当前的位置信息
}

type SearchResult struct {
	documentID int     // 检索出的文档编号
	Score      float64 // 检索得分
}

// 检索结果的结构体转化为哈希表
type SearchResultHash struct {
	HashMap map[int]*SearchResult
	Item    []*SearchResult
}

// 进行全文检索 query 查询
func (env *WiserEnv) Search(q string) {
	// 1. 判断查询字符串的长度是否大于 N-gram 中的 N 的取值
	var result = &SearchResultHash{}
	if len(q) < env.TokenLen {
		fmt.Println("too short query.")
		return
	} else { // 2. 如果长度大于N，就将词元从查询字符串中提取出来
		queryTokens, _ := env.splitQueryToTokens(q)
		// 3. 以刚刚提取出来的词元作为参数，开始进行检索处理
		env.searchDocs(queryTokens, result)
	}
	// 4. 打印检索结果
	printSearchResults(result)
}

// 从查询字符串中提取出词元的信息
// text 查询字符串
// query_tokens 按词元编号存储位置信息序列的关联数组
//              若传入的是指向nil的指针，则新建一个关联数组
func (env *WiserEnv) splitQueryToTokens(text string) (*QueryTokenHash, error) {
	// 将文档编号设置为0
	err := env.TextToPostingsLists(0, text)
	return env.IIBuffer, err
}

// 检索文档
// results 检索结果
// tokens 从查询中提取出的词元信息
func (env *WiserEnv) searchDocs(tokens *QueryTokenHash, results *SearchResultHash) {
	var err error
	nTokens := len(tokens.HashMap)
	if nTokens == 0 {
		return
	}
	var cursors = make([]*DocSearchCursor, nTokens)
	for i := range cursors {
		cursors[i] = &DocSearchCursor{}
	}
	// 按照文档频率的升序堆tokens排序
	// 比较出现过词元a和词元b的文档数
	sort.Slice(tokens.Items, func(i, j int) bool {
		return tokens.Items[i].DocsCount < tokens.Items[j].DocsCount
	})
	// 初始化
	if nTokens != 0 {
		var token = tokens.Items[0]
		for i := 0; i < nTokens; i++ {
			if token.TokenID == 0 {
				// 当前的token在构建索引的过程中从未出现过
				return
			}
			cursors[i].Documents, _, err = FetchPostings(token.TokenID)
			if err != nil {
				fmt.Printf("decode postings error! : %d\n", token.TokenID)
				return
			}
			if cursors[i].Documents == nil {
				// 虽然当前的token存在，但是由于更新或删除导致其他倒排类别为空
				return
			}
			cursors[i].Current = cursors[i].Documents
			if i + 1 < nTokens {
				token = tokens.Items[i+1]
			}
		}
		for cursors[0].Current != nil {
			var docId, nextDocId int
			// 将拥有文档最少的词元称为A
			docId = cursors[0].Current.DocumentID
			// 对于除词元A以外的词元，不断获取其下一个DocumentID，
			// 直到当前的document_id不小于词元A的document_id为止
			for i := 1; i < nTokens; i++ {
				cur := cursors[i]
				for cur.Current != nil && cur.Current.DocumentID < docId {
					cur.Current = cur.Current.Next
				}
				if cur.Current == nil {
					return
				}
				// 对于除词元A以外的词元，如果其document_id不等于词元A的document_id
				// 那么就将这个document_id设定为next_doc_id
				if cur.Current.DocumentID != docId {
					nextDocId = cur.Current.DocumentID
					break
				}
			}
			if nextDocId > 0 {
				// 不断获取A的下一个document_id，直到其当前的document_id不小于next_doc_id为止
				for cursors[0].Current != nil && cursors[0].Current.DocumentID < nextDocId {
					cursors[0].Current = cursors[0].Current.Next
				}
			} else {
				score := calcTfIdf(tokens, cursors[0], nTokens, env.IndexedCount)
				addSearchResult(results, docId, score)

				cursors[0].Current = cursors[0].Current.Next
			}
		}
	}
	sort.Slice(results.Item, func(i, j int) bool {
		return results.Item[i].Score > results.Item[j].Score
	})
}

// 以检索结果中的文档编号为查询条件，从文档数据库中取出相应的文档标题，
// 最后输出获取到的标题和检索的得分
func printSearchResults(res *SearchResultHash) {
	if res == nil {
		return
	}
	n := len(res.HashMap)
	for k, v := range res.HashMap {
		r := v
		delete(res.HashMap, k)
		title, _ := dao.GetDocumentTitle(r.documentID)
		fmt.Printf("document_id: %d title: %s score: %.2f\n", r.documentID, title, r.Score)
	}
	fmt.Printf("Total %d document are found!\n", n)
}

// 用TF-IDF计算得分
// query_tokens 查询
// doc_cursors 用于文档检索的游标的集合
// n_query_tokens 查询中的词元数
// indexed_count 建立过索引的文档总数
// return 得分
func calcTfIdf(queryTokens *QueryTokenHash, docCursors *DocSearchCursor,
	nQueryTokens int, indexedCount int) float64 {
	var dcur = docCursors.Current
	var score float64
	for i := 0; i < nQueryTokens; i++ {
		qt := queryTokens.Items[i]
		idf := indexedCount / qt.DocsCount
		score += float64(dcur.PositionsCount * idf)
		dcur = dcur.Next
	}
	return score
}

// 将文档添加到检索结果中
// results 指向检索结果的指针
// document_id 要添加的文档的编号
// score 得分
func addSearchResult(result *SearchResultHash, documentID int, score float64) {
	r := &SearchResult{}
	if v, ok := result.HashMap[documentID]; ok {
		r = v
		r.Score += score
	} else {
		r.documentID = documentID
		r.Score = 0
		result.HashMap = make(map[int]*SearchResult)
		result.HashMap[documentID] = r
	}
}
