package service

import (
	"errors"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

type UserService struct {
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
}

type UserInfoService struct {
	Token string `form:"token" json:"token"`
}

func (service *UserService) Register() serializer.Response {

	var user model.User
	result := model.DB.Model(&model.User{}).Where("name=?", service.UserName).First(&user)
	if result.RowsAffected == 1 {
		return serializer.Response{
			StatusCode: 1,
			StatusMsg:  "User already exist",
		}
	}
	// atomic.AddInt64(&userIdSequence, 1)
	user.Name = service.UserName
	user.SetPassword(service.Password)
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			StatusCode: 1,
			StatusMsg:  "DataBase save user failed",
		}
	}
	return serializer.Response{
		StatusCode: 0,
		StatusMsg:  "register successfully",
	}
}

func (service *UserService) Login() serializer.UserLoginResponse {
	var user model.User

	if err := model.DB.Where("name=?", service.UserName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.UserLoginResponse{
				Response: serializer.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			}
		}
	}
	// 用户名没错就验证密码
	if user.CheckPassword(service.Password) == false {
		return serializer.UserLoginResponse{
			Response: serializer.Response{StatusCode: 1, StatusMsg: "Password error"},
		}
	}
	// 密码正确，签发token
	token, _ := utils.GenerateToken(user.Id, user.Name)
	return serializer.UserLoginResponse{
		Response: serializer.Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	}

}

func (service *UserInfoService) Info() serializer.UserResponse {
	claim, _ := utils.ParseToken(service.Token)
	username := claim.UserName
	var user model.User
	if err := model.DB.Where("name=?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.UserResponse{
				Response: serializer.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			}
		}
	}
	return serializer.UserResponse{
		Response: serializer.Response{StatusCode: 0},
		User:     user,
	}
}
