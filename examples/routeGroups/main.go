package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/manuelarte/milogo"
	"github.com/manuelarte/milogo/pkg/config"
)

var (
	manuel = User{
		Name:    "Manuel",
		Surname: "Example",
		Age:     99,
		Address: Address{
			Street: "mystreet",
			Number: "mynumber",
		},
	}
	milo = User{
		Name:    "Milo",
		Surname: "Example",
		Age:     0,
		Address: Address{
			Street: "mystreet",
			Number: "mynumber",
		},
	}

	usersByName = map[string]User{"manuel": manuel, "milo": milo}
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	usersGroup := r.Group("/users", milogo.Milogo())
	usersGroup.GET("", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, []User{manuel, milo})
	})
	usersGroup.GET("/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		if user, ok := usersByName[name]; ok {
			c.IndentedJSON(http.StatusOK, user)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	milogoOption, _ := config.WithWrapField("data")
	wrappedUsersGroup := r.Group("/wrapped/users", milogo.Milogo(milogoOption))
	wrappedUsersGroup.GET("", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, Response[[]User]{
			Data:     []User{manuel, milo},
			Metadata: map[string]string{"link": "/wrapped/users"},
		})
	})
	wrappedUsersGroup.GET("/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		if user, ok := usersByName[name]; ok {
			c.IndentedJSON(http.StatusOK, Response[User]{
				Data:     user,
				Metadata: map[string]string{"link": "/wrapped/users/" + name},
			})
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	return r
}

func main() {
	r := setupRouter()

	ctx := context.Background()
	go func() {
		time.Sleep(time.Second)
		w := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(ctx, "GET", "/users?fields=name", nil)
		r.ServeHTTP(w, req)
		fmt.Println(w.Body.String())

		w = httptest.NewRecorder()
		req, _ = http.NewRequestWithContext(ctx, "GET", "/users/manuel?fields=name,surname", nil)
		r.ServeHTTP(w, req)
		fmt.Println(w.Body.String())

		w = httptest.NewRecorder()
		req, _ = http.NewRequestWithContext(ctx, "GET", "/wrapped/users/manuel?fields=name,surname", nil)
		r.ServeHTTP(w, req)
		fmt.Println(w.Body.String())
		os.Exit(1)
	}()

	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
