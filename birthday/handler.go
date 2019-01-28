package birthday

import (
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/rs/zerolog/log"
	"github.com/vimsucks/borthday/birthday/tpl"
	"github.com/vimsucks/borthday/util"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func RegisterRoute(bot *telebot.Bot) {
	bot.Handle("/lunar", func(m *telebot.Message) {
		payloads := strings.Fields(m.Payload)
		if len(payloads) < 2 {
			bot.Send(m.Sender, "命令错误，格式： /lunar 1996-10-03 周诗怡")
		}
		lunar, name := util.ParseDate(payloads[0]), payloads[1]
		uid := int64(m.Sender.ID)
		solar, err := util.LunarToSolar(payloads[0])
		if err != nil {
			bot.Send(m.Sender, "日期格式错误，格式： /lunar 1996-10-03 周诗怡")
			return
		}
		dbBirthday, err := GetBirthdayByUser(uid, name)
		if err == nil && dbBirthday != nil {
			dbBirthday.SolarBirthday = *solar
			dbBirthday.LunarBirthday = lunar
			err = UpdateBirthday(dbBirthday)
			if err != nil {
				log.Error().Err(err)
				bot.Send(m.Sender, "更新数据库失败："+err.Error())
				return
			}
			msg, _ := util.RenderTemplate(tpl.UpdateSuccess, dbBirthday)
			bot.Send(m.Sender, msg)
			return
		}
		b := Birthday{}
		b.UID = uid
		b.Name = name
		b.LunarBirthday = lunar
		b.SolarBirthday = *solar
		err = CreateBirthday(&b)
		if err != nil {
			log.Error().Err(err)
			bot.Send(m.Sender, "插入数据库失败："+err.Error())
			return
		}
		msg, _ := util.RenderTemplate(tpl.CreateSuccess, b)
		bot.Send(m.Sender, msg)
	})

	bot.Handle("/solar", func(m *telebot.Message) {
		payloads := strings.Fields(m.Payload)
		if len(payloads) < 2 {
			bot.Send(m.Sender, "命令错误，格式： /solar 1996-10-03 周诗怡")
			return
		}
		solar, name := util.ParseDate(payloads[0]), payloads[1]
		uid := int64(m.Sender.ID)
		lunar, err := util.SolarToLunar(payloads[0])
		if err != nil {
			bot.Send(m.Sender, "日期格式错误，格式： /solar 1996-10-03 周诗怡")
			return
		}
		dbBirthday, err := GetBirthdayByUser(uid, name)
		if err == nil && dbBirthday != nil {
			dbBirthday.SolarBirthday = solar
			dbBirthday.LunarBirthday = util.CCTime(lunar)
			err = UpdateBirthday(dbBirthday)
			if err != nil {
				log.Error().Err(err)
				bot.Send(m.Sender, "更新数据库失败："+err.Error())
				return
			}
			msg, _ := util.RenderTemplate(tpl.UpdateSuccess, dbBirthday)
			bot.Send(m.Sender, msg)
			return
		}
		b := Birthday{}
		b.UID = uid
		b.Name = name
		b.LunarBirthday = util.CCTime(lunar)
		b.SolarBirthday = solar
		err = CreateBirthday(&b)
		if err != nil {
			log.Error().Err(err)
			bot.Send(m.Sender, "插入数据库失败："+err.Error())
			return
		}
		msg, _ := util.RenderTemplate(tpl.CreateSuccess, b)
		bot.Send(m.Sender, msg)
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
/*
	bot.Handle("/test_solar", func(m *telebot.Message) {
		birthday, err := GetBirthdaySolarBetween("10-01", "10-05")
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

	bot.Handle("/test_lunar", func(m *telebot.Message) {
		birthday, err := GetBirthdayLunarBetween("08-18", "08-25")
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
	})*/
}
