package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/manuelarte/milogo"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(milogo.Milogo())

	// Get user value
	r.GET("/users", func(c *gin.Context) {
		users := []*User{
			{
				Name:    "John",
				Surname: "Doe",
				Age:     99,
			},
			{
				Name:    "Milo",
				Surname: "Doe",
				Age:     99,
			},
		}
		c.IndentedJSON(http.StatusOK, &users)
	})

	return r
}

func main() {
	r := setupRouter()

	go func() {
		time.Sleep(time.Second)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users?fields=name,surname", nil)
		r.ServeHTTP(w, req)
		fmt.Println(w.Body.String())
		os.Exit(1)
	}()

	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
