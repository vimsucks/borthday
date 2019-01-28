package main

import (
	"flag"
	"github.com/robfig/cron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vimsucks/borthday/birthday"
	"github.com/vimsucks/borthday/birthday/tpl"
	"github.com/vimsucks/borthday/config"
	"github.com/vimsucks/borthday/util"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	config.ParseFlag()
	if config.Flag.Help {
		flag.Usage()
		os.Exit(0)
	}
	// 解析配置文件
	config.ParseConfig(config.Flag.ConfigFile)

	initLogger()
	birthday.ConnectToSQLite(config.APP.SQLiteFilePath)

	// 连接 Redis 和 OA数据库，因为要用到 config

	bot, err := tb.NewBot(tb.Settings{
		Token:  config.BOT.Token,
		Poller: &tb.LongPoller{Timeout: config.BOT.PollerTimeout * time.Second},
	})

	if err != nil {
		log.Fatal().Err(err)
		return
	}

	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})

	startCron(bot)

	birthday.RegisterRoute(bot)

	bot.Start()
}

func initLogger() {
	var logWriter io.Writer = getLogFile(config.APP.LogFile)
	if config.APP.LogStdout {
		logWriter = io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logWriter)
	}

	var level zerolog.Level
	switch config.APP.LogLevel {
	case config.PanicLogLevel:
		level = zerolog.PanicLevel
	case config.FatalLogLevel:
		level = zerolog.FatalLevel
	case config.ErrorLogLevel:
		level = zerolog.ErrorLevel
	case config.WarnLogLevel:
		level = zerolog.WarnLevel
	case config.DebugLogLevel:
		level = zerolog.DebugLevel
	case config.InfoLogLevel:
		level = zerolog.InfoLevel
	default:
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	//log.Logger = zerolog.New(logWriter).With().Timestamp().Logger()
	log.Logger = zerolog.New(logWriter).With().Timestamp().Logger()
}

func getLogFile(filePath string) *os.File {
	var file *os.File
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err = os.Create(filePath)
	} else if err != nil {
		log.Fatal().Str("file", filePath).Err(err).Msg("创建文件失败")
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal().Str("file", filePath).Err(err).Msg("打开文件失败")
		}
	}
	return file
}

func startCron(bot *tb.Bot) {
	c := cron.New()
	c.AddFunc("0 0 1 * * ?", func() {
		today := time.Now()
		fivesDayAfter := today.AddDate(0, 0, 5)
		birthdays, err := birthday.GetBirthdaySolarBetween(util.MonthDay(&today), util.MonthDay(&fivesDayAfter))
		notified := make(map[string]int)
		if err != nil {
			log.Error().Err(err)
		}
		for _, b := range birthdays {
			receiver := &tb.User{}
			receiver.ID = int(b.UID)
			key := b.Name + strconv.FormatInt(b.UID, 10)
			if _, ok := notified[key]; !ok {
				msg, _ := util.RenderTemplate(tpl.IncomingBirthday, b)
				bot.Send(receiver, msg)
				notified[key] = 1
			}
		}
		lunarTodayCal, _ := util.SolarToLunar(util.DateStr(&today))
		lunarToday := util.CCTime(lunarTodayCal)
		lunarFivesDayAfter := lunarToday.AddDate(0, 0, 5)
		birthdays, err = birthday.GetBirthdayLunarBetween(util.MonthDay(&lunarToday), util.MonthDay(&lunarFivesDayAfter))
		if err != nil {
			log.Error().Err(err)
		}
		for _, b := range birthdays {
			receiver := &tb.User{}
			receiver.ID = int(b.UID)
			key := b.Name + strconv.FormatInt(b.UID, 10)
			if _, ok := notified[b.Name + strconv.FormatInt(b.UID, 10)]; !ok {
				msg, _ := util.RenderTemplate(tpl.IncomingBirthday, b)
				bot.Send(receiver, msg)
				notified[key] = 1
			}
		}
	})
	c.Start()
}
