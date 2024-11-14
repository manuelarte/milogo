[![Go](https://github.com/manuelarte/milogo/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/milogo/actions/workflows/go.yml)
![coverage](https://raw.githubusercontent.com/manuelarte/milogo/badges/.badges/main/coverage.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/milogo)](https://goreportcard.com/report/github.com/manuelarte/milogo)
![version](https://img.shields.io/github/v/release/manuelarte/milogo)
# Milogo
Rest Partial Response (aka Field Selection) Pattern middleware for [Gin](https://gin-gonic.com/). This gin plugin allows you to select a subset of fields to be returned from your endpoints.
e.g.
```
> /products/1
{
 "code": "1",
 "price": "200",
 "description": "Very nice product",
 "manufacturedBy": "company"
}
```
```
> /products/1?fields=code,price
{
 "code": "1",
 "price": "200",
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
```
/users/1?fields=name,surname
```

## ✨ Features

- [Support for multiple fields filtering](./examples/simple). 

```
> /users/manuel?fields=name,surname
{
 "name": "Manuel",
 "surname": "Example"
}
```

- [Support for arrays](./examples/simple-array)

```
> /users?fields=name
[
  {
    "name": "Manuel"
  }
]
```

- [Support for nested jsons](./examples/nested). Check [example]

```
> /users/manuel?fields=name,surname,address(street,number)
{
 "name": "Manuel",
 "surname": "Example",
 "address": {
   "street": "mystreet",
   "zipcode": "myzipcode"
 }
}
```

- [Support for json wrapped](./examples/wrapped). 
```
> /users/manuel?fields=name
{
 "data": {
    "name": "Manuel"
 }
}
```

## 🤝 Contributing

Feel free to create a PR or suggest improvements or ideas.

## 🔗 Contact

- 📧 manueldoncelmartos@gmail.com
