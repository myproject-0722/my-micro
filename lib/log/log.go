package log

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func InitLog(logPath string, logFileName string) {
	log.SetLevel(log.DebugLevel)

	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht),   // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(24*time.Hour),    // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		log.Printf("failed to create rotatelogs: %s", err)
		return
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{})
	log.AddHook(lfHook)

	//TODO
}
