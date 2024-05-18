package error

// 使用这个命令，可以code生成对应的错误
//
//go:generate errgen -p common.go
const (
	// UserNotFoundCodeCode  用户不存在
	UserNotFoundCode = AccountCodeInit + iota + 1
	// AmountInsufficientCode 用户余额不足
	AmountInsufficientCode
	// TokenValidateFailedCode token验证失败
	TokenValidateFailedCode
	// TokenExpireCode Token到期
	TokenExpireCode
	// LoginFailedCode 登录账户密码验证失败
	LoginFailedCode
)

var (
	UserNotFound        = UserNotFoundCode.Error("")
	AmountInsufficient  = AmountInsufficientCode.Error("")
	TokenValidateFailed = TokenValidateFailedCode.Error("")
	TokenExpire         = TokenExpireCode.Error("")
	LoginFailed         = LoginFailedCode.Error("")
)
