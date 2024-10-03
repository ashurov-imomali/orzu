package usecase

import "gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"

type IUseCase interface {
	Ping() (*models.Resp, error)
	GetUserByInn(inn string) (r *models.PreCheckResponse, err error)
	SendOtp(otpReq *models.OTP) error
	ConfirmOtp(otp *models.OTP) (string, error)
	GetCondition(sSrvId, sOrzuId string) ([]models.Condition, error)
	CreateTrnash(req *models.CreateTranshReq) (*models.Resp, error)
	GetServices() ([]models.Service, error)
	CheckCard(orzuId, pan string) error
	CheckToken(token string) error
}
