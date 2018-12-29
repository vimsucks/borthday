package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/BurntSushi/toml.v0"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	APP appConfig
	BOT botConfig
}

type appConfig struct {
	Mode           string
	LogFile        string
	LogStdout      bool
	LogLevel       string
	SQLiteFilePath string
}

type botConfig struct {
	Token         string
	PollerTimeout time.Duration // 用作与秒相乘
}

const (
	ModeDebug     = "debug"
	ModeRelease   = "release"
	ModeTest      = "test"
	PanicLogLevel = "panic"
	FatalLogLevel = "fatal"
	ErrorLogLevel = "error"
	WarnLogLevel  = "warn"
	InfoLogLevel  = "info"
	DebugLogLevel = "debug"
)

var (
	Conf config
	APP  *appConfig
	BOT *botConfig
)

func SetDefaultConfig(conf *config) {
	if conf.APP.Mode == "" || (conf.APP.Mode != ModeDebug && conf.APP.Mode != ModeRelease && conf.APP.Mode != ModeTest) {
		conf.APP.Mode = ModeDebug
	}

	pwd, _ := filepath.Abs(".")
	if conf.APP.LogFile == "" {
		conf.APP.LogFile = path.Join(pwd, "borthday.log")
	}

	if conf.APP.LogLevel != PanicLogLevel &&
		conf.APP.LogLevel != FatalLogLevel &&
		conf.APP.LogLevel != ErrorLogLevel &&
		conf.APP.LogLevel != WarnLogLevel &&
		conf.APP.LogLevel != DebugLogLevel &&
		conf.APP.LogLevel != InfoLogLevel {
		conf.APP.LogLevel = DebugLogLevel
	}

	if conf.APP.SQLiteFilePath == "" {
		conf.APP.SQLiteFilePath = "./borthday.db"
	}

	if conf.BOT.PollerTimeout == 0 {
		conf.BOT.PollerTimeout = 10
	}
}

func ParseConfig(configFile string) {
	if _, err := toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatal().Err(err).Msg("解析配置文件失败")
	}

	SetDefaultConfig(&Conf)

	ConfToGlobal(&Conf)
}

// 移除尾部的 /
func removeTrailingSlash(str string) string {
	if strings.HasSuffix(str, "/") {
		return str[0 : len(str)-1];
	}
	return str
}

func ConfToGlobal(conf *config) {
	APP = &conf.APP
	BOT = &conf.BOT
}
