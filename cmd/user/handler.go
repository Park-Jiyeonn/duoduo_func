package main

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"simple_tiktok/cmd/user/dal/db"
	user "simple_tiktok/kitex_gen/user"
	"simple_tiktok/util/errno"
	"simple_tiktok/util/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, request *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// TODO: Your code here...
	//fmt.Println("=============================================")
	//fmt.Println("这里是第三层，应该叫服务层吧。。。，如果看得见代表能走到这里")
	//fmt.Println("=============================================")
	resp = new(user.RegisterResponse)

	users, err := db.QueryUser(ctx, request.Username)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "数据库查询用户存不存在失败"
		return resp, err
	}
	if len(users) > 0 {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户已存在")
	}
	EncodePassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("密码加密失败")
	}
	newUser := &db.User{
		Username: request.Username,
		Password: string(EncodePassword),
	}
	id, err := db.CreateUser(ctx, newUser)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建用户失败")
	}

	token, err := jwt.GetToken(request.Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建jwt失败")
	}

	resp.StatusCode = 0
	resp.StatusMsg = "success"
	resp.UserId = id
	resp.Token = token
	return resp, nil
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, request *user.LoginRequest) (resp *user.LoginResponse, err error) {
	// TODO: Your code here...
	resp = new(user.LoginResponse)
	users, err := db.QueryUser(ctx, request.Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询用户存不存在失败！")
	}

	//密码认证
	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(request.Password))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户名或密码错误")
	}

	//颁发token
	token, err := jwt.GetToken(request.Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("颁发token失败")
	}

	resp.StatusCode = 0
	message := "登录成功"
	Id := int64(users[0].ID)
	resp.StatusMsg = &message
	resp.Token = &token
	resp.UserId = &Id
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, request *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	// TODO: Your code here...
	resp = new(user.UserInfoResponse)
	var message = ""
	resp.StatusMsg = &message

	fmt.Println("===================================")
	fmt.Println(request)
	fmt.Println("===================================")
	//数据库查一查，有没有这个人
	ThisUser, err := db.GetUser(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询失败")
	}
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("没有此用户")
	}
	resp.StatusCode = 0
	message = "请求成功"
	resp.User = &user.UserInfo{
		Id:            int64(ThisUser.ID),
		Name:          ThisUser.Username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return resp, nil
}
