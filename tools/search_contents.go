package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	SearchAPI = "https://api.vocus.cc/api/search"
)

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

// SearchContentsTool 搜尋內容
func SearchContentsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
