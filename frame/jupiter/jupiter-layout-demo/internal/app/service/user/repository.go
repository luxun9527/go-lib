package user

import (
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/model/db"
)

type Repository interface {
	Get(id int) (user db.User,err error)
}
