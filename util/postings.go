package util

import (
	"encoding/json"
)

// 将倒排列表转换成字节序列
func EncodePostings(postingsMap map[int][]int) (string, error) {
	buf, err := json.Marshal(postingsMap)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// 对倒排列表进行还原或解码
func DecodePostings(buf string) (map[int][]int, error) {
	postingsMap := map[int][]int{}

	err := json.Unmarshal([]byte(buf), postingsMap)
	if err != nil {
		return nil, err
	}

	return postingsMap, nil
}

// 获取将两个倒排列表合并后得到的倒排列表
func MergePostings(pa, pb map[int][]int) map[int][]int {
	mergePostings := map[int][]int{}
	allKeysSet := NewSet()

	for key := range pa {
		allKeysSet.Add(key)
	}
	for key := range pb {
		allKeysSet.Add(key)
	}

	for _, key := range allKeysSet.List() {
		subSet := NewSet()
		al, ok := pa[key]
		if ok {
			subSet.Add(al...)
		}
		bl, ok := pb[key]
		if ok {
			subSet.Add(bl...)
		}
		mergePostings[key] = subSet.SortList()
	}
	return mergePostings
}
