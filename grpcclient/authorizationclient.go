package grpcclient

import (
	"authentication/authorization"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthorizationGrpcClient struct {
	authorization.AuthorizationClient
	grpc.ClientConn
}

func GetAuthorizationGrpcClient() (*AuthorizationGrpcClient, error) {
	conn, err := grpc.Dial("localhost:81", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	} else {
		client := authorization.NewAuthorizationClient(conn)
		return &AuthorizationGrpcClient{
			AuthorizationClient: client,
			ClientConn:          *conn,
		}, nil
	}
}

func (u *AuthorizationGrpcClient) CreateJWT(id string) (string, error) {
	request := &authorization.Createjwtinput{Userid: id}
	response, err := u.AuthorizationClient.CreateJWT(context.Background(), request)
	if err != nil {
		return "", err
	} else {
		return response.Token, nil
	}
}

func (u *AuthorizationGrpcClient) DeleteJWT(token string) (bool, error) {
	request := &authorization.Deletejwtinput{Token: token}
	response, err := u.AuthorizationClient.DeleteJWT(context.Background(), request)
	if err != nil {
		return false, err
	} else {
		return response.Res, nil
	}
}

func (u *AuthorizationGrpcClient) ValidateJWT(token string) (string, error) {
	request := &authorization.Validatejwtinput{Token: token}
	response, err := u.AuthorizationClient.ValidateJWT(context.Background(), request)
	if err != nil {
		return "", err
	} else {
		return response.Userid, nil
	}
}
