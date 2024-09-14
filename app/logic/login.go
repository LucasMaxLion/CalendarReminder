package logic

import (
	"CalendarReminder/app/model"
	"CalendarReminder/app/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// User 用户登录，一般不会把数据库直接暴露在外面，
type User struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

func Login(c *gin.Context) {
	var user User
	// 绑定参数，前端返回的数据会通过绑定参数，来让我们后端获得数据
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(),
		})
	}
	fmt.Printf("user:%+v\n", user)

	//使用name查询数据
	ret := model.GetUser(user.Name)
	//这里的密码要加密的，但我为了方便就没有加密
	if ret.Id < 1 || ret.Password != user.Password {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "账号密码错误！",
		})
		return
	}
	c.SetCookie("uid", ret.Uid, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
	})
	return
}

type CUser struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
}

func CreateUser(context *gin.Context) {
	var user CUser

	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	fmt.Printf("user:%+v", user)

	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: "参数错误",
		})
		return
	}
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同!",
		})
		return
	}

	if oldUser := model.GetUser(user.Name); oldUser.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在！",
		})
		return
	}

	nameLen := len(user.Name)
	password := len(user.Password)
	if nameLen > 16 || nameLen < 8 || password > 16 || password < 8 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "账号或密码大于8小于16",
		})
		return
	}

	// 密码不能是纯数字
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.Password) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字",
		})
		return
	}

	newUser := model.User{
		Name:        user.Name,
		Password:    user.Password,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Uid:         strconv.Itoa(int(tools.GetUid())),
	}

	if err := model.CreateUser(&newUser); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10007,
			Message: "新用户创建失败!",
		})
		return
	}
	//返回添加成功
	context.JSON(http.StatusOK, tools.ECode{Code: 0})
}
