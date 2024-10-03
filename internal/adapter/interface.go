package adapter

import "gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"

type IAdapter interface {
	Pong() (r *models.Resp, err error)
	GetUserByInn(inn string) (r *models.PreCheckResponse, err error)
	SendOtp(otp *models.OTP) error
	ConfirmOtp(otp *models.OTP) error
	GetServices() ([]models.Service, error)
	GetCondition(srvId, orzuId string) ([]models.Condition, error)
	CreateTranshe(req *models.CreateTranshReq) (*models.Resp, error)
	CheckCard(orzuId, pan string) error
}
