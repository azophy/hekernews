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
	"golang.org/x/crypto/bcrypt"
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
	Id           string `json:"id" xml:"id" form:"id" query:"id"`
	Name         string `json:"name" xml:"name" form:"name" query:"name"`
	Email        string `json:"email" xml:"email" form:"email" query:"email"`
	Username     string `json:"username" xml:"username" form:"username" query:"username"`
	PasswordHash string `query:"password_hash"`
	CreatedAt time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
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
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(200) NOT NULL,
        email VARCHAR(200) NOT NULL,
        username VARCHAR(200) NOT NULL,
        password_hash VARCHAR(200) NOT NULL,
        created_at VARCHAR(200),
        updated_at VARCHAR(200)
    );
    `

	_, err := db_conn.Exec(query)
	return err
}

func main() {
	JWT_SECRET := []byte(os.Getenv("JWT_SECRET"))
	APP_PORT := os.Getenv("APP_PORT")
	if APP_PORT == "" {
		APP_PORT = "3000"
	}

	dbConnection := os.Getenv("DB_CONNECTION")
	if dbConnection == "" {
		dbConnection = "sqlite3"
	}
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		dbURI = "./testdb.sqlite3"
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

		query := `SELECT * FROM users where username = ?`
		rows, err := db_conn.Query(query, username)
		if err != nil {
			return c.String(http.StatusForbidden, "Incorrect username or password")
		}
    user := new(User)
		for rows.Next() {
			var createdAt, updatedAt string
			if err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.PasswordHash, &createdAt, &updatedAt); err != nil {
			  fmt.Println(err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, "Something bad happened on the server")
			}
			user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
			user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
      fmt.Printf("name: %s\n", user.Name)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect username or password")
		}

		// Set custom claims
    user.PasswordHash = "" // remove password
		claims := &jwtCustomClaims{
      *user,
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
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/")
	})

	e.File("/register", "public/register.html")
	e.POST("/register", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		username := c.FormValue("username")
		password := c.FormValue("password")

		timeNow := time.Now().Format(time.RFC3339)
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		_, err = db_conn.Exec("INSERT INTO users(name, email, username, password_hash, created_at, updated_at) values(?, ?, ?, ?, ?, ?)", name, email, username, passwordHash, timeNow, timeNow)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		return c.HTML(http.StatusOK, "Registration succeed. Login <a href='login'>here</a> with your credential")
	})

	// Restricted group
	r := e.Group("/member")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.HTML(http.StatusForbidden, "You can't access this page. Click <a href='login'>here</a> to login")
		},
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
		cookie.Path = "/"
		c.SetCookie(cookie)

		return c.HTML(http.StatusOK, "you have been logged out. back to <a href='/'>home</a>")
	})
	r.File("/new_post", "public/new_post.html")

	e.GET("/api/posts", func(c echo.Context) error {
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

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}
