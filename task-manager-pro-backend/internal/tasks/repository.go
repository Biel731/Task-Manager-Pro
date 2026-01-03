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
	if err := database.DB.Preload("Tags").Where("id = ? AND user_id = ?", id, userID).
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

func DeleteTask(userID uint, taskID uint) error {
	var task Task

	// 1) Buscar a task do usuário
	if err := database.DB.
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error; err != nil {
		return err
	}

	// 2) Limpar as associações na tabela task_tags
	if err := database.DB.
		Model(&task).
		Association("Tags").
		Clear(); err != nil {
		return err
	}

	// 3) Deletar a task
	if err := database.DB.Delete(&task).Error; err != nil {
		return err
	}

	return nil
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
	db := database.DB.
		Model(&Task{}).
		Preload("Tags").
		Where("tasks.user_id = ?", userID)

	// filtros simples
	if filter.Status != "" {
		db = db.Where("tasks.status = ?", strings.ToUpper(filter.Status))
	}
	if filter.Priority != "" {
		db = db.Where("tasks.priority = ?", strings.ToUpper(filter.Priority))
	}

	// 🔍 Query: busca em title, description E tags.name
	if filter.Query != "" {
		q := "%" + strings.TrimSpace(filter.Query) + "%"

		db = db.
			Joins("LEFT JOIN task_tags ON task_tags.task_id = tasks.id").
			Joins("LEFT JOIN tags ON tags.id = task_tags.tag_id").
			Where(`
      (
        tasks.title ILIKE ?
        OR tasks.description ILIKE ?
        OR tags.name ILIKE ?
      )
    `, q, q, q).
			Group("tasks.id")
	}

	// filtro por tag exata (se você quiser continuar suportando)
	if filter.Tags != "" {
		db = db.
			Joins("JOIN task_tags tt ON tt.task_id = tasks.id").
			Joins("JOIN tags t ON t.id = tt.tag_id").
			Where("t.name = ?", strings.TrimSpace(filter.Tags)).
			Distinct("tasks.id")
	}

	var tasks []Task
	if err := db.Order("tasks.created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
