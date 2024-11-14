package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Address struct {
	Street  string `json:"street"`
	Number  string `json:"number"`
	ZipCode string `json:"zipcode"`
}

type User struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type RestResponse[T any] struct {
	Data     T   `json:"data"`
	Metadata any `json:"_metadata"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use()

	// Get user value
	r.GET("/users/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		user := User{
			Name:    name,
			Surname: "Example",
			Age:     1,
			Address: Address{
				Street:  "mystreet",
				Number:  "mynumber",
				ZipCode: "myzipcode",
			},
		}
		c.IndentedJSON(http.StatusOK, RestResponse[User]{
			Data: user,
			Metadata: map[string]string{
				"_link": "/users/" + name,
			},
		})
	})

	return r
}

func main() {
	r := setupRouter()

	go func() {
		time.Sleep(time.Second)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/manuel?fields=name,surname,address", nil)
		r.ServeHTTP(w, req)
		fmt.Printf(w.Body.String())
		os.Exit(1)
	}()

	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
