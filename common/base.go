package common

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BaseInfra struct {
	Conf   *TaoConf
	Logger *zap.SugaredLogger
}

var GlbBaInfa *BaseInfra
var Logger *zap.SugaredLogger

func BaseInit(path string) *BaseInfra {
	GlbBaInfa = &BaseInfra{}
	GlbBaInfa.Conf = &TaoConf{}
	GlbBaInfa.Conf.LoadTaoConf(path)

	syncer := initLogWriter(GlbBaInfa.Conf)
	encoder := initEncoder()
	level, err := zapcore.ParseLevel(GlbBaInfa.Conf.Log.Level)
	if err != nil {
		log.Fatalf("ParseLevel:%s failed", GlbBaInfa.Conf.Log.Level)
		level = zapcore.InfoLevel
	}

	highPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return (l >= level)
	})

	console := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	var core zapcore.Core
	if GlbBaInfa.Conf.Log.Env == "dev" {
		core = zapcore.NewTee(zapcore.NewCore(encoder, syncer, highPriority), zapcore.NewCore(consoleEncoder, console, highPriority))
	} else {
		core = zapcore.NewCore(encoder, syncer, highPriority)
	}

	GlbBaInfa.Logger = zap.New(core).Sugar()
	defer GlbBaInfa.Logger.Sync()
	Logger = GlbBaInfa.Logger
	return GlbBaInfa
}

func initLogWriter(c *TaoConf) zapcore.WriteSyncer {
	logger := &lumberjack.Logger{
		Filename:   c.Log.File,
		MaxSize:    c.Log.MaxSize,
		MaxBackups: c.Log.MaxBackups,
		MaxAge:     c.Log.MaxAge,
		Compress:   false,
	}
	return zapcore.AddSync(logger)
}

func initEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(cfg)
}
