# go-mongodb-crud

A lightweight, generic, type-safe CRUD layer for MongoDB written in Go. Designed to eliminate repetitive boilerplate across your models while keeping your code simple, predictable, and easy to maintain.

This package provides:

* Generic CRUD operations powered by Go generics
* A QueryBuilder DSL for filtering, searching, pagination, and custom queries
* Soft delete & hard delete support
* Reflection-based filtering
* Easy model wiring through `BaseModel[T]`
* Flexible `ListOptions` for building complex list endpoints

---

## Installation

```sh
go get github.com/chokey2nv/go-mongodb-crud
```

---

## Quick Start

### Define your model

Your model must implement the `Identifiable` interface:

```go
type User struct {
    Id    string `bson:"id"`
    Name  string `bson:"name"`
    Email string `bson:"email"`
}

func (u User) GetId() string {
    return u.Id
}
```

---

### Initialize the model store

```go
import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/chokey2nv/go-mongodb-crud/crud"
)

var UserModel *crud.BaseModel[User]

func InitUserModel(db *mongo.Database) {
    UserModel = crud.NewBaseModel[User](db, "users")
}
```

---

## CRUD Operations

### Create

```go
user := User{Id: "123", Name: "John", Email: "john@example.com"}
inserted, err := UserModel.Insert(ctx, user)
```

### Update

```go
updated, err := UserModel.Update(ctx, "123", User{Name: "New Name"})
```

### Get (auto-id routing)

```go
u, err := UserModel.Get(ctx, User{Id: "123"})
```

### Exists

```go
exists, err := UserModel.Exists(ctx, User{Name: "John"})
```

### Count

```go
count, err := UserModel.Count(ctx, User{Email: "john@example.com"})
```

### Hard Delete

```go
err := UserModel.Delete(ctx, "123")
```

### Soft Delete

```go
err := UserModel.Archive(ctx, "123")
```

---

# List & Advanced Querying

The `ListOptions[T]` struct provides a powerful yet simple way to build list queries.

### Example

```go
results, total, err := UserModel.List(ctx, &crud.ListOptions[User]{
    Limit:    20,
    Skip:     0,
    SortBy:   "createdAt",
    SortDesc: true,
    Search:   "john",
    SearchIn: []string{"name", "email"},
    Filter:   User{Email: "john@example.com"},
})
```

---

# Custom Query Builder Logic

```go
opt.CustomQuery = func(q *crud.QueryBuilder) {
    q.Eq("role", "admin")
}
```

---

# Custom Aggregation Pipeline

```go
opt.CustomPipeline = func(p mongo.Pipeline) mongo.Pipeline {
    return append(p, bson.D{{Key: "$sort", Value: bson.D{{"age", 1}}}})
}
```

---

# Architecture Overview

### `BaseModel[T]`

Handles:

* Insert
* Update
* Get
* Delete (soft & hard)
* List / FindMany
* Reflection-based filtering
* Search queries

### `QueryBuilder`

Fluent query builder used internally by `List`, `FindMany`, and `Get`.

### `ListOptions[T]`

Flexible struct for filtering, sorting, searching, IDs, and custom MongoDB pipelines.

---

# Example File Structure

```
cmd/app/main.go
models/user.go
internal/db/mongo.go
go-mongodb-crud/crud/*
```

---

# Contributing

PRs welcome.

---

# License

MIT
