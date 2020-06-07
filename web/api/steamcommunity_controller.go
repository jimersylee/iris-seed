package api

import (
	"fmt"
	"github.com/jimersylee/iris-seed/commons/response"
	"github.com/jimersylee/iris-seed/services"
	"github.com/kataras/iris"
)

type SteamCommunityController struct {
	Ctx iris.Context
}

func (c *SteamCommunityController) All() *response.WebApiRes {

	uri := c.Ctx.Request().RequestURI
	fmt.Println(uri)
	uriNeed := uri[len("/api/steamcommunity/") : len(uri)-1]
	route := c.Ctx.GetCurrentRoute()
	fmt.Println(route)
	fmt.Println(uriNeed)
	ip := "182.255.44.92"
	webUrl := "http://api.steampowered.com/" + uriNeed
	responseStr := services.IpService.Proxy(ip, webUrl)
	return response.JsonData(responseStr)

}
