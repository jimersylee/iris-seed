package app

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/iris-contrib/middleware/prometheus"
	datasource "github.com/jimersylee/iris-seed/datasources"
	"github.com/jimersylee/iris-seed/repositories"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/_examples/mvc/login/web/controllers"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"time"
)

func RunApp() {
	app := initIris()
	initConfig()
	initLog(app)
	initPrometheus(app)
	initDoc(app)
	initRouter(app)
	initDataSource(app)

	_ = app.Run(iris.Addr(":10002"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}

func initConfig() {

}

//初始化数据源
func initDataSource(app *iris.Application) {

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
		ctx.View("shared/error.html")
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

	// Prepare our repositories and services.
	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}
	repo := repositories.NewUserRepository(db)
	userService := services.NewUserService(repo)

	// "/user" based mvc application.
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))

	// http://localhost:8080/noexist
	// and all controller's methods like
	// http://localhost:8080/users/1
	// http://localhost:8080/user/register
	// http://localhost:8080/user/login
	// http://localhost:8080/user/me
	// http://localhost:8080/user/logout

}

//初始化文档
func initDoc(app *iris.Application) {
	//api文档自动生成开始
	yaag.Init(&yaag.Config{On: true, DocTitle: "iris-seed", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Stage": ""}})
	app.Use(irisyaag.New())
	//api文档自动生成结束

}
