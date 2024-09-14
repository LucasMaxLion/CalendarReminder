package router

import (
	"CalendarReminder/app/logic"
	"github.com/gin-gonic/gin"
)

func New() {
	r := gin.Default()
	r.POST("/login", logic.Login)

	r.POST("/reminders", logic.CreateReminder)
	r.GET("/reminders", logic.GetReminders)
	r.DELETE("/reminders/:id", logic.DeleteReminder)
	r.PATCH("/reminders/:id", logic.UpdateReminder)

	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}
