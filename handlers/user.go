package handlers

import (
	"errors"
	"gallery-api/header"
	"gallery-api/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	us models.UserService
}

func NewUserHandler(us models.UserService) *UserHandler {
	return &UserHandler{us}
}

type SignupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *UserHandler) Signup(c *gin.Context) {
	req := new(SignupReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	user := new(models.User)
	user.Email = req.Email
	user.Password = req.Password
	if err := uh.us.Create(user); err != nil {
		Error(c, 500, err)
		return
	}
	c.JSON(201, gin.H{
		"token": user.Token,
		"email": user.Email,
	})
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *UserHandler) Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	user := new(models.User)
	user.Email = req.Email
	user.Password = req.Password
	token, err := uh.us.Login(user)
	if err != nil {
		Error(c, 401, err)
		return
	}
	c.JSON(201, gin.H{
		"token": token,
	})
}

func (uh *UserHandler) Logout(c *gin.Context) {
	token := header.GetToken(c)
	if token == "" {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	err := uh.us.Logout(token)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}
