package app

import (
	"CalendarReminder/app/model"
	"CalendarReminder/app/router"
	"CalendarReminder/app/schedule"
)

func Start() {
	model.NewMysql()
	model.NewRdb()
	defer func() {
		model.MySQLClose()
	}()
	defer func() {
		model.RedisClose()
	}()
	go schedule.Start()
	router.New()
}
