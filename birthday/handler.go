package birthday

import (
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/rs/zerolog/log"
	"github.com/vimsucks/borthday/util"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func RegisterRoute(bot *telebot.Bot) {
	bot.Handle("/lunar", func(m *telebot.Message) {
		payloads := strings.Split(m.Payload, " ")
		if len(payloads) < 2 {
			bot.Send(m.Sender, "命令错误，格式： /lunar 1996-10-03 王书喜")
		}
		solar, err := util.LunarToSolar(payloads[0])
		if err != nil {
			bot.Send(m.Sender, "日期格式错误，格式： /lunar 1996-10-03 王书喜")
			return
		}
		b := Birthday{}
		b.UID = int64(m.Sender.ID)
		b.Name = payloads[1]
		b.LunarBirthday = util.ParseTime(payloads[0])
		b.SolarBirthday = *solar
		err = CreateBirthday(&b)
		if err != nil {
			log.Error().Err(err)
			bot.Send(m.Sender, "插入数据库失败："+err.Error())
			return
		}
		bot.Send(m.Sender, fmt.Sprintf("阳历%d年%d月%d生日", solar.Year(), solar.Month(), solar.Day()))
	})

	bot.Handle("/solar", func(m *telebot.Message) {
		payloads := strings.Split(m.Payload, " ")
		if len(payloads) < 2 {
			bot.Send(m.Sender, "命令错误，格式： /solar 1996-10-03 王书喜")
			return
		}
		lunar, err := util.SolarToLunar(payloads[0])
		if err != nil {
			bot.Send(m.Sender, "日期格式错误，格式： /solar 1996-10-03 王书喜")
			return
		}
		b := Birthday{}
		b.UID = int64(m.Sender.ID)
		b.Name = payloads[1]
		b.LunarBirthday = util.CCTime(lunar)
		b.SolarBirthday = util.ParseTime(payloads[0])
		err = CreateBirthday(&b)
		if err != nil {
			log.Error().Err(err)
			bot.Send(m.Sender, "插入数据库失败："+err.Error())
			return
		}
		bot.Send(m.Sender, fmt.Sprintf("农历%d年%d月%d生日", lunar.Year, lunar.Month, lunar.Day))
	})

	bot.Handle("/list", func(m *telebot.Message) {
		query := sqrl.Eq{"uid": m.Sender.ID}
		birthday, err := GetBirthday(query)
		if err != nil {
			log.Error().Err(err)
			bot.Send(m.Sender, "数据库查询失败："+err.Error())
			return
		}
		message := ""
		for _, b := range birthday {
			message += fmt.Sprintf("%s 农历生日%s 阳历生日%s\n", b.Name, util.DateStr(&b.LunarBirthday), util.DateStr(&b.SolarBirthday))
		}
		bot.Send(m.Sender, message)
	})
}
