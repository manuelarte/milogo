package main

import (
	"fmt"
	"github.com/manuelarte/milogo"
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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(milogo.Milogo())

	// Get user value
	r.GET("/users/:name", func(c *gin.Context) {
		user := User{
			Name:    c.Params.ByName("name"),
			Surname: "Example",
			Age:     1,
			Address: Address{
				Street:  "mystreet",
				Number:  "mynumber",
				ZipCode: "myzipcode",
			},
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
		req, _ := http.NewRequest("GET", "/users/manuel?fields=name,surname,address", nil)
		r.ServeHTTP(w, req)
		fmt.Printf("All the address fields:\n%s", w.Body.String())

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/users/manuel?fields=name,surname,address(number,zipcode)", nil)
		r.ServeHTTP(w, req)
		fmt.Printf("Some address fields:\n%s", w.Body.String())

		os.Exit(1)
	}()

	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
