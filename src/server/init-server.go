package main

import (
	"INIT-SGGW/InIT-backend/initializer"
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type HiResponse struct {
	Body struct {
		Message string `json:"message" example:"Hello KN InIT!" doc:"Greeting message"`
	}
}

func main() {
	logger := initializer.CreateLogger()
	fmt.Println("Hello InIT!")
	r := chi.NewRouter()
	r.Use(initializer.New(logger))

	api := humachi.New(r, huma.DefaultConfig("KN INIT Website API", "1.0.0"))

	huma.Get(api, "/hi", func(ctx context.Context, input *struct{}) (*HiResponse, error) {
		resp := &HiResponse{}
		resp.Body.Message = "Hello KN InIT!"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "welcome-user",
		Method:      http.MethodGet,
		Path:        "/hi/{name}",
		Summary:     "Great user",
		Description: "Great user with his name",
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*HiResponse, error) {
		resp := &HiResponse{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	http.ListenAndServe(":3131", r)
}
