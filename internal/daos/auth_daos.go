package daos

import (
	"context"
)

type AuthRepository interface {
	LoginUser(ctx context.Context, username string) (res User, err error)
	RegisterUser(ctx context.Context, params User) (res uint, err error)
}
