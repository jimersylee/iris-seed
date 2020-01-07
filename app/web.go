package app

import (
	"flag"
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/iris-contrib/middleware/prometheus"
	"github.com/jimersylee/iris-seed/commons/api_token"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/commons/web_session"
	"github.com/jimersylee/iris-seed/config"
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jimersylee/iris-seed/web/api"
	"github.com/jimersylee/iris-seed/web/frontend"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"time"
)

func RunApp() {
	initConfig()
	app := initIris()
	initLog(app)
	initPrometheus(app)
	initDoc(app)
	initRouter(app)
	initDataSource(app)
	//初始化web session管理
	web_session.InitSessionManager()
	//初始化api token 管理
	api_token.InitTokenManager()

	_ = app.Run(iris.Addr(":"+config.Conf.Port), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}

func initConfig() {
	var configFile = flag.String("config", "./iris-seed.yaml", "配置文件路径")
	config.InitConfig(*configFile)
}

//初始化数据源
func initDataSource(app *iris.Application) {
	// 连接数据库
	db.OpenDB(&db.DBConfiguration{
		Dialect:        "mysql",
		Url:            config.Conf.MySqlUrl,
		MaxIdle:        5,
		MaxActive:      20,
		EnableLogModel: config.Conf.ShowSql,
		Models:         datamodels.Models,
	})
}

//初始化iris框架
func initIris() *iris.Application {
	app := iris.New()
	app.Use(recover.New())

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		_, _ = ctx.Writef("Not Found")
	})
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		_ = ctx.View("shared/error.html")
	})

	// Load the template files.
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.StaticWeb("/public", "./web/public")

	return app
}

//初始化日志
func initLog(app *iris.Application) {
	app.Logger().SetLevel("debug")
	app.Use(logger.New())
}

//初始化监控
//init monitor
func initPrometheus(app *iris.Application) {
	//集成prometheus监控开始,访问/metrics
	m := prometheus.New("go-bbs", 300, 1200, 5000)
	app.Use(m.ServeHTTP)
}

//初始化路由
func initRouter(app *iris.Application) {
	app.Handle("GET", "/", func(ctx iris.Context) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		_, _ = ctx.Writef("Slept for %d milliseconds", sleep)
	})
	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.WriteString("pong")
	})
	app.Get("/hello", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{"message": "hello"})
	})
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	// "/user" based mvc application.
	user := mvc.New(app.Party("/user"))
	user.Handle(new(frontend.UserController))
	apiUser := mvc.New(app.Party("/api/user"))
	apiUser.Handle(new(api.ApiUserController))

	// http://localhost:17001/noexist
	// and all controller's methods like
	// http://localhost:17001/users/1
	// http://localhost:17001/user/register
	// http://localhost:17001/user/login
	// http://localhost:17001/user/me
	// http://localhost:17001/user/logout

}

//初始化文档
func initDoc(app *iris.Application) {
	//api文档自动生成开始
	yaag.Init(&yaag.Config{On: true, DocTitle: "iris-seed", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Stage": ""}})
	app.Use(irisyaag.New())
	//api文档自动生成结束

}
