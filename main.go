package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	HotContentAPI = "https://api.vocus.cc/api/top5-contents"
	SearchAPI     = "https://api.vocus.cc/api/search"
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

func main() {
	s := server.NewMCPServer(
		"Vocus Content Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// 註冊工具 1: 取得熱門內容
	hotTool := mcp.NewTool("get_hot_contents",
		mcp.WithDescription("search Vocus hot contents with structured output"),
		mcp.WithOutputSchema[[]GetHotContentsToolResponse](),
	)
	s.AddTool(hotTool, handleGetHotContents)

	// 註冊工具 2: 搜尋內容
	searchTool := mcp.NewTool("search_content",
		mcp.WithDescription("search Vocus contents by title with structured output"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("title keywords to search for"),
		),
	)
	s.AddTool(searchTool, handleSearchContent)

	log.Println("Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

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

// handleGetHotContents 取得熱門內容
func handleGetHotContents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

type SearchContentAPIResponse struct {
	Contents []struct {
		ContentID string `json:"contentId"`
		Title     string `json:"title"` // 標題
		Type      string `json:"type"`
	} `json:"contents"`
	Creators []struct {
		ID       string `json:"_id"`
		Fullname string `json:"fullname"`
	} `json:"creators"`
	Salons []struct {
		ID   string `json:"_id"`
		Name string `json:"name"`
	} `json:"salons"`
	Tags []string `json:"tags"`
}

type SearchContentToolResponse struct {
	Contents []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"contents"`
	Creators []struct {
		Fullname string `json:"fullname"`
		URL      string `json:"url"`
	} `json:"creators"`
	Salons []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"salons"`
	Tags []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tags"`
}

// handleSearchContent 搜尋內容
func handleSearchContent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := request.RequireString("title")
	if err != nil {
		return mcp.NewToolResultJSON(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("參數錯誤: %v", err),
		})
	}

	var searchResults SearchContentAPIResponse

	resp, err := httpClient.R().
		SetContext(ctx).
		SetQueryParam("title", title).
		SetResult(&searchResults).
		Get(SearchAPI)
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

	var searchContentToolResponse SearchContentToolResponse

	for _, content := range searchResults.Contents {
		contentURL := fmt.Sprintf("https://vocus.cc/%s/%s", content.Type, content.ContentID)
		searchContentToolResponse.Contents = append(searchContentToolResponse.Contents, struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}{
			Title: content.Title,
			URL:   contentURL,
		})
	}

	for _, creator := range searchResults.Creators {
		creatorURL := fmt.Sprintf("https://vocus.cc/user/%s", creator.ID)
		searchContentToolResponse.Creators = append(searchContentToolResponse.Creators, struct {
			Fullname string `json:"fullname"`
			URL      string `json:"url"`
		}{
			Fullname: creator.Fullname,
			URL:      creatorURL,
		})
	}

	for _, salon := range searchResults.Salons {
		salonURL := fmt.Sprintf("https://vocus.cc/salon/%s", salon.ID)
		searchContentToolResponse.Salons = append(searchContentToolResponse.Salons, struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			Name: salon.Name,
			URL:  salonURL,
		})
	}

	for _, tag := range searchResults.Tags {
		tagURL := fmt.Sprintf("https://vocus.cc/tags/%s", tag)
		searchContentToolResponse.Tags = append(searchContentToolResponse.Tags, struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			Name: tag,
			URL:  tagURL,
		})
	}

	return mcp.NewToolResultJSON(APIResponse{
		Success: true,
		Data:    searchContentToolResponse,
	})
}
