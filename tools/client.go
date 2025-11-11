package tools

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// APIResponse 統一的 API 回應格式
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// 建立全域 resty 客戶端
var httpClient = resty.New().
	SetTimeout(30 * time.Second).
	SetRetryCount(2).
	SetRetryWaitTime(1 * time.Second)
