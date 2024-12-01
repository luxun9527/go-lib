package errs

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code codes.Code

func WarpMessage(err error, msg string) error {
	s, ok := status.FromError(err)
	if ok {
		msg = s.Message() + ":" + msg
		return status.New(s.Code(), msg).Err()
	}
	return errors.Wrap(err, msg)
}

func (c Code) error() error {
	return status.Error(codes.Code(c), "")
}

func (c Code) Error(msg string) error {
	return status.Error(codes.Code(c), msg)
}

func (c Code) String() string {
	return ""
}

func (c Code) DtmErrorMsg() string {
	return fmt.Sprintf("=%d=", int32(c))
}
