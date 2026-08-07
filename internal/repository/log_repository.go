package repository

import (
	"context"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"gorm.io/gorm"
)

type RequestLogRepository struct {
	db *gorm.DB
}

func NewRequestLogRepository(db *gorm.DB) *RequestLogRepository {
	return &RequestLogRepository{db: db}
}

func (r *RequestLogRepository) Create(ctx context.Context, log *domain.RequestLog) error {
	return r.db.WithContext(ctx).Table("logs").Create(log).Error
}

func (r *RequestLogRepository) GetAll(ctx context.Context, methodFilter string, statusCodeFilter int, page int, perPage int) ([]domain.RequestLog, int, error) {
	query := r.db.WithContext(ctx).Table("logs")

	if methodFilter != "" {
		query = query.Where("method = ?", methodFilter)
	}

	if statusCodeFilter > 0 {
		query = query.Where("status_code = ?", statusCodeFilter)
	}

	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	var logs []domain.RequestLog
	if err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}
