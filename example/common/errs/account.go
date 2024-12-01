//go:generate errgen -p account.go

package errs

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
	CaptchaValidateFailedCode
)

var (
	UserNotFound          = UserNotFoundCode.error()
	AmountInsufficient    = AmountInsufficientCode.error()
	TokenValidateFailed   = TokenValidateFailedCode.error()
	TokenExpire           = TokenExpireCode.error()
	LoginFailed           = LoginFailedCode.error()
	CaptchaValidateFailed = CaptchaValidateFailedCode.error()
)
