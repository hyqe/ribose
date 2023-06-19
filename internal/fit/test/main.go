// Run this file and try out the example.
//
//	go run internal/fit/test/main.go
//
//	curl -X POST http://localhost:8080/Users/Create -d '{"id":"1", "name":"John Doe", "email":"john.doe@example.com"}' -H "Content-Type: application/json"
//	curl -X POST http://localhost:8080/Users/Get -d '{"id":"1"}' -H "Content-Type: application/json"
package main

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/hyqe/ribose/internal/fit"
	"github.com/hyqe/ribose/internal/fit/status"
)

func main() {
	u := NewUsers()
	service := fit.NewRPC(u)

	app := fiber.New()
	service.MountFiberApp(app)
	app.Listen(":8080")
}

type Users struct {
	sync.RWMutex
	lookup map[string]User
}

func NewUsers() *Users {
	return &Users{
		lookup: map[string]User{},
	}
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserCreateRequest = User
type UserCreateResponse = User

func (u *Users) Create(ctx context.Context, in *UserCreateRequest) (*UserCreateResponse, status.Status) {
	out := *in
	u.RLock()
	defer u.RUnlock()
	u.lookup[in.ID] = out

	return &out, status.OK
}

type UserGetRequest struct {
	ID string `json:"id"`
}
type UserGetResponse = User

func (u *Users) Get(ctx context.Context, in *UserGetRequest) (*UserGetResponse, status.Status) {
	u.RLock()
	defer u.RUnlock()
	user, ok := u.lookup[in.ID]
	if !ok {
		return &user, status.NotFound
	}
	return &user, status.OK
}
