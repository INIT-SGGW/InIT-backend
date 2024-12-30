package main

import (
	"INIT-SGGW/InIT-backend-00.Gateway/initializer"
	"INIT-SGGW/InIT-backend-00.Gateway/service"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type HiResponse struct {
	Body struct {
		Message string `json:"message" example:"I'm Alive!" doc:"Health check"`
	}
}

func main() {
	logger := initializer.CreateLogger()
	fmt.Println("Hello InIT!")
	r := chi.NewRouter()
	r.Use(initializer.New(logger))
	registerServ := service.NewServiceConnection("localhost:8089", logger, time.Second*10)

	standardApiRouter := chi.NewRouter()

	r.Mount("/v1/api", standardApiRouter)

	api := humachi.New(r, huma.DefaultConfig("KN INIT Website API", "1.0.0"))

	huma.Get(api, "/hearthbeat", func(ctx context.Context, input *struct{}) (*HiResponse, error) {
		resp := &HiResponse{}
		resp.Body.Message = "I'm Alive!"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "register-user",
		Method:      http.MethodPost,
		Path:        "/api/v1/register/user",
		Summary:     "Register user",
		Description: "Register user and send confirmation email to provided adress with unique token for account verification",
	}, registerServ.HandleRegisterUserRequest)

	http.ListenAndServe(":3131", r)
}
