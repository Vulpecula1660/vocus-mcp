package main

import (
	"log"

	"mcp/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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
		mcp.WithOutputSchema[[]tools.GetHotContentsToolResponse](),
	)
	s.AddTool(hotTool, tools.GetHotContentsTool)

	// 註冊工具 2: 搜尋內容
	searchTool := mcp.NewTool("search_contents",
		mcp.WithDescription("search Vocus contents by title with structured output"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("title keywords to search for"),
		),
	)
	s.AddTool(searchTool, tools.SearchContentsTool)

	log.Println("Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
