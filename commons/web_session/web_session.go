package web_session

import (
	context2 "context"
	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/jimersylee/go-steam-proxy/config"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	CurrentUser = "CurrentUser"
)

func InitSessionManager() {
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

func SetCurrentUser(ctx context.Context, userId int64) {
	store := Start(ctx)
	store.Set(CurrentUser, strconv.FormatInt(userId, 10))
	err := store.Save()
	if err != nil {
		logrus.Error(err)
	}
}

func GetCurrentUser(ctx context.Context) int64 {
	return GetCurrentUserByRequest(ctx.ResponseWriter(), ctx.Request())
}

func GetCurrentUserByRequest(w http.ResponseWriter, r *http.Request) int64 {
	val, exists := StartByRequest(w, r).Get(CurrentUser)
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

func DelCurrentUser(ctx context.Context) {
	store := Start(ctx)
	store.Delete(CurrentUser)
	err := store.Save()
	if err != nil {
		logrus.Error(err)
	}
}
