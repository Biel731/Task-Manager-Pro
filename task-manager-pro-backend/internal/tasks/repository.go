package tasks

import (
	"strings"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
	"gorm.io/gorm"
)

func createOrGetTags(userID uint, names []string) ([]Tag, error) {
	if len(names) == 0 {
		return []Tag{}, nil
	}

	var tags []Tag
	for _, n := range names {
		name := strings.TrimSpace(n)
		if name == "" {
			continue
		}

		var tag Tag
		err := database.DB.Where("user_id = ? AND name = ?", userID, name).First(&tag).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				tag = Tag{
					UserID: userID,
					Name:   name,
				}
				if err := database.DB.Create(&tag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func CreateTask(userID uint, input CreateTaskInput) (*Task, error) {
	tags, err := createOrGetTags(userID, input.Tags)
	if err != nil {
		return nil, err
	}

	task := &Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    strings.ToUpper(input.Priority),
		Status:      strings.ToUpper(input.Status),
		DueDate:     input.DueDate,
		Tags:        tags,
	}

	if err := database.DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(userID uint, id uint, input UpdateTaskInput) (*Task, error) {
	var task Task
	if err := database.DB.Preload("Tags").Where("id = ? AND user = ?", id, userID).
		First(&task).Error; err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Priority != nil {
		task.Priority = strings.ToUpper(*input.Priority)
	}
	if input.Status != nil {
		task.Status = strings.ToUpper(*input.Status)
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.Tags != nil {
		tags, err := createOrGetTags(userID, *input.Tags)
		if err != nil {
			return nil, err
		}
		task.Tags = tags
	}

	if err := database.DB.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTask(UserID uint, id uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, UserID).Delete(&Task{}).Error
}

func GetTaskByID(userID uint, id uint) (*Task, error) {
	var task Task
	err := database.DB.Preload("Tags").Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func ListTasks(userID uint, filter TaskFilter) ([]Task, error) {
	db := database.DB.Preload("Tag").Where("user_id = ?", userID)

	if filter.Status != "" {
		db = db.Where("status = ?", strings.ToUpper(filter.Status))
	}
	if filter.Priority != "" {
		db = db.Where("priority = ?", strings.ToUpper(filter.Priority))
	}
	if filter.Query != "" {
		q := "%" + filter.Query + "%"
		db = db.Where("title LIKE ? OR description ILIKE ?", q, q)
	}
	if filter.Tags != "" {
		db = db.Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
			Joins("JOIN tags ON tags.id = task_tags.tag_id").Where("tags.name = ?", filter.Tags)
	}

	var tasks []Task

	if err := db.Order("create_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
