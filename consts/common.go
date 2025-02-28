package consts

import (
	"context"
	"time"
)

const (
	DefaultDialTimeout            = 2 * time.Second
	DefaultRunCommandFailedCounts = 3
)

// RunCommandConfig 用于存储远程运行参数
type RunCommandConfig struct {
	Command                string
	RunCommandFailedCounts int
}

// 异步调用参数
type AsyncCallParam struct {
	Concurrency int
	Operation   func(ctx context.Context, item interface{}) (interface{}, error)
	Items       []interface{}
}

type AsyncCallResult struct {
	Ret interface{}
	Err error
}

type AsyncCallResultList struct {
	RetList []AsyncCallResult
}
