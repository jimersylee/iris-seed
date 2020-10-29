package api_token

import (
	context2 "context"
	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/jimersylee/iris-seed/commons/redis_manager"
	"github.com/jimersylee/iris-seed/config"
	"github.com/jimersylee/iris-seed/services/cache"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	ApiCurrentUser = "ApiCurrentUser"
)

func InitTokenManager() {
	if config.Conf.RedisAddr != "" {
		config.Conf.Redis.Addr = config.Conf.RedisAddr
	}
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			Addr:     config.Conf.Redis.Addr,
			Password: config.Conf.Redis.Password,
		})),
		session.SetCookieName("mlog_session_id"),
		session.SetExpired(86400),
		session.SetEnableSIDInURLQuery(false),
		session.SetEnableSIDInHTTPHeader(false),
	)

}

func Start(ctx context.Context) session.Store {
	return StartByRequest(ctx.ResponseWriter(), ctx.Request())
}

func StartByRequest(w http.ResponseWriter, r *http.Request) session.Store {
	store, err := session.Start(context2.Background(), w, r)
	if err != nil {
		logrus.Error(err)
	}
	return store
}

func SetApiCurrentUser(token string, userId int64) {
	redis_manager.GetClient().Set(token, userId, 10*time.Minute)
}

func GetApiCurrentUser(ctx context.Context) int64 {
	return GetApiCurrentUserByRequest(ctx.Request())
}

func GetApiCurrentUserByRequest(r *http.Request) int64 {
	token := r.Header.Get("X-USER-TOKEN")
	logrus.Info("token:" + token)
	if len(token) <= 0 {
		return 0
	}
	//有token,根据token去查用户
	userId := cache.UserTokenCache.GetUserIdByToken(token)
	return userId
}

func DelApiCurrentUser(ctx context.Context) {
	token := ctx.GetHeader("X-USER-TOKEN")
	cache.UserTokenCache.Delete(token)
}
