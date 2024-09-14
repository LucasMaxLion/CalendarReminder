package schedule

import (
	"CalendarReminder/app/model"
	"CalendarReminder/app/tools"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var scheduledTasks = make(map[string]*time.Timer)
var mutex sync.Mutex

func AddScheduledTask(uid string, timer *time.Timer) {
	mutex.Lock()
	defer mutex.Unlock()
	scheduledTasks[uid] = timer
}

func RemoveScheduledTask(uid string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(scheduledTasks, uid)
}

func CreateScheduledTask(reminder model.Reminder) {
	reminder.UpdateReminderStatus("1", "0")
	duration := time.Until(reminder.RemindTime)
	timer := time.AfterFunc(duration, func() {

		tools.SendEmail(reminder.Email, "日历提醒", reminder.Content)

		reminder.UpdateReminderStatus("1", "1")
		RemoveScheduledTask(reminder.Uid)
	})
	AddScheduledTask(reminder.Uid, timer)

}

func Start() {
	// 循环从 Redis 获取任务
	for {
		task, err := model.Rdb.BRPop(context.Background(), 0, "reminders").Result()
		if err != nil {
			// 处理错误，这里简单打印错误信息
			fmt.Println("Error:", err)
			continue
		}
		if task == nil || len(task) != 2 {
			continue
		}

		var reminder model.Reminder
		if err := json.Unmarshal([]byte(task[1]), &reminder); err != nil {
			// 处理反序列化错误
			fmt.Println("Error unmarshalling task:", err)
			continue
		}

		// 启动一个新的协程来处理任务
		wg.Add(1)
		go func() {
			defer wg.Done() // 确保在协程完成时减少 WaitGroup 计数
			CreateScheduledTask(reminder)
		}()
	}
}
