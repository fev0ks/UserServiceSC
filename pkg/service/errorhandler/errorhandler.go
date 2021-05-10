package errorhandler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func LogMsg(msg string) {
	if msg != "" {
		log.Printf("msg: %v", msg)
	}
}

func NewStatusError(code codes.Code, errorMsg string) error {
	LogMsg(errorMsg)
	return status.Error(code, errorMsg)
}

func NewInvalidArgumentError(msg string) error {
	return NewStatusError(codes.InvalidArgument, msg)
}

func NewNotFoundError(msg string) error {
	return NewStatusError(codes.NotFound, msg)
}

func NewInternalError(msg string) error {
	return NewStatusError(codes.Internal, msg)
}
