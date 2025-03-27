package apiserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost         string `env:"DB_HOST" envDefault:"localhost"`
	DBPort         string `env:"DB_PORT" envDefault:"5432"`
	DBUser         string `env:"DB_USER" envDefault:"postgres"`
	DBPassword     string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName         string `env:"DB_NAME" envDefault:"postgres"`
	ExternalAPIURL string `env:"EXTERNAL_API_URL" envDefault:"https://jsonplaceholder.typicode.com"`
	APIPort        string `env:"API_PORT" envDefault:"8080"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func CloseDB() error {
	if db != nil {
		return db.Close()
	}

	return nil
}

func migrateDB(cfg Config) error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	commands := []string{
		"psqldef",
		fmt.Sprintf("-U%s", cfg.DBUser),
		fmt.Sprintf("-W%s", cfg.DBPassword),
		fmt.Sprintf("-p%s", cfg.DBPort),
		cfg.DBName,

		fmt.Sprintf("-f%s", filepath.Join(dir, "schema.sql")),
	}

	cmd := exec.Command(commands[0], commands[1:]...)
	out, err := cmd.CombinedOutput()

	log.Println("[psqldef] ", string(out))

	if err != nil {
		return err
	}

	return nil
}

func InitServer() (*gin.Engine, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		return nil, err
	}

	// DBのマイグレーションを実行
	if err := migrateDB(cfg); err != nil {
		return nil, err
	}

	r := gin.Default()

	// example: DBから取得するエンドポイント
	r.GET("/users", func(c *gin.Context) {
		// DBから取得する
		rows, err := db.Query("SELECT id, name FROM users")

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			rows.Scan(&user.ID, &user.Name)
			users = append(users, user)
		}
		c.JSON(200, gin.H{
			"users": users,
		})
	})

	// example: userを作成するエンドポイント
	r.POST("/user", func(c *gin.Context) {

		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err = db.Exec("INSERT INTO users (name) VALUES ($1)", user.Name)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "user created"})
	})

	// example: 外部サービスと通信するエンドポイント
	r.GET("/external-request", func(c *gin.Context) {
		resp, err := http.Get(cfg.ExternalAPIURL + "/posts")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		var data any
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"posts": data,
		})
	})
	return r, nil
}
