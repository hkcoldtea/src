# web.go
A simple Go-lang Web framework.
<br/><br/>一個基於 Go-lang 套件包 net/http 的簡易 Web 框架，實現了靜態、動態路由映射，分組映射，靜態文件傳輸，模板渲染。
## Getting Start
支持動態、靜態路由，分組路由配置，以及靜態文件獲取，模板渲染頁面。
***
### 靜態路由：
```golang
func main() {
    s := server.InitServer()
    s.Get("/level1/level2", func(c *server.Context) {
		  c.JSON(http.StatusOK, server.Content{
        "username": "yanyibin",
        "password": "yyb",
      })
    })
    s.Run("localhost:9999")
}
```
### 動態路由：
可由Context中PathParams中獲得動態路由名
<br/>利用trie樹進行實現
```golang
func main() {
    s := server.InitServer()
    s.Get("/level1/:v1", func(c *server.Context) {
      c.JSON(http.StatusOK, server.Content{
        "pathParam": c.PathParams["v1"]
        "username": "yanyibin",
        "password": "yyb",
      })
    })
    s.Run("localhost:9999")
}
```
### 分組路由：
```golang
func main() {
    s := server.InitServer()
    g := s.SetGroup("/group1")
    {
        g.Get("/level1/:v1", func(c *server.Context) {
          c.JSON(http.StatusOK, server.Content{
            "pathParam": c.PathParams["v1"]
            "username": "yanyibin",
            "password": "yyb",
          })
        })

        g.Get("/level2/:v2", func(c *server.Context) {
          c.JSON(http.StatusOK, server.Content{
            "pathParam": c.PathParams["v2"]
            "username": "yanyibin",
            "password": "yyb",
          })
        })
    }
    s.Run("localhost:9999")
}
```
### 靜態文件訪問、模板渲染：
通過localhost:9999/student 即可獲取test.tmpl對應頁面。
```golang
type student struct {
	Name string
	Age  int
}

func main() {
    s := server.InitServer()
    s.LoadTemplate("test/templates/*")
    s.StaticResource("/static/css", "test/static")

    s1 := &student{Name: "yanyibin", Age: 23}
	  s2 := &student{Name: "ty", Age: 23}

    s.Get("/student", func(c *server.Context) {
      c.HTML(http.StatusOK, "test.tmpl", server.Content{
        "title":    "yanyibin",
        "students": [2]*student{s1, s2},
      })
    })
    s.Run("localhost:9999")
}
```
***
## 目錄結構
### 目錄結構描述
```
.
├── README.md                   // 讀我檔案/敘述文件
├── server
│   ├── context.go              // 請求上下文
│   ├── cookie.go               // 儲存在用戶端的瀏覽資訊
│   ├── group.go                // 服務url前綴分組
│   ├── router.go               // 請求路由
│   ├── server.go               // 服務相關
│   ├── error.go                // 顯示功能失效
|
├── util
|   |── file.go                 // 檔案系統相關
|   |── string.go               // 字符串處理工具
|   |── time.go                 // 與時間資訊相關
|   |── trie.go                 // 實現動態路由 trie樹
|
├── test                        // 靜態文件測試用包
|   |── static                  // html & css
|   |── templates               // tmpl模板
|
├── test.go                     // 測試啓動的程序
```
