package logic

const (
	// bi-gram
	NGram = 2
)

type Page struct {
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

type InvertedIndex struct {
	PostingsMap   map[int][]int
	PostingsCount int // 该词元在所有文档中的出现次数之和
}

type WiserEnv struct {
	TokenLen                int                   // 词元的长度。NGram中N的取值
	IIBuffer                map[int]InvertedIndex // 用于更新倒排索引的缓冲区（Buffer）
	IIBufferUpdateThreshold int                   // 缓冲区中文档数的阈值
	IIBufferCount           int                   // 用户更新倒排索引的缓冲区中的文档数
	IndexedCount            int                   // 建立了索引的文档数
}

func NewEnv(IIBufferUpdateThreshold int) *WiserEnv {
	return &WiserEnv{
		TokenLen:                NGram,
		IIBuffer:                map[int]InvertedIndex{},
		IIBufferUpdateThreshold: IIBufferUpdateThreshold,
		IIBufferCount:           0,
		IndexedCount:            0,
	}
}

// 倒排列表（以文档编号和位置信息为元素的链表结构）
type PostingsList struct {
	DocumentID     int           // 文档编号
	Positions      []int         // 位置信息的数组 TODO: 不知道是啥类型
	PositionsCount int           // 位置信息的条数
	Next           *PostingsList // 指向下一个倒排列表
}

// 倒排索引（以词元编号为键，以倒排列表为值的关联数组）
type InvertedIndexHash struct {
	TokenID      int           // 词元编号（Token ID）
	PostingsList *PostingsList // 指向包含该词元的倒排列表的指针
	//PostingsMap   map[int][]int
	DocsCount     int // 出现过该词元的文档数
	PostingsCount int // 该词元在所有文档中的出现次数之和
	// TODO: 将该结构体转化为哈希表
}

type CompressMethod struct {
	CompressNone   bool // 不压缩
	CompressGolomb bool // 使用 Golomb 编码压缩
}

//type WiserEnv struct {
//	DBPath                  string               // 数据库的路径
//	TokenLen                int                  // 词元的长度。NGram中N的取值
//	Compress                *CompressMethod      // 压缩倒排列表等数据的方法
//	EnablePharseSearch      int                  // 是否进行短语检索
//	IIBuffer                []*InvertedIndexHash // 用于更新倒排索引的缓冲区（Buffer）
//	IIBufferCount           int                  // 用户更新倒排索引的缓冲区中的文档数
//	IIBufferUpdateThreshold int                  // 缓冲区中文档数的阈值
//	IndexedCount            int                  // 建立了索引的文档数
//}
