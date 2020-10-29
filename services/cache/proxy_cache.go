package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/jimersylee/go-steam-proxy/commons/redis_manager"
	"github.com/jimersylee/go-steam-proxy/entities"
	"github.com/sirupsen/logrus"
	"net/http"
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

func (this *proxyCache) IncrHttpStatusTimesByIpAndStatus(ip string, httpStatus int) {
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
	for httpStatus, times := range result {
		times, _ := strconv.Atoi(times)
		httpStatusInt, _ := strconv.Atoi(httpStatus)
		//不等于200的计入失败总次数
		if httpStatusInt != http.StatusOK {
			totalTimes += times
		}
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

//获取统计数据
func (this *proxyCache) GetStatistic() *entities.StatisticDTO {
	allIps := this.IpPoolGetAll()
	statisticDTO := new(entities.StatisticDTO)
	for _, ip := range allIps {
		logrus.Info("GetStatistic:" + REDIS_KEY_IP_2_HASH + ip)
		result, err := this.redisClient.HGetAll(REDIS_KEY_IP_2_HASH + ip).Result()
		if err != nil {
			return nil
		}
		oneIp := new(entities.OneIp)
		for httpStatus, times := range result {
			times, _ := strconv.Atoi(times)
			httpStatusInt, _ := strconv.Atoi(httpStatus)
			oneCode := new(entities.OneCode)
			oneCode.Code = httpStatusInt
			oneCode.Times = times
			_ = append(oneIp.CodeStatus, oneCode)
		}
		_ = append(statisticDTO.Ips, oneIp)
	}

	return statisticDTO
}
