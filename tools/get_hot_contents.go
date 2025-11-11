package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	HotContentAPI = "https://api.vocus.cc/api/top5-contents"
)

type HotContentsAPIResponse struct {
	Type    string `json:"type"`
	Content struct {
		ID       string `json:"_id"`
		Title    string `json:"title"`    // 標題
		Abstract string `json:"abstract"` // 摘要
		User     struct {
			FullName string `json:"fullname"` // 暱稱
		} `json:"user"` // 創作者資料
		Salon struct {
			Name string `json:"name"` // 沙龍名稱
		} `json:"salon"` // 所屬沙龍資料
	} `json:"content"`
}

type GetHotContentsToolResponse struct {
	Title    string `json:"title"`    // 標題
	Abstract string `json:"abstract"` // 摘要
	Author   string `json:"author"`   // 創作者暱稱
	Salon    string `json:"salon"`    // 所屬沙龍名稱
	URL      string `json:"url"`      // 內容連結
}

// GetHotContentsTool 取得熱門內容
func GetHotContentsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var contents []*HotContentsAPIResponse

	resp, err := httpClient.R().
		SetContext(ctx).
		SetResult(&contents).
		Get(HotContentAPI)
	if err != nil {
		return mcp.NewToolResultJSON(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("無法連接到 API: %v", err),
		})
	}

	if !resp.IsSuccess() {
		return mcp.NewToolResultJSON(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("API 回應錯誤，狀態碼: %d, 內容: %s", resp.StatusCode(), resp.String()),
		})
	}

	toolResponse := make([]GetHotContentsToolResponse, 0, len(contents))
	for _, item := range contents {
		contentURL := fmt.Sprintf("https://vocus.cc/%s/%s", item.Type, item.Content.ID)
		toolResponse = append(toolResponse, GetHotContentsToolResponse{
			Title:    item.Content.Title,
			Abstract: item.Content.Abstract,
			Author:   item.Content.User.FullName,
			Salon:    item.Content.Salon.Name,
			URL:      contentURL,
		})
	}

	return mcp.NewToolResultJSON(APIResponse{
		Success: true,
		Data:    toolResponse,
	})
}
