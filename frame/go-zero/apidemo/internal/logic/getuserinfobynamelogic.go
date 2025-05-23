package logic

import (
	"context"
	"go-lib/frame/go-zero/apidemo/internal/svc"
	"go-lib/frame/go-zero/apidemo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// User mapped from table <user>
type User struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username  string `gorm:"column:username;not null;comment:用户名" json:"username"`      // 用户名
	Age       int32  `gorm:"column:age;not null;comment:年龄" json:"age"`                 // 年龄
	Fav       string `gorm:"column:fav;not null;comment:爱好" json:"fav"`                 // 爱好
	CreatedAt int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"` // 修改时间
	DeletedAt int64  `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"` // 删除时间
}

func NewGetUserInfoByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByNameLogic {
	return &GetUserInfoByNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByNameLogic) GetUserInfoByName(req *types.GetUserInfoByNameReq) (resp *types.Response, err error) {

	//INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`) VALUES ('',0,'',1692947238,1692947238)

	resp = &types.Response{Message: "test"}
	return
}
