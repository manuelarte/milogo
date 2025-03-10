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
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Age             int    `json:"age"`
	FavouriteColour string `json:"favouriteColour"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(milogo.Milogo())

	// Get user value
	r.GET("/users/:name", func(c *gin.Context) {
		user := User{
			Name:            c.Params.ByName("name"),
			Surname:         "Doe",
			Age:             18,
			FavouriteColour: "red",
		}
		c.IndentedJSON(http.StatusOK, user)
	})

	return r
}

func main() {
	r := setupRouter()

	go func() {
		time.Sleep(time.Second)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/john?fields=name,surname", nil)
		r.ServeHTTP(w, req)
		fmt.Println(w.Body.String())
		os.Exit(1)
	}()

	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
