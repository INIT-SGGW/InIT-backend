package main

import (
	"INIT-SGGW/InIT-backend-00.Gateway/initializer"
	"INIT-SGGW/InIT-backend-00.Gateway/model"
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "INIT-SGGW/InIT-backend-00.Gateway/proto-messages/protogen/register"

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
	}, func(ctx context.Context, input *model.RegisterUserRequest) (*model.RegisterUserResponse, error) {
		resp := &model.RegisterUserResponse{}
		resp.Body.Status = "created"
		return resp, nil
	})

	conn, err := grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_ = pb.NewRegisterUserSenderClient(conn)
	grpc.WithTransportCredentials(insecure.NewCredentials())

	// for i := 0; i < 10; i++ {
	// 	resp, err := client.SendRegisterUserRequestSession(context.Background(), &pb.RegisterUserRequestRPC{RequestSessionId: int32(i), FirstName: "Adrian", IsAggrementFulfielled: i%2 == 0})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(resp)
	// }

	http.ListenAndServe(":3131", r)
}
