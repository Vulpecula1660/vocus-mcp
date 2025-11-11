# Vocus MCP Server

這是一個基於 [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) 的 **MCP (Model Context Protocol)** 伺服器實作範例。
此伺服器透過 [Vocus API](https://vocus.cc) 提供兩個工具 (Tools)：

1. 🔥 取得 Vocus 平台的 Top 5 熱門內容
2. 🔍 根據標題關鍵字搜尋內容

伺服器以 **Streamable HTTP** 模式運行，可作為 MCP client（如 LLM 插件、Agent 或開發者工具）呼叫的後端服務。

---

## 🚀 功能介紹

### 1. 取得熱門內容 (`get_hot_contents`)

* 說明：從 `https://api.vocus.cc/api/top5-contents` 取得目前 Vocus 上的熱門內容。
* 回傳內容：

  * 標題 (`title`)
  * 摘要 (`abstract`)
  * 創作者 (`author`)
  * 所屬沙龍 (`salon`)
  * 內容連結 (`url`)

**範例輸出：**

```json
{
  "success": true,
  "data": [
    {
      "title": "最新文章標題",
      "abstract": "文章摘要",
      "author": "作者名稱",
      "salon": "沙龍名稱",
      "url": "https://vocus.cc/article/xxxxx"
    }
  ]
}
```

---

### 2. 搜尋內容 (`search_content`)

* 說明：根據輸入的標題關鍵字搜尋內容。
* 參數：

  * `title` (string, required)：搜尋關鍵字。
* 呼叫目標 API：`https://api.vocus.cc/api/search?title={關鍵字}`
* 回傳內容：

  * 相關內容、創作者、沙龍與標籤的清單及對應連結。

**範例輸出：**

```json
{
  "success": true,
  "data": {
    "contents": [
      { "title": "文章標題", "url": "https://vocus.cc/article/xxxx" }
    ],
    "creators": [
      { "fullname": "作者名稱", "url": "https://vocus.cc/user/xxxx" }
    ],
    "salons": [
      { "name": "沙龍名稱", "url": "https://vocus.cc/salon/xxxx" }
    ],
    "tags": [
      { "name": "標籤名稱", "url": "https://vocus.cc/tags/xxxx" }
    ]
  }
}
```

---

## 🧩 架構說明

* **主程式**：建立一個 MCP 伺服器，並以 Streamable HTTP 啟動。
* **Tools**：

  * `get_hot_contents`：呼叫 Vocus 熱門內容 API。
  * `search_content`：呼叫 Vocus 搜尋 API。
* **錯誤處理**：

  * 所有錯誤皆會包裝在 `APIResponse` 中回傳。
* **HTTP Client**：

  * 使用 [`resty`](https://github.com/go-resty/resty) 作為 HTTP 客戶端。
  * 設定 30 秒逾時、最多重試 2 次。

---

## ⚙️ 環境需求

* Go 1.21 或以上
* 網路可連線到 `https://api.vocus.cc`

---

## 📦 安裝與執行

### 1️⃣ 下載專案

```bash
git clone https://github.com/yourname/vocus-mcp-server.git
cd vocus-mcp-server
```

### 2️⃣ 安裝相依套件

```bash
go mod tidy
```

### 3️⃣ 執行伺服器

```bash
go run main.go
```

伺服器啟動後會在 `http://localhost:8080` 運行，並輸出：

```
Starting StreamableHTTP server on :8080
```

---

## 🤖 與各 AI Client 整合方式

此 MCP Server 可直接整合到支援 **Model Context Protocol (MCP)** 的多種 AI 工具中。

---

### 🍒 Cherry Studio

1. 在 `"MCP"` 的 JSON 設定增加以下配置：

   ```json
   {
     "mcpServers": {
       "vocus": {
         "name": "Vocus Content Server",
         "type": "streamableHttp",
         "baseUrl": "http://localhost:8080/mcp"
       }
     }
   }
   ```
2. 重新啟動 Cherry Studio。
   完成後，你即可在 Cherry 的工具清單中找到 **Vocus Content Server**，並直接呼叫：

   * `get_hot_contents`
   * `search_contents`

---

## 🧠 MCP Tool 註冊概念

此伺服器透過 `mark3labs/mcp-go` 的 `server.AddTool` 註冊兩個工具：

```go
s.AddTool(hotTool, handleGetHotContents)
s.AddTool(searchTool, handleSearchContents)
```

每個 Tool 都包含：

* 工具名稱 (`get_hot_contents`, `search_contents`)
* 描述 (`WithDescription`)
* 輸入與輸出格式 (`WithString`, `WithOutputSchema`)
* 對應處理函式（負責呼叫外部 API 並回傳 JSON 結果）

---

## 📚 主要套件

| 套件                                                                 | 用途              |
| ------------------------------------------------------------------ | --------------- |
| [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) | MCP Protocol 實作 |
| [github.com/go-resty/resty/v2](https://github.com/go-resty/resty)  | HTTP 客戶端        |
| Go 標準庫 `context`, `log`, `fmt`, `time`                             | 系統基本操作          |

---
