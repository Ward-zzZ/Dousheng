package initialize

import (
	"os"
	"path"
	"runtime"
	"time"

	"tiktok-demo/shared/consts"

	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger to init logrus
func InitLogger() {

	logFilepath := consts.KlogFilePath
	if err := os.MkdirAll(logFilepath, 0o777); err != nil {
		panic(err)
	}

	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilepath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			panic(err)
		}
	}

	// Provide compression and deletion
	logger := kitexlogrus.NewLogger()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time.
		MaxAge:     10,   // A file can exist for a maximum of 10 days.
		Compress:   true, // Compress with gzip.
	}

	if runtime.GOOS == "linux" {
		logger.SetOutput(lumberjackLogger)
		logger.SetLevel(klog.LevelDebug)
	} else {
		logger.SetLevel(klog.LevelDebug)
	}

	klog.SetLogger(logger)

}
