package commons

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gogo/protobuf/proto"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func InitLog() error {
	var outPutWriter io.Writer
	if Conf.App.Env == "dev" || Conf.Log.Path == "" {
		outPutWriter = os.Stdout
	} else {
		// 创建目录
		if err := createLogPath(Conf.Log.Path); err != nil {
			return err
		}

		//设置日志切割
		path := Conf.Log.Path + "/log"
		writer, _ := rotatelogs.New(
			path+".%Y%m%d",
			rotatelogs.WithLinkName(path),
			// 最多保存 7 天日志
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// 1 分割一次日志
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		outPutWriter = io.MultiWriter(writer, os.Stdout)
	}

	level, err := logrus.ParseLevel(Conf.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	// 设置SLS的相关配置信息，包括Endpoint、Project、Logstore等
	//设置日志开始
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.Infof("日志级别为:%s", level)

	if Conf.Log.SLS.Endpoint != "" {
		//设置日志结束
		// 创建SLS客户端
		producerConfig := producer.GetDefaultProducerConfig()
		producerConfig.AccessKeyID = Conf.Log.SLS.AccessKeyId
		producerConfig.AccessKeySecret = Conf.Log.SLS.AccessKeySecret
		producerConfig.Endpoint = Conf.Log.SLS.Endpoint
		producerInstance := producer.InitProducer(producerConfig)
		producerInstance.Start()

		hook := &SLSHook{
			producer: producerInstance,
			project:  Conf.Log.SLS.Project,
			logstore: Conf.Log.SLS.Logstore,
			topic:    Conf.Log.SLS.Topic,
			source:   GetIntranetIp(),
		}

		// 设置日志格式化
		// 添加SLS Hook到logger
		Log.AddHook(hook)
	}

	Log.SetLevel(level)
	Log.SetOutput(outPutWriter)
	return nil
}

func createLogPath(logPath string) error {
	_, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(Conf.Log.Path, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// SLSHook 定义一个实现了logrus.Hook接口的SLSHook结构体
type SLSHook struct {
	producer *producer.Producer
	project  string
	logstore string
	topic    string
	source   string
}

func (hook *SLSHook) Fire(entry *logrus.Entry) error {
	// 根据 entry 构建日志内容
	var content []*sls.LogContent
	content = append(content, &sls.LogContent{
		Key:   proto.String("log"),
		Value: proto.String(entry.Message),
	})
	content = append(content, &sls.LogContent{
		Key:   proto.String("level"),
		Value: proto.String(entry.Level.String()),
	})
	content = append(content, &sls.LogContent{
		Key:   proto.String("time"),
		Value: proto.String(entry.Time.String()),
	})
	traceId := entry.Data["traceId"]
	if traceId != nil {
		content = append(content, &sls.LogContent{
			Key:   proto.String("traceId"),
			Value: proto.String(traceId.(string)),
		})
	}

	log := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content,
	}
	err := hook.producer.SendLog(hook.project, hook.logstore, hook.topic, hook.source, log)
	return err
}

func (hook *SLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
