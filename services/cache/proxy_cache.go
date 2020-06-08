package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/jimersylee/iris-seed/commons/redis_manager"
	"github.com/sirupsen/logrus"
)

var (
	REDIS_KEY_429     = "ip:429:"
	REDIS_KEY_IP_POOL = "ip:pool"
)

var ProxyCache *proxyCache

func newProxyCache() *proxyCache {
	client := redis_manager.GetClient()
	return &proxyCache{redisClient: client}
}

type proxyCache struct {
	redisClient *redis.Client
}

func (this *proxyCache) GetUserIdByToken(token string) int64 {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	if len(token) == 0 {
		return 0
	}
	tokeneee := this.redisClient.Get(token)
	if tokeneee == nil {
		return 0
	}
	userId, err := tokeneee.Int64()
	if err != nil {
		return 0
	}

	return userId

}

func (this *proxyCache) Delete(token string) {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	logrus.Info("token:" + token)
	this.redisClient.Del(token)
}
func (this *proxyCache) incr429TimesByIp(ip string) {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	this.redisClient.Incr("")
}

func (this *proxyCache) IpPoolAdd(ip string) {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	this.redisClient.SAdd(REDIS_KEY_IP_POOL, ip)
}
func (this *proxyCache) IpPoolDel(ip string) {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	this.redisClient.SRem(REDIS_KEY_IP_POOL, ip)
}

func (this *proxyCache) IpPoolGetAll() []string {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	members := this.redisClient.SMembers(REDIS_KEY_IP_POOL).Val()
	return members
}
func (this *proxyCache) GetRedisClient() *redis.Client {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	return this.redisClient
}
