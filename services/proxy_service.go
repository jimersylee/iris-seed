package services

import (
	"fmt"
	"github.com/jimersylee/iris-seed/commons"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

type ProxyInterface interface {
	//注册ip
	ProxyCommunity(ctx iris.Context) bool
	ProxyApi(ctx iris.Context)
}

var ProxyService = NewProxyService()

func NewProxyService() *ProxyServiceImpl {
	return &ProxyServiceImpl{}
}

type ProxyServiceImpl struct {
}

func (p *ProxyServiceImpl) ProxyCommunity(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	uriNeed := uri[len("/api/steamcommunity/") : len(uri)-1]
	//todo：找出能用的ip
	ipModel := IpService.FindOne(commons.NewSqlCnd().Where("status=1").Asc("request_times"))
	logrus.Info("ipModel:", ipModel)
	if ipModel == nil {
		logrus.Error("can't find available ip")
		ctx.ResponseWriter().WriteHeader(500)
		return
	}
	ip := ipModel.Ip
	webUrl := "http://steamcommunity.com/" + uriNeed
	responseStr, statusCode := IpService.Proxy(ip, webUrl)
	fmt.Println(responseStr)
	ctx.ResponseWriter().WriteHeader(statusCode)
	ctx.WriteString(responseStr)
}

func (p *ProxyServiceImpl) ProxyApi(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	uriNeed := uri[len("/api/steamapi/") : len(uri)-1]
	//todo：找出能用的ip
	ipModel := IpService.FindOne(commons.NewSqlCnd().Where("status=1").Asc("request_times"))
	logrus.Info("ipModel:", ipModel)
	if ipModel == nil {
		logrus.Error("can't find available ip")
		ctx.ResponseWriter().WriteHeader(500)
		return
	}
	ip := ipModel.Ip
	webUrl := "http://api.steampowered.com/" + uriNeed
	responseStr, statusCode := IpService.Proxy(ip, webUrl)
	fmt.Println(responseStr)
	ctx.ResponseWriter().WriteHeader(statusCode)
	ctx.WriteString(responseStr)
}
