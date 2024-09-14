package model

import (
	"CalendarReminder/app/tools"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func AddReminder(creatorId, content string, remindTime time.Time) error {

	rem := Reminder{
		Uid:         strconv.Itoa(int(tools.GetUid())) + creatorId,
		CreatorId:   creatorId,
		Content:     content,
		RemindTime:  remindTime,
		Begin:       "0",
		Active:      "0",
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	Conn.Table("reminder").Save(&rem)
	return nil
}

func GetRemindersByCreator(creatorID string) ([]Reminder, error) {
	var filteredReminders []Reminder
	Conn.Table("reminder").Where("credtor_id = ?", creatorID).Find(&filteredReminders)
	return filteredReminders, nil
}

func DeleteReminder(uid string) error {
	var reminder Reminder
	Conn.Table("reminder").Where("uid = ?", uid).Find(&reminder)
	if err := Conn.Delete(&reminder).Error; err != nil {
		return err
	}
	return nil
}

func UpdateReminder(uid, content string, remindTime time.Time) error {
	var reminder Reminder
	if err := Conn.Table("reminder").Where("uid = ?", uid).Find(&reminder).Error; err != nil {
		return err
	}
	reminder.RemindTime = remindTime
	reminder.Content = content

	return nil
}

func UpdateTaskInRedis(uid, content string, remindTime time.Time) error {
	// 读取 Redis 中的所有任务
	tasks, err := Rdb.LRange(context.Background(), "reminders", 0, -1).Result()
	if err != nil {
		return err
	}

	// 查找需要更新的任务
	for i, task := range tasks {
		var reminder Reminder
		if err := json.Unmarshal([]byte(task), &reminder); err != nil {
			return err
		}

		// 如果找到匹配的任务，更新内容和提醒时间
		if reminder.Uid == uid {
			reminder.Content = content
			reminder.RemindTime = remindTime
			reminder.UpdatedTime = time.Now()

			// 序列化更新后的任务
			updatedTask, err := json.Marshal(reminder)
			if err != nil {
				return err
			}

			// 写回 Redis 中的特定位置
			if err := Rdb.LSet(context.Background(), "reminders", int64(i), string(updatedTask)).Err(); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("未找到 UID 为 %s 的任务", uid)
}

func DeleteTaskFromRedis(uid string) error {
	// 读取 Redis 中的所有任务
	tasks, err := Rdb.LRange(context.Background(), "reminders", 0, -1).Result()
	if err != nil {
		return err
	}

	// 查找需要删除的任务索引
	var index int
	for i, task := range tasks {
		var reminder Reminder
		if err := json.Unmarshal([]byte(task), &reminder); err != nil {
			return err
		}
		if reminder.Uid == uid {
			index = i
			break
		}
	}

	// 如果找到了任务，使用 LREM 命令移除它
	if index != 0 {
		if _, err := Rdb.LRem(context.Background(), "reminders", 1, tasks[index]).Result(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("未找到 UID 为 %s 的任务", uid)
	}

	return nil
}
