package repository

import (
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"
	"time"
)

type IRepository interface {
	GetRCache(key string, data interface{}) (bool, error)
	SetRCache(key string, data []byte, duration time.Duration) error
	CreateCredit(credit *models.OrzuCredit) error
	GetClientByOrzuId(orzuId int) (*models.OrzuClient, bool, error)
	CreateClient(client *models.OrzuClient) error
	UpdateCreditTranshId(id int64, transhId float64) error
}
