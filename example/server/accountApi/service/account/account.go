package account

var UserService = &userService{}

type userService struct{}

func (*userService) GetAccountInfo() error {

	//accountInfo, err := rpcClient.AccountClient.GetAccountInfo()
	//if err != nil {
	//	zlog.Errorf("GetAccountInfo err: %v", err)
	//	return err
	//}
	return nil
}
