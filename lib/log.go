package lib

import (
	"fmt"
	"github.com/cxpgo/golib/model"
	"go.uber.org/zap/zapcore"
	"time"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Trace struct {
	LogTag      string
	TraceId     string
	SpanId      string
	Caller      string
	SrcMethod   string
	HintCode    int64
	HintContent string
}

type TraceContext struct {
	Trace
	CSpanId string
}

func InitLog(cfg model.Log) {
	//是否开始日志
	if !cfg.Zap.On {
		return
	}
	//Zap日志配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  cfg.Zap.TimeKey,
		LevelKey: cfg.Zap.LevelKey,
		NameKey:  cfg.Zap.NameKey,
		//CallerKey:      bc.Log.ZapConf.CallerKey,
		MessageKey:    cfg.Zap.MessageKey,
		StacktraceKey: cfg.Zap.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	//日志是否转储
	var pRotate *lumberjack.Logger
	if cfg.Rotate.On {
		pRotate = &lumberjack.Logger{
			Filename:   cfg.Zap.MainPath,      // 日志文件路径
			MaxSize:    cfg.Rotate.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: cfg.Rotate.MaxBackups, // 日志文件最多保存多少个备份
			MaxAge:     cfg.Rotate.MaxAge,     // 文件最多保存多少天
			Compress:   cfg.Rotate.Compress,   // 是否压缩
			LocalTime:  cfg.Rotate.LocalTime,
		}
	}

	// 设置日志级别
	var zLevel zapcore.Level
	err := zLevel.Set(cfg.Zap.Level)
	if err != nil {
		fmt.Printf("init zap log level is error = %v", err)
		zLevel = zap.InfoLevel
	}

	//主日志
	var mainWrite zapcore.WriteSyncer
	if cfg.Zap.Stdout && cfg.Rotate.On {
		mainWrite = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(pRotate))
	} else if cfg.Rotate.On {
		mainWrite = zapcore.NewMultiWriteSyncer(zapcore.AddSync(pRotate))
	} else if cfg.Zap.Stdout {
		mainWrite = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	} else {
		mainWrite = zapcore.NewMultiWriteSyncer()
	}

	mainLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zLevel
	})

	var core zapcore.Core
	//是否单独开启ErrorLog
	if cfg.Zap.ErrorPath != "" {
		var errorWrite zapcore.WriteSyncer
		if pRotate != nil {
			errHook := *pRotate
			errHook.Filename = cfg.Zap.ErrorPath
			errorWrite = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&errHook))
		} else {
			errorWrite = zapcore.NewMultiWriteSyncer()
		}
		errorLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), mainWrite, mainLevelEnabler),
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), errorWrite, errorLevelEnabler),
		)
	} else {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), mainWrite, mainLevelEnabler)
	}

	// 开启文件及行号
	zapLog = zap.New(core, zap.AddCaller())
	// 开启开发模式，堆栈跟踪
	if cfg.Zap.Development {
		zapLog.WithOptions(zap.Development())
	}
	zapSugarLog = zapLog.Sugar()
	zapSugarLog.Info("===>Log Init Successful<===")
}

type golibLog struct {
}

var zapLog *zap.Logger
var zapSugarLog *zap.SugaredLogger

func (l *golibLog) GetLog() *zap.Logger {
	return zapLog
}

func (l *golibLog) GetSugarLog() *zap.SugaredLogger {
	return zapSugarLog
}

func (l *golibLog) Close() {
	if zapLog != nil {
		zapLog.Sync()
	}
	if zapSugarLog != nil {
		zapSugarLog.Sync()
	}
}

func (l *golibLog) Debug(args ...interface{}) {
	zapSugarLog.Debug(args...)
}
func (l *golibLog) Debugw(msg map[string]interface{}, ctx *TraceContext) {
	zapSugarLog.Debugw(JSONMarshalToString(msg), parseTraceToSugar(ctx)...)
}
func (l *golibLog) Debugf(template string, args ...interface{}) {
	zapSugarLog.Debugf(template, args...)
}

func (l *golibLog) Info(args ...interface{}) {
	zapSugarLog.Info(args...)
}
func (l *golibLog) Infow(msg map[string]interface{}, ctx *TraceContext) {
	zapSugarLog.Infow(JSONMarshalToString(msg), parseTraceToSugar(ctx)...)
}
func (l *golibLog) Infof(template string, args ...interface{}) {
	zapSugarLog.Infof(template, args...)
}
func (l *golibLog) Printf(template string, args ...interface{}) {
	zapSugarLog.Infof(template, args...)

}

func (l *golibLog) Warn(args ...interface{}) {
	zapSugarLog.Warn(args...)
}
func (l *golibLog) Warnw(msg map[string]interface{}, ctx *TraceContext) {
	zapSugarLog.Warnw(JSONMarshalToString(msg), parseTraceToSugar(ctx)...)
}
func (l *golibLog) Warnf(template string, args ...interface{}) {
	zapSugarLog.Warnf(template, args...)
}

func (l *golibLog) Error(args ...interface{}) {
	zapSugarLog.Error(args...)
}
func (l *golibLog) Errorw(msg map[string]interface{}, ctx *TraceContext) {
	zapSugarLog.Errorw(JSONMarshalToString(msg), parseTraceToSugar(ctx)...)
}
func (l *golibLog) Errorf(template string, args ...interface{}) {
	zapSugarLog.Errorf(template, args...)
}
func (l *golibLog) DPanic(args ...interface{}) {
	zapSugarLog.DPanic(args...)
}
func (l *golibLog) DPanicw(msg map[string]interface{}, ctx *TraceContext) {
	zapSugarLog.DPanicw(JSONMarshalToString(msg), parseTraceToSugar(ctx)...)
}
func (l *golibLog) DPanicf(template string, args ...interface{}) {
	zapSugarLog.DPanicf(template, args...)
}

//map格式化为string
func parseTraceToSugar(trace *TraceContext) []interface{} {
	paramMap := map[string]interface{}{}
	if trace.LogTag != "" {
		paramMap["LogTag"] = trace.LogTag
	}
	paramMap["TraceId"] = trace.TraceId
	paramMap["SpanId"] = trace.SpanId

	if trace.Caller != "" {
		paramMap["Caller"] = trace.Caller
	}
	if trace.SrcMethod != "" {
		paramMap["SrcMethod"] = trace.SrcMethod
	}
	if trace.HintCode != 0 {
		paramMap["HintCode"] = trace.HintCode
	}
	if trace.HintContent != "" {
		paramMap["HintContent"] = trace.HintContent
	}

	var params []interface{}
	for _key, _vale := range paramMap {
		params = append(params, _key, _vale)
	}

	//log.Printf("====>%v", params)
	return params
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02-15:04:05.000"))
}
