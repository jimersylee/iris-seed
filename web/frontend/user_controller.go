package frontend

import (
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
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
	Service services.UserService
	//Session，使用来自main.go的依赖注入绑定
	Session *sessions.Session
}

const userIDKey = "UserID"

func (c *UserController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}
func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}
func (c *UserController) logout() {
	c.Session.Destroy()
}

var registerStaticView = mvc.View{
	Name: "user/register.html",
	Data: iris.Map{"Title": "User Registration"},
}

// GetRegister 处理 GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}
	return registerStaticView
}

// PostRegister 处理 POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	//从表单中获取名字，用户名和密码
	var (
		firstname = c.Ctx.FormValue("firstname")
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)
	//创建新用户，密码将由服务进行哈希处理
	u, err := c.Service.Create(password, datamodels.User{
		Username:  username,
		Firstname: firstname,
	})
	//将用户的id设置为此会话，即使err！= nil，
	//零id无关紧要因为.getCurrentUserID()检查它。
	//如果错误！= nil那么它将被显示，见下面的mvc.Response.Err：err
	c.Session.Set(userIDKey, u.ID)
	return mvc.Response{
		//如果不是nil，则会显示此错误
		Err: err,
		//从定向 /user/me.
		Path: "/user/me",
		//当从POST重定向到GET请求时，您应该使用此HTTP状态代码，
		//但是如果你有一些（复杂的）选择
		//在线搜索甚至是HTTP RFC。
		//状态“查看其他”RFC 7231，但虹膜可以自动修复它
		//但很高兴知道你可以设置自定义代码;
		//代码：303，
	}
}

var loginStaticView = mvc.View{
	Name: "user/login.html",
	Data: iris.Map{"Title": "User Login"},
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.logout()
	}
	return loginStaticView
}

// PostLogin handles
// PostLogin处理POST: http://localhost:8080/user/register.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)
	u, found := c.Service.GetByUsernameAndPassword(username, password)
	if !found {
		return mvc.Response{
			Path: "/user/register",
		}
	}
	c.Session.Set(userIDKey, u.ID)
	return mvc.Response{
		Path: "/user/me",
	}
}

// GetMe 处理P GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		//如果没有登录，则将用户重定向到登录页面。
		return mvc.Response{Path: "/user/login"}
	}
	u, found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		//如果session存在但由于某种原因用户不存在于“数据库”中
		//然后注销并重新执行该函数，它会将客户端重定向到
		// /user/login页面。
		c.logout()
		return c.GetMe()
	}
	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User":  u,
		},
	}
}

// AnyLogout处理 All/AnyHTTP 方法：http://localhost:8080/user/logout
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}
	c.Ctx.Redirect("/user/login")
}
