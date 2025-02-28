#!/bin/bash

# 检查是否提供参数
if [ $# -eq 0 ]; then
    echo "用法: $0 <ConfigName>"
    exit 1
fi

CONFIG_NAME=$1
CONFIG_NAME_LOWER=$(echo "$CONFIG_NAME" | tr '[:upper:]' '[:lower:]')
# 将参数转换为下划线命名方式（snake_case）
CONFIG_NAME_SNAKE=$(echo "$CONFIG_NAME" | sed -r 's/([A-Z]+)/_\1/g' | tr '[:upper:]' '[:lower:]' | sed 's/^_//')

# 创建目录
mkdir -p examples/${CONFIG_NAME_LOWER}

# 创建 .go 文件
touch examples/${CONFIG_NAME_LOWER}/${CONFIG_NAME_SNAKE}.go

echo "package main
import \"fmt\"

// go run examples/${CONFIG_NAME_LOWER}/${CONFIG_NAME_SNAKE}.go
func main() {
    fmt.Println(\"Hello, World!\")
}
" > examples/${CONFIG_NAME_LOWER}/${CONFIG_NAME_SNAKE}.go

echo "已创建 examples/${CONFIG_NAME_LOWER}/${CONFIG_NAME_SNAKE}.go"
