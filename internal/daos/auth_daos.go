package daos

import (
	"context"
)

type AuthRepository interface {
	LoginUser(ctx context.Context, params User) (res User, err error)
	CreateUser(ctx context.Context, params User) (res uint, err error)
}
