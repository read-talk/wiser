package util

import (
	"fmt"
	"time"
)

// 检查输入的字符（UTF-32）是否不属于索引对象
// ustr 输入的字符
// 返回是否是空白字符 true: 是空白字符，false: 不是空白字符
func IsIgnoredChar(c rune) bool {
	switch c {
	case ' ', '\f', '\n', '\r', '\t',
		'!', '"', '#', '$', '%', '&', '\'',
		'(', ')', '*', '+', ',', '-', '.',
		'/', ':', ';', '<', '=', '>', '?',
		'@', '[', '\\', ']', '^', '_', '`',
		'{', '|', '}', '~',
		'、', '。', '（', '）', '！', '，', '：', '；', '“', '”',
		'a', 'b', 'c', 'd', 'e', 'f', 'g',
		'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return true
	default:
		return false
	}
}

// 将输入的字符串分隔为N-gram
// ustr 输入的字符串
// n N-gram 中 N 的取值。建议将其设为大于 1 的值
// start 词元的起始位置
// 返回分隔出来的词元的长度
func NgramNext(ustr []rune, start *int, n int) (int, int) {
	totalLen := len(ustr)
	// 读取时跳过文本开头的空格等字符
	for {
		if *start >= totalLen {
			break
		}
		// 当不是空白字符的时候就跳出循环
		if !IsIgnoredChar(ustr[*start]) {
			break
		}
		*start++
	}
	tokenLen := 0
	position := *start

	// 不断取出最多包含n个字符的词元，直到遇到不属于索引对象的字符或到达了字符串的尾部
	for {
		if *start >= totalLen {
			break
		}
		if tokenLen >= n {
			break
		}
		// 当是空白字符的时候就结束索引
		if IsIgnoredChar(ustr[*start]) {
			break
		}
		*start++
		tokenLen++
	}

	if tokenLen >= n {
		*start = position + 1
	}

	return tokenLen, position
}

var preTime *time.Time

func PrintTimeDiff() {
	currentTime := time.Now()
	if preTime != nil {
		timeDiff := currentTime.UnixNano() - preTime.UnixNano()
		fmt.Printf("[time] %s (diff %d)\n", currentTime, timeDiff)
	} else {
		fmt.Printf("[time] %s\n", currentTime)
	}
	preTime = &currentTime
}
