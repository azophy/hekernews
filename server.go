package main

import (
  "net/http"
  "time"

	"github.com/labstack/echo/v4"
)

type Post struct {
  id  string
	Title  string `json:"title" xml:"title" form:"title" query:"title"`
	Content string `json:"content" xml:"content" form:"content" query:"content"`
  CreatedAt time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
  UpdatedAt time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

func main() {
	e := echo.New()

  e.File("/", "public/index.html")

	//e.GET("/", func(c echo.Context) error {
		//return c.String(http.StatusOK, "Hello, World!")
	//})
  e.GET("/api/posts", func(c echo.Context) error {
    posts := []Post{
      {"1","Test aja 1", "content 1", time.Now(), time.Now() },
      {"2","Test aja 2", "content 2", time.Now(), time.Now() },
      {"3","Test aja 3", "content 3", time.Now(), time.Now() },
      {"4","Test aja 4", "content 4", time.Now(), time.Now() },
      {"5","Test aja 5", "content 5", time.Now(), time.Now() },
    }
		return c.JSON(http.StatusOK, posts)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
