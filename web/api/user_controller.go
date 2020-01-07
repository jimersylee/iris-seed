package api

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/jimersylee/iris-seed/commons/api_token"
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jimersylee/iris-seed/entities"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

// UserController是我们的/用户控制器。
// UserController负责处理以下请求：
// GET              /user/register
// POST             /user/register
// GET                 /user/login
// POST             /user/login
// GET                 /user/me
//所有HTTP方法 /user/logout
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
func (c *UserController) PostRegister() *commons.WebApiResult {
	//从表单中获取名字，用户名和密码
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)
	//创建新用户，密码将由服务进行哈希处理
	u, err := services.UserService.Create(password, datamodels.User{
		Name: username,
	})
	if err != nil {
		return commons.JsonErrorMsg(err.Error())
	}
	return commons.JsonData(u)

}

// PostLogin handles
// PostLogin处理POST: http://localhost:17001/user/register.
func (c *UserController) PostLogin() *commons.WebApiResult {
	var (
	//username = c.Ctx.FormValue("username")
	//password = c.Ctx.FormValue("password")

	)
	ee := &entities.LoginDTO{}
	err := c.Ctx.ReadJSON(ee)
	if err != nil {
		commons.JsonErrorMsg("解析错误")
	}
	user, found := services.UserService.GetByUsernameAndPassword(ee.Username, ee.Password)
	logrus.Info("username:"+ee.Username)
	if !found {
		return commons.JsonErrorCode(111, "账号未找到")
	}
	token := services.UserTokenService.UpdateToken(user.ID)
	api_token.SetApiCurrentUser(token, user.ID)
	return commons.JsonData(token)
}

// GetMe 处理P GET: http://localhost:17001/user/me.
func (c *UserController) GetMe() *commons.WebApiResult {
	if !c.isLoggedIn() {
		//如果没有登录，则将用户重定向到登录页面。
		return commons.JsonErrorMsg("未登录")
	}
	u := services.UserService.GetByID(c.getCurrentUserID())
	if u == nil {
		//如果session存在但由于某种原因用户不存在于“数据库”中
		//然后注销并重新执行该函数，它会将客户端重定向到
		// /user/login页面。
		return commons.JsonErrorCode(404, "未找到用户")
	}
	return commons.JsonData(u)

}

// AnyLogout处理 All/AnyHTTP 方法：http://localhost:17001/user/logout
func (c *UserController) AnyLogout() *commons.WebApiResult {
	c.logout()
	return commons.JsonSuccess()
}
