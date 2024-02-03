package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
)

type Post struct {
	Id        string    `json:"id" xml:"id" form:"id" query:"id"`
	Title     string    `json:"title" xml:"title" form:"title" query:"title"`
	Content   string    `json:"content" xml:"content" form:"content" query:"content"`
	CreatedAt time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

type User struct {
	Id       string `json:"id" xml:"id" form:"id" query:"id"`
	Name     string `json:"name" xml:"name" form:"name" query:"name"`
	Username string `json:"username" xml:"username" form:"username" query:"username"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	User
	jwt.RegisteredClaims
}

func migrateDb(db_conn *sql.DB) error {
	fmt.Println("attempting migration....")
	query := `
    CREATE TABLE IF NOT EXISTS posts(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title VARCHAR(200),
        content TEXT,
        created_at VARCHAR(200),
        updated_at VARCHAR(200)
    );
    `

	_, err := db_conn.Exec(query)
	return err
}

func main() {
	JWT_SECRET := []byte("secret")

	dbConnection := os.Getenv("DB_CONNECTION")
	if dbConnection == "" {
		dbConnection = "sqlite3"
	}
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		dbURI := "./testdb.sqlite3"
	}
	db_conn, err := sql.Open(dbConnection, dbURI)
	if err != nil {
		panic(err)
	}
	defer db_conn.Close()

	migrateDb(db_conn)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.File("/", "public/index.html")

	e.GET("/login", func(c echo.Context) error {
		return c.File("public/login.html")
	})
	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Throws unauthorized error
		if username != "fulan" || password != "fulan" {
			return echo.ErrUnauthorized
		}

		// Set custom claims
		claims := &jwtCustomClaims{
			User{Id: "id-1", Name: "Bapak Fulan", Username: "Fulan", Email: "fulan@example.com"},
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString(JWT_SECRET)
		if err != nil {
			return err
		}

		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = t
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/")
	})

	// Restricted group
	r := e.Group("/member")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		//ErrorHandler: func(c echo.Context, err error) error {
		//return c.Redirect(http.StatusSeeOther, "/login")
		//},
		TokenLookup: "cookie:token",
		SigningKey:  JWT_SECRET,
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("/whoami", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		return c.JSON(http.StatusOK, claims)
	})
	r.GET("/logout", func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = ""
		cookie.MaxAge = -1 // https://stackoverflow.com/a/59736764
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/")
	})
	r.File("/new_post", "public/new_post.html")

	e.GET("/api/posts", func(c echo.Context) error {
		//posts := []Post{
		//{"1","Test aja 1", "content 1", time.Now(), time.Now() },
		//{"2","Test aja 2", "content 2", time.Now(), time.Now() },
		//{"3","Test aja 3", "content 3", time.Now(), time.Now() },
		//{"4","Test aja 4", "content 4", time.Now(), time.Now() },
		//{"5","Test aja 5", "content 5", time.Now(), time.Now() },
		//}
		rows, err := db_conn.Query("SELECT * FROM posts")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer rows.Close()

		var all []Post
		for rows.Next() {
			var post Post
			var createdAt, updatedAt string
			if err := rows.Scan(&post.Id, &post.Title, &post.Content, &createdAt, &updatedAt); err != nil {
				fmt.Println(err.Error())
				return err
			}

			post.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
			post.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

			all = append(all, post)
		}

		return c.JSON(http.StatusOK, all)
	})

	e.POST("/api/posts", func(c echo.Context) error {
		title := c.FormValue("title")
		content := c.FormValue("content")
		timeNow := time.Now().Format(time.RFC3339)

		_, err := db_conn.Exec("INSERT INTO posts(title, content, created_at, updated_at) values(?, ?, ?, ?)", title, content, timeNow, timeNow)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		return c.Redirect(http.StatusSeeOther, "/")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
