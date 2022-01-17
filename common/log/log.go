package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/shooketh/sakura/common/config"
	"github.com/shooketh/sakura/common/constants"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func Init() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := setLevel(); err != nil {
		return err
	}

	l := &lumberjack.Logger{
		Filename: fmt.Sprintf("%s/%s", config.Config.Log.Path, constants.LogFileName),
		MaxSize:  config.Config.Log.MaxSize,
		MaxAge:   config.Config.Log.MaxAge,
		Compress: config.Config.Log.Compress,
	}

	o := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	w := io.MultiWriter(o, l)
	Logger = zerolog.New(w).With().Timestamp().Caller().Logger()

	return nil
}

func setLevel() error {
	switch config.Config.Log.Level {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		return fmt.Errorf("set log level failed")
	}
	return nil
}
