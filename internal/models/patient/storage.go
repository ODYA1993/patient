package patient

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, person *Person) error
	FindAll(ctx context.Context) ([]*Person, error)
	Update(ctx context.Context, person *Person) error
	Delete(ctx context.Context, guid interface{}) error
}
