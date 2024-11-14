package main

type Address struct {
	Street string `json:"street"`
	Number string `json:"number"`
}

type User struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type Response[T any] struct {
	Data     T   `json:"data"`
	Metadata any `json:"_metadata"`
}
