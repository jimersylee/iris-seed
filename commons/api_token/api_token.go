package api_token

import (
	context2 "context"
	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/jimersylee/iris-seed/config"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

func SetApiCurrentUser(ctx context.Context, userId int64) {
	store := Start(ctx)
	store.Set(ApiCurrentUser, strconv.FormatInt(userId, 10))
	err := store.Save()
	if err != nil {
		logrus.Error(err)
	}
}

func GetApiCurrentUser(ctx context.Context) int64 {
	return GetApiCurrentUserByRequest(ctx.ResponseWriter(), ctx.Request())
}

func GetApiCurrentUserByRequest(w http.ResponseWriter, r *http.Request) int64 {
	val, exists := StartByRequest(w, r).Get(ApiCurrentUser)
	if exists {
		switch val.(type) {
		case string:
			userId, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return 0
			}
			return userId
		}
	}
	return 0
}

func DelApiCurrentUser(ctx context.Context) {
	store := Start(ctx)
	store.Delete(ApiCurrentUser)
	err := store.Save()
	if err != nil {
		logrus.Error(err)
	}
}
