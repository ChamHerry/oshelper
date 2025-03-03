package utils

import (
	"context"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/ChamHerry/oshelper/consts"
)

// CamelToUnderline 驼峰转下划线
func CamelToUnderline(s string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	out := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(out)
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// AsyncCall 异步调用函数
func AsyncCall(ctx context.Context, param consts.AsyncCallParam) (ret consts.AsyncCallResultList, err error) {
	var wg sync.WaitGroup
	if param.Concurrency <= 0 {
		param.Concurrency = runtime.NumCPU()
	}
	sem := make(chan struct{}, param.Concurrency)
	results := make(chan consts.AsyncCallResult, len(param.Items))
	for _, item := range param.Items {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// 执行操作
			itemRet, err := param.Operation(ctx, item)
			// 将结果发送到通道
			results <- consts.AsyncCallResult{
				Ret: itemRet,
				Err: err,
			}
		}(item)
	}
	wg.Wait()
	close(results)
	for result := range results {
		ret.RetList = append(ret.RetList, result)
	}
	return ret, nil
}

// ConvertSliceToInterfaceSlice 使用 reflect 将任何类型的切片转换为 []interface{}
func ConvertSliceToInterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("ConvertSliceToInterfaceSlice: input is not a slice")
	}

	s := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		s[i] = v.Index(i).Interface()
	}
	return s
}
