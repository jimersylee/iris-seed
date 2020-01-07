package cache

import (
	"github.com/goburrow/cache"
	"github.com/jimersylee/iris-seed/commons/db"
	"github.com/jimersylee/iris-seed/datamodels"
	"github.com/jimersylee/iris-seed/repositories"
	"time"
)

var UserTokenCache = newUserTokenCache()

type userTokenCache struct {
	cache cache.LoadingCache
}

func newUserTokenCache() *userTokenCache {
	return &userTokenCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserTokenRepository.FindByToken(db.GetDB(), key.(string))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(60*time.Minute),
		),
	}
}

func (this *userTokenCache) Get(token string) *datamodels.UserToken {
	if len(token) == 0 {
		return nil
	}
	val, err := this.cache.Get(token)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*datamodels.UserToken)
	}
	return nil
}

func (this *userTokenCache) Invalidate(token string) {
	this.cache.Invalidate(token)
}
