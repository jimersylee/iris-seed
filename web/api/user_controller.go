package api

import (
	"github.com/jimersylee/go-steam-proxy/commons"
	"github.com/jimersylee/go-steam-proxy/commons/api_token"
	"github.com/jimersylee/go-steam-proxy/commons/response"
	"github.com/jimersylee/go-steam-proxy/entities"
	"github.com/jimersylee/go-steam-proxy/models"
	"github.com/jimersylee/go-steam-proxy/services"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

// UserController是我们的/用户控制器。
type UserController struct {
	//每个请求都由Iris自动绑定上下文，
	//记住，每次传入请求时，iris每次都会创建一个新的UserController，
	//所以所有字段都是默认的请求范围，只能设置依赖注入
	//自定义字段，如服务，对所有请求都是相同的（静态绑定）
	//和依赖于当前上下文的会话（动态绑定）。
	Ctx iris.Context
	//我们的UserService，它是一个接口
	//从主应用程序绑定。
}

func (c *UserController) getCurrentUserID() int64 {
	userId := api_token.GetApiCurrentUser(c.Ctx)
	return userId
}
func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}
func (c *UserController) logout() {
	api_token.DelApiCurrentUser(c.Ctx)
}

// PostRegister 处理 POST: http://localhost:17001/user/register.
func (c *UserController) PostRegister() *response.WebApiRes {
	loginDTO := &entities.LoginDTO{}
	err := c.Ctx.ReadJSON(loginDTO)
	if err != nil {
		response.JsonErrorCode(commons.ErrorCodeParse)
	}
	//创建新用户，密码将由服务进行哈希处理
	err = services.UserService.Create(loginDTO.Password, models.User{
		Name: loginDTO.Username,
	})
	if err != nil {
		return response.JsonErrorMsg(err.Error())
	}
	user, b := services.UserService.GetUserByUsernameAndPassword(loginDTO.Username, loginDTO.Password)
	if b {
		return response.JsonData(user)
	}

	return response.JsonErrorCode(commons.ErrorCodeRegisterFailed)

}

// PostLogin handles
// PostLogin处理POST: http://localhost:17001/user/register.
func (c *UserController) PostLogin() *response.WebApiRes {
	loginDTO := &entities.LoginDTO{}
	err := c.Ctx.ReadJSON(loginDTO)
	if err != nil {
		response.JsonErrorCode(commons.ErrorCodeParse)
	}
	user, found := services.UserService.GetUserByUsernameAndPassword(loginDTO.Username, loginDTO.Password)
	logrus.Info("username:" + loginDTO.Username)
	if !found {
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}

	token := services.UserTokenService.UpdateToken(user.ID)
	api_token.SetApiCurrentUser(token, user.ID)
	return response.JsonData(token)
}

// GetMe 处理P GET: http://localhost:17001/user/me.
func (c *UserController) GetMe() *response.WebApiRes {
	if !c.isLoggedIn() {
		//如果没有登录，则将用户重定向到登录页面。
		return response.JsonErrorMsg("未登录")
	}
	u := services.UserService.GetByID(c.getCurrentUserID())
	if u == nil {
		//如果session存在但由于某种原因用户不存在于“数据库”中
		//然后注销并重新执行该函数，它会将客户端重定向到
		// /user/login页面。
		return response.JsonErrorCode(commons.ErrorCodeNotFound)
	}
	return response.JsonData(u)

}

// GetMe 处理P GET: http://localhost:17001/user/me.
func (c *UserController) GetJimersylee() *response.WebApiRes {
	user, _ := services.UserService.GetByUsername("jimersylee")
	return response.JsonData(user)

}

// AnyLogout处理 All/AnyHTTP 方法：http://localhost:17001/user/logout
func (c *UserController) AnyLogout() *response.WebApiRes {
	c.logout()
	return response.JsonSuccess()
}
