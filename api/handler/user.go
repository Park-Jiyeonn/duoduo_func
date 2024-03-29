package handler

import (
	"context"
	api_model "duoduo_fun/api/model"
	"duoduo_fun/dal/db"
	"duoduo_fun/dal/db/model"
	"duoduo_fun/dal/redis"
	"duoduo_fun/pkg/errno"
	"duoduo_fun/util/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func UserRegister(c *gin.Context) {
	req := &api_model.User{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp := new(api_model.RegisterResponse)
	user, err := db.GetUserByName(req.Username)
	if user != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "用户已存在"
	}
	EncodePassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "密码加密失败"
	}
	id, err := db.CreateUser(req.Username, string(EncodePassword))
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "创建用户失败"
	}

	token, err := jwt.GenToken(id, req.Username)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "创建jwt失败"
	}

	resp.StatusCode = 0
	resp.StatusMsg = "success"
	resp.UserId = id
	resp.Token = token
}

func userLogin(req *api_model.User) (resp *api_model.LoginResponse, err error) {
	resp = new(api_model.LoginResponse)
	user, err := db.GetUserByName(req.Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询用户存不存在失败！")
	}

	//密码认证
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户名或密码错误")
	}

	//颁发token
	token, err := jwt.GenToken(int(user.ID), user.Name)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("颁发token失败")
	}

	resp.StatusCode = 200
	message := "登录成功"
	resp.StatusMsg = message
	resp.Token = token
	resp.UserId = int(user.ID)
	return resp, nil
}
func UserLogin(c *gin.Context) {
	req := &api_model.User{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := userLogin(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func getUserInfo(ctx context.Context, req *api_model.UserInfoRequest) (resp *api_model.UserInfoResponse, err error) {
	var user *model.User
	if redis.UserIsExists(ctx, req.UserId) != 0 {
		user, err = redis.GetUserInfo(ctx, req.UserId)
		resp.StatusCode = 1
		if err != nil {
			return resp, err
		}
	} else {
		user, err = db.GetUserById(ctx, req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("数据库查询失败")
		}

		//	数据库查询成功后，需要将信息写入缓存
		redis.SetUserInfo(ctx, user)
	}

	resp.StatusCode = 200
	resp.StatusMsg = "请求成功"
	resp.User = api_model.UserInfo{
		Id:   int(user.ID),
		Name: user.Name,
	}
	return resp, nil
}
func GetUserInfo(c *gin.Context) {
	req := &api_model.UserInfoRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := getUserInfo(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
