package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/jimersylee/iris-seed/commons/redis_manager"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	REDIS_KEY_IP_2_HASH = "ip:hash:"
	REDIS_KEY_IP_POOL   = "ip:pool"
	//错误次数出现的次数阈值
	ERROR_TIMES_THRESHOLD = 10
	//错误次数统计间隔
	ERROR_TIMES_INTERVAL_SECONDS = 30
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
func (this *proxyCache) incrErrorTimesByIp(ip string, httpStatus int) {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	this.redisClient.HIncrBy(REDIS_KEY_IP_2_HASH+ip, strconv.Itoa(httpStatus), 1)
}
func (this *proxyCache) CalcIpNeedToBeBanned(ip string) bool {
	if UserTokenCache == nil {
		this = newProxyCache()
	}
	result, err := this.redisClient.HGetAll(REDIS_KEY_IP_2_HASH + ip).Result()
	if err != nil {
		return false
	}
	totalTimes := 0
	for _, v := range result {
		times, _ := strconv.Atoi(v)
		totalTimes += times
	}
	if totalTimes > ERROR_TIMES_THRESHOLD {
		logrus.Infof("CalcIpNeedToBeBanned,need to ban，ip:%s,error times:%d", ip, totalTimes)
		return true
	}
	return false
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
