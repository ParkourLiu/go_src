package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

var (
	factory *guidFactory
	r       *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	factory = NewGUIDFactory(nodeNo)
}

// StringService provides operations on strings.
type IdGeneraterService interface {
	GenerateUniqueIdV1(uint32) ([]string, error)
	//	GenerateUniqueIdInt64V1(uint32) ([]int64, error)
}

type idGeneraterService struct{}

func (service idGeneraterService) GenerateUniqueIdInt64V1(count uint32) ([]int64, error) {
	if count <= 0 {
		return []int64{}, nil
	}

	//generate id
	suffixIds := []int64{}
	for {
		id, err := factory.NewGUID()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// 混淆处理。添加最后一位随机数
		id = mix(id)

		suffixIds = append(suffixIds, id)
		if uint32(len(suffixIds)) >= count {
			break
		}
	}

	return suffixIds, nil
}

func (service idGeneraterService) GenerateUniqueIdV1(count uint32) ([]string, error) {
	if count <= 0 {
		return []string{}, nil
	}

	//generate id
	ids, _ := service.GenerateUniqueIdInt64V1(count)
	suffixIds := []string{}
	var idStr string
	for _, id := range ids {
		//		idStr, _ = addZeroForNum(strconv.FormatInt(id, 36), 23)
		idStr = strconv.FormatInt(id, 36)
		idStr, _ = addZeroForNum(idStr, 14)
		suffixIds = append(suffixIds, idStr)
	}

	return suffixIds, nil
}

func mix(src int64) int64 {
	return src*int64(10) + int64(r.Intn(10))
}

func addZeroForNum(str string, strLength int) (string, error) {
	var buffer bytes.Buffer
	strLen := len(str)
	if strLength <= strLen {
		return str, nil
	}
	for i := 0; i < strLength-strLen; i++ {
		buffer.WriteString("0")
	}
	buffer.WriteString(str)
	return buffer.String(), nil
}
