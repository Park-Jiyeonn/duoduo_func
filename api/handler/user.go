package handler

import (
	"duoduo_fun/api/model"
	"duoduo_fun/dal/db"
	"duoduo_fun/pkg/errno"
	"duoduo_fun/util/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func UserRegister(c *gin.Context) {
	req := &model.User{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp := new(model.RegisterResponse)
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

func userLogin(req *model.User) (resp *model.LoginResponse, err error) {
	resp = new(model.LoginResponse)
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
	req := &model.User{}
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
