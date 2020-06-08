package services

import (
	"github.com/jimersylee/iris-seed/commons"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"strings"
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

func (p *ProxyServiceImpl) Proxy(ctx iris.Context) {
	uri := ctx.Request().RequestURI

	//判断是代理api还是社区
	uriNeed := ""
	webUrl := ""
	if strings.Index(uri, "steamapi") > 0 {
		uriNeed = uri[len("/api/steamapi/"):]
		webUrl = "http://api.steampowered.com/" + uriNeed
	} else {
		uriNeed = uri[len("/api/steamcommunity/"):]
		webUrl = "http://steamcommunity.com/" + uriNeed
	}
	//todo：找出能用的ip
	ipModel := IpService.FindOne(commons.NewSqlCnd().Where("status=1").Asc("request_times"))
	logrus.Infof("ipModel:%s,uri:%s,uriNeed:%s,webUrl:%s", ipModel, uri, uriNeed, webUrl)
	if ipModel == nil {
		logrus.Error("can't find available ip")
		ctx.ResponseWriter().WriteHeader(500)
		return
	}
	ip := ipModel.Ip
	responseStr, statusCode := IpService.Proxy(ip, webUrl)
	logrus.Infof("steam返回数据：%s", responseStr)
	ctx.ResponseWriter().WriteHeader(statusCode)
	ctx.WriteString(responseStr)

}
