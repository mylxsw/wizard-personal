package controller

import "github.com/pkg/errors"

type OperationResponse struct {
	Message string `json:"message"`
}

func Success() (*OperationResponse, error) {
	return &OperationResponse{Message: "操作成功"}, nil
}

func Error(err error, message string) (*OperationResponse, error) {
	return nil, errors.Wrap(err, message)
}
