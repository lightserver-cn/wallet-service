package logger

import (
	fileRotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"server/pkg/consts"

	"io"
	"strings"
	"time"
)

var Logger *zap.SugaredLogger

func InitZapLogger(filepath, infoFilename, warnFilename, errFilename, fileExt, callerLoc string) (*zap.SugaredLogger, error) {
	cfg := &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			TimeKey:     "ts",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(consts.TimestampFormat1))
			},
			CallerKey: "caller",
			EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
				fullPath := caller.FullPath()
				encoder.AppendString(fullPath[strings.Index(fullPath, callerLoc)+len(callerLoc):])
			},
			EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendInt64(int64(d) / (consts.Thousand * consts.Thousand))
			},
		},
		OutputPaths:      []string{"stdout", filepath + "/output.txt"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := newZapLogger(filepath, infoFilename, warnFilename, errFilename, fileExt, cfg)

	return logger.Sugar(), err
}

func newZapLogger(logFilepath, infoFilename, warnFilename, errFilename, fileExt string, cfg *zap.Config) (*zap.Logger, error) {
	encoder := getEncoder(cfg)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})

	errLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel
	})

	infoWriter, err := getLogWriter(logFilepath+"/"+infoFilename, fileExt)
	if err != nil {
		return nil, err
	}

	warnWriter, err2 := getLogWriter(logFilepath+"/"+warnFilename, fileExt)
	if err2 != nil {
		return nil, err2
	}

	errWriter, err3 := getLogWriter(logFilepath+"/"+errFilename, fileExt)
	if err3 != nil {
		return nil, err3
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
		zapcore.NewCore(encoder, errWriter, errLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
}

// getEncoder
func getEncoder(conf *zap.Config) zapcore.Encoder {
	var enc zapcore.Encoder

	switch conf.Encoding {
	case "json":
		enc = zapcore.NewJSONEncoder(conf.EncoderConfig)
	case "console":
		enc = zapcore.NewConsoleEncoder(conf.EncoderConfig)
	default:
		panic("unknown encoding")
	}

	return enc
}

// getLogWriter
func getLogWriter(filePath, fileExt string) (zapcore.WriteSyncer, error) {
	warnIoWriter, err := getWriter(filePath, fileExt)
	if err != nil {
		return nil, err
	}

	return zapcore.AddSync(warnIoWriter), nil
}

// getWriter 日志文件切割，按天
func getWriter(filename, fileExt string) (io.Writer, error) {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	hook, err := fileRotatelogs.New(
		filename+"_%Y%m%d."+fileExt,
		fileRotatelogs.WithLinkName(filename),
		fileRotatelogs.WithMaxAge(consts.Day30),
		fileRotatelogs.WithRotationTime(consts.Day),
	)

	return hook, err
}
