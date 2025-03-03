package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ChamHerry/oshelper/consts"
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

// 解析RPM包名
// 	regex := regexp.MustCompile()

var rpmRegex = regexp.MustCompile(`^(.+)-([0-9A-Za-z\._]+(?:\.[0-9]+)*)-([0-9]+)\.([^.]+(?:\.[^.]+)*)\.([^.]+)\.rpm$`)

func ParseRPMName(rpmName string) (consts.PackageInfo, error) {
	packageInfo := consts.PackageInfo{}
	// 使用预编译的正则表达式
	match := rpmRegex.FindStringSubmatch(rpmName)
	if match == nil {
		return packageInfo, fmt.Errorf("invalid RPM file name format: %s", rpmName)
	}
	if len(match) != 6 {
		return consts.PackageInfo{}, fmt.Errorf("invalid RPM name: %s", rpmName)
	}
	// 提取匹配的分组
	packageInfo.Name = match[1]          // 包名
	packageInfo.Version = match[2]       // 版本号
	packageInfo.ReleaseNumber = match[3] // 发布号
	packageInfo.OS = match[4]            // 操作系统
	packageInfo.Architecture = match[5]  // 架构
	// 处理包名可能包含子版本的情况，比如 `antlr3-C`
	if strings.Contains(packageInfo.Name, "-") {
		subparts := strings.Split(packageInfo.Name, "-")
		if len(subparts) > 1 && isVersion(subparts[len(subparts)-1]) {
			packageInfo.Name = strings.Join(subparts[:len(subparts)-1], "-")
			packageInfo.Version = subparts[len(subparts)-1] + "." + packageInfo.Version
		}
	}
	return packageInfo, nil
}

// 判断是否是版本号的辅助函数
func isVersion(value string) bool {
	regex := regexp.MustCompile(`^[0-9]+(\.[0-9]+)*$`)
	return regex.MatchString(value)
}
