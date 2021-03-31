package module

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service base service
type Service struct {
}

// Error new grpc error
func (s Service) Error(code codes.Code, msg string) error {
	return status.Error(code, msg)
}
