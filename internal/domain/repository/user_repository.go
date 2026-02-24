package repository

import (
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAll(page, limit int) ([]model.User, int64, error)
}
