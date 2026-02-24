package repository

import (
	"github.com/google/uuid"
	"github.com/senoagung27/warehousex/internal/model"
)

type RequestRepository interface {
	Create(req *model.Request) error
	FindByID(id uuid.UUID) (*model.Request, error)
	FindAll(page, limit int, filters map[string]interface{}) ([]model.Request, int64, error)
	Update(req *model.Request) error
	FindByIDWithTx(tx interface{}, id uuid.UUID) (*model.Request, error)
	UpdateWithTx(tx interface{}, req *model.Request) error
}
