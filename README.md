[![Go](https://github.com/manuelarte/milogo/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/milogo/actions/workflows/go.yml)
![coverage](https://raw.githubusercontent.com/manuelarte/milogo/badges/.badges/main/coverage.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/milogo)](https://goreportcard.com/report/github.com/manuelarte/milogo)
![version](https://img.shields.io/github/v/release/manuelarte/milogo)
# Milogo
Rest Partial Response (aka Field Selection) Pattern middleware for [Gin](https://gin-gonic.com/). This gin middleware allows you to select a subset of fields to be returned from your endpoints.

e.g. Imagine that you have the following rest endpoint that returns an user with the fields, `id, name, surname, age, address`:
> /users/1
```json
{
 "id": 1,
 "name": "John",
 "surname": "Doe",
 "age": 18,
 "address": {
   "street": "mystreet",
   "city": "mycity",
   "country": "mycountry",
   "zipcode": "1111"
 }
}
```
We can call the endpoint and, with the query parameter fields, filter out the fields that we are interested:
> /users/1?**fields=name,surname**
```json
{
 "name": "John",
 "surname": "Doe"
}
```

## 📝 How To Install It And Use It

- Run the command:

> go get -u -d github.com/manuelarte/milogo

- Add milogo middleware
```go
r := gin.Default()
r.Use(Milogo())
```

- Call your endpoints adding the query parameter `fields` with the fields you want to filter:

> /users/1?**fields=name,surname**


## ✨ Features

- [Support for multiple fields filtering](./examples/simple). 

> /users/1?fields=name,surname
```json
{
 "name": "John",
 "surname": "Doe"
}
```

- [Support for arrays](./examples/simple-array)

> /users?fields=name
```json
[
  {
    "name": "John"
  }
]
```

- [Support for nested jsons](./examples/nested).

> /users/1?fields=name,surname,address(street,zipcode)
```json
{
 "name": "John",
 "surname": "Doe",
 "address": {
   "street": "mystreet",
   "zipcode": "myzipcode"
 }
}
```

- [Support for wrapped json](./examples/wrapped). 
> /users/1?fields=name
```json
{
 "data": {
    "name": "John"
 }
}
```

- [Middleware applied to route groups with different configuration](./example/routeGroups)

Milogo middleware, as any other gin middleware, can be applied to different route groups with different configurations.

## 🤝 Contributing

Feel free to create a PR or suggest improvements or ideas.

## 🔗 Contact

- 📧 manueldoncelmartos@gmail.com
