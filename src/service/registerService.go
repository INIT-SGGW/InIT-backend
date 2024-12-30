package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"INIT-SGGW/InIT-backend-00.Gateway/model"
	pb "INIT-SGGW/InIT-backend-00.Gateway/proto-messages/protogen/register"
)

type RegisterServiceContext struct {
	client  pb.RegisterUserSenderClient
	timeout time.Duration
	logger  *zap.Logger
}

type RegisterServiceConnection interface {
	HandleRegisterUserRequest(ctx context.Context, input *model.RegisterUserRequest) (*model.RegisterUserResponse, error)
	SendRequestRPC(ctx context.Context, req *pb.RegisterUserRequestRPC) *pb.RegisterUserResponseRPC
}

func NewServiceConnection(listAddr string, logger *zap.Logger, timeout time.Duration) *RegisterServiceContext {
	defer logger.Sync()

	conn, err := grpc.NewClient(listAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	ctx := &RegisterServiceContext{
		client:  pb.NewRegisterUserSenderClient(conn),
		logger:  logger,
		timeout: timeout,
	}

	return ctx
}

func (rs RegisterServiceContext) HandleRegisterUserRequest(ctx context.Context, input *model.RegisterUserRequest) (*model.RegisterUserResponse, error) {
	defer rs.logger.Sync()

	clientCtx, cancel := context.WithTimeout(ctx, rs.timeout)
	defer cancel()

	serviceReq := pb.RegisterUserRequestRPC{
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
		Email:     input.Body.Email,
		Password:  input.Body.Password,
		DateOfBirth: &date.Date{
			Year:  int32(input.Body.DateOfBirth.Year()),
			Month: int32(input.Body.DateOfBirth.Month()),
			Day:   int32(input.Body.DateOfBirth.Day()),
		},
		IsAggrementFulfielled: input.Body.IsAggrementFulfielled,
		PrivilageLevel:        "User",
	}
	respRPC := rs.SendRequestRPC(clientCtx, &serviceReq)
	resp := model.RegisterUserResponse{}
	resp.Body.Status = respRPC.Status.Status

	if respRPC.Status.Status == "Error-Internal" {
		resp.Body.Error = *respRPC.Status.Errors
		resp.Status = http.StatusInternalServerError
		return &resp, fmt.Errorf("[Internal Error] %v", respRPC.Status.Errors)
	}

	if respRPC.Status.Status == "BadRequest" {
		resp.Body.Error = *respRPC.Status.Errors
		resp.Status = http.StatusBadRequest
		return &resp, fmt.Errorf("[BadRequest] %v", respRPC.Status.Errors)
	}
	resp.Body.Status = respRPC.Status.Status
	resp.Status = http.StatusCreated

	rs.logger.Info("Sucesfully created user",
		zap.Any("resp", resp))

	return &resp, nil

}

func (rs RegisterServiceContext) SendRequestRPC(ctx context.Context, req *pb.RegisterUserRequestRPC) *pb.RegisterUserResponseRPC {
	defer rs.logger.Sync()
	resp, err := rs.client.SendRegisterUserRequestSession(context.Background(), req)

	if err != nil {
		rs.logger.Error("Error sending request to register service",
			zap.Error(err))

		return &pb.RegisterUserResponseRPC{Status: &pb.RPCStatusMessage{
			Status: "Error-Internal",
		}}
	}
	fmt.Println(resp)
	rs.logger.Info("Sucesfully send request",
		zap.Any("resp", resp))

	return resp
}
