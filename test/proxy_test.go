package test

import (
	"fmt"
	"github.com/jimersylee/go-steam-proxy/app"
	"github.com/jimersylee/go-steam-proxy/services"
	"github.com/jimersylee/go-steam-proxy/services/cache"
	"testing"
	"time"
)

func TestIpPool(t *testing.T) {
	app.RunApp()
	cache.ProxyCache.IpPoolAdd("127.0.0.1")
	all := cache.ProxyCache.IpPoolGetAll()
	for _, v := range all {
		if v == "127.0.0.1" {
			t.Log("getAll ok")
		}
	}
	cache.ProxyCache.IpPoolDel("127.0.0.1")
	all = cache.ProxyCache.IpPoolGetAll()
	if len(all) > 0 {
		t.Error("del error")
	}
	t.Log("del ok")
}
func TestCheckIp(t *testing.T) {
	app.RunApp()
	services.ProxyService.CheckIpAlive()
}

func TestChangeIp(t *testing.T) {
	app.RunApp()
	services.ProxyService.ChangeIp("38.21.249.98")
}

func TestCheckIpStatus(t *testing.T) {
	app.RunApp()
	services.ProxyService.CheckIpStatus()
}

func TestBuildTestData(t *testing.T) {
	app.RunApp()
	cache.ProxyCache.GetRedisClient().HIncrBy(cache.REDIS_KEY_IP_2_HASH+"127.0.0.1", "500", 3)
	cache.ProxyCache.GetRedisClient().HIncrBy(cache.REDIS_KEY_IP_2_HASH+"127.0.0.1", "429", 8)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
}
