package logic

import (
	"CalendarReminder/app/model"
	"CalendarReminder/app/schedule"
	"CalendarReminder/app/tools"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CreateReminder(c *gin.Context) {
	var rem model.Reminder
	if err := c.ShouldBindJSON(&rem); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    0,
			Message: "未获得信息",
		})
		return
	}
	if err := CheckTime(rem.RemindTime); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(),
		})
		return
	}
	uid, _ := c.Cookie("uid")
	err := model.AddReminder(uid, rem.Content, rem.RemindTime)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "新增失败",
		})
		return
	}
	publishTaskToRedis(rem)

	c.JSON(http.StatusOK, tools.ECode{Message: "新增成功"})
}

func GetReminders(c *gin.Context) {
	uid, _ := c.Cookie("uid")
	reminders, err := model.GetRemindersByCreator(uid)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "获取列表数据失败"})
		return
	}
	c.JSON(http.StatusOK, reminders)
}

func DeleteReminder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "获取留言id失败"})
		return
	}
	err = model.DeleteReminder(strconv.Itoa(id))
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "删除失败"})
		return
	}
	if err = model.DeleteTaskFromRedis(strconv.Itoa(id)); err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "删除失败"})
		return
	}
	schedule.RemoveScheduledTask(strconv.Itoa(id))
	c.JSON(http.StatusOK, tools.ECode{Message: "删除成功"})
}

func UpdateReminder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "获取ID失败"})
		return
	}
	var updates map[string]interface{}
	if err = c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, tools.ECode{Message: "解析失败"})
		return
	}

	var content string
	var remindTime time.Time
	if v, ok := updates["content"]; ok {
		content = v.(string)
	}
	if v, ok := updates["remind_time"]; ok {
		remindTime, err = time.Parse(time.RFC3339, v.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, tools.ECode{Message: "时间格式不正确"})
			return
		}
		if err = CheckTime(remindTime); err != nil {
			c.JSON(http.StatusOK, tools.ECode{
				Message: err.Error(),
			})
			return
		}
	}
	err = model.UpdateReminder(strconv.Itoa(id), content, remindTime)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "更新失败"})
		return
	}
	if err = model.UpdateTaskInRedis(strconv.Itoa(id), content, remindTime); err != nil {
		c.JSON(http.StatusOK, tools.ECode{Message: "更新失败"})
		return
	}
	c.JSON(http.StatusOK, tools.ECode{Message: "更新成功"})
}

func CheckTime(rem_time time.Time) error {
	currentTime := time.Now()

	// 比较提醒时间是否小于当前时间
	if rem_time.Before(currentTime) {
		return &Error{message: "您输入的时间比当前时间小，这是不对的"}
	}
	return nil
}

type Error struct {
	message string
}

// 实现 error 接口的 Error 方法
func (e *Error) Error() string {
	return e.message
}

func publishTaskToRedis(rem model.Reminder) {
	jsonTask, _ := json.Marshal(rem)
	model.Rdb.RPush(context.Background(), "reminders", jsonTask).Result()
}
