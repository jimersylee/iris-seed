package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/jimersylee/go-steam-proxy/commons/redis_manager"
	"github.com/sirupsen/logrus"
)

var UserTokenCache *userTokenCache

func newUserTokenCache() *userTokenCache {
	client := redis_manager.GetClient()
	return &userTokenCache{redisClient: client}
}

type userTokenCache struct {
	redisClient *redis.Client
}

func (this *userTokenCache) GetUserIdByToken(token string) int64 {
	if UserTokenCache == nil {
		this = newUserTokenCache()
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

func (this *userTokenCache) Delete(token string) {
	if UserTokenCache == nil {
		this = newUserTokenCache()
	}
	logrus.Info("token:" + token)
	this.redisClient.Del(token)
}
