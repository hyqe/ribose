package status

import (
	"fmt"

	"github.com/hyqe/ribose/internal/fit/codes"
	"github.com/lib/pq"
)

// convenience vars for common response.
var (
	OK        = Status{Code: codes.OK}
	Created   = Status{Code: codes.Created}
	NoContent = Status{Code: codes.NoContent}
	NotFound  = Status{Code: codes.NotFound}
)

type Status struct {
	codes.Code
	Message string
}

func (s Status) Error() string {
	return fmt.Sprintf("%v: %v", s.Code, s.Message)
}

func New(code codes.Code, v any) Status {
	return Status{
		Code:    code,
		Message: fmt.Sprint(v),
	}
}

func Newf(code codes.Code, format string, v ...any) Status {
	return Status{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

// Pg converts an pq.Error into a Status.
func Pg(e *pq.Error) Status {
	return Status{
		Code:    codes.Pg(e.Code),
		Message: e.Error(),
	}
}
