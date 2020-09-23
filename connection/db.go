package connection

import (
	"context"

	"github.com/latihan/model"
)

type SimpleDatabase interface {
	Insert(ctx context.Context, query string) error
	Select(ctx context.Context, query string) ([]model.Student, error)
}
