package storage

import (
	"context"
)

type AStorage interface {
	Create(ctx context.Context) error
}
