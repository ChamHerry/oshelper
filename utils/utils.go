package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ChamHerry/oshelper/consts"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// ParseToJSON 解析字符串为 JSON 格式
func ParseToJSON(data string, key string) (string, error) {
	parsedData := make(map[string]string)
	// 按行分割数据
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		// 跳过空行
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 按 key 分割键值对
		parts := strings.SplitN(line, key, 2)
		if len(parts) != 2 {
			continue
		}
		// g.Log().Debugf(context.Background(), "parts:%v", parts)
		// 去掉键和值的多余引号和空格
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		parsedData[CamelToUnderline(key)] = value
	}
	// 转换为 JSON 格式
	jsonData, err := json.MarshalIndent(parsedData, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err) // 使用 %w 以便错误链
	}
	// g.Log().Debugf(context.Background(), "jsonData:%v", string(jsonData))
	return string(jsonData), nil
}

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
