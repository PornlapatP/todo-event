package adapter

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"
	"time"

	taskdomain "todoe/domain/task/domain"
)

type CSVRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewCSVRepository(filePath string) *CSVRepository {
	return &CSVRepository{filePath: filePath}
}
func (r *CSVRepository) AppendCompletedTask(task taskdomain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if task.CompletedAt == nil {
		return nil
	}

	fileExists := true
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		fileExists = false
	}

	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if !fileExists {
		if err := writer.Write([]string{
			"id",
			"title",
			"status",
			"created_at",
			"completed_at",
			"sla_seconds",
			"sla_human",
		}); err != nil {
			return err
		}
	}

	completedAt := *task.CompletedAt
	sla := completedAt.Sub(task.CreatedAt)

	return writer.Write([]string{
		task.ID.Hex(),
		task.Title,
		string(task.Status),
		task.CreatedAt.Format(time.RFC3339),
		completedAt.Format(time.RFC3339),
		strconv.FormatInt(int64(sla.Seconds()), 10),
		sla.String(),
	})
}
