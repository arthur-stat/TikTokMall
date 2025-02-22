package service

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// getUserIDByToken 根据token寻找用户id（redis -> mysql）
func getUserIDByToken(token string) string {

}
