package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"
	"gorm.io/gorm"
	"time"
)

type repos struct {
	c *redis.Client
	p *gorm.DB
}

func NewRepos(c *redis.Client, p *gorm.DB) IRepository {
	return &repos{c, p}
}

func (r *repos) SetRCache(key string, data []byte, duration time.Duration) error {
	return r.c.Set(context.Background(),
		key,
		data,
		duration).Err()
}

func (r *repos) GetRCache(key string, data interface{}) (bool, error) {
	cmd := r.c.Get(context.Background(), key)
	bytes, err := cmd.Bytes()
	if err != nil {
		return false, err
	}
	return errors.Is(err, redis.Nil), json.Unmarshal(bytes, data)
}

func (r *repos) CreateCredit(credit *models.OrzuCredit) error {
	return r.p.Create(credit).Error
}

func (r *repos) GetClientByOrzuId(orzuId int) (*models.OrzuClient, bool, error) {
	var c models.OrzuClient
	err := r.p.Where("orzu_id = ?", orzuId).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.OrzuClient{}, true, nil
	}
	return &c, false, err
}

func (r *repos) CreateClient(client *models.OrzuClient) error {
	return r.p.Create(client).Error
}

func (r *repos) UpdateCreditTranshId(id int64, transhId float64) error {
	return r.p.Where("id = ?", id).Update("transh_id", transhId).Error
}
