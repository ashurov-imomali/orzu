package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/adapter"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/repository"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/logger"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/utils"
	"strconv"
	"strings"
	"time"
)

const (
	Created   = 1
	Confirmed = 2
	Completed = 3
)

type useCase struct {
	l logger.Logger
	i adapter.IAdapter
	r repository.IRepository
}

func New(l logger.Logger, i adapter.IAdapter, r repository.IRepository) IUseCase {
	return &useCase{l: l, i: i, r: r}
}

func (u *useCase) Ping() (*models.Resp, error) {
	return u.i.Pong()
}

func (u *useCase) GetUserByInn(inn string) (r *models.PreCheckResponse, err error) {
	if len(inn) < 1 { //todo inn length
		return nil, errors.New("too Short INN")
	}
	resp, err := u.i.GetUserByInn(inn)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(resp.Payload)
	if err != nil {
		return nil, err
	}
	if err := u.r.SetRCache(fmt.Sprintf("orzu_%d", resp.Payload.OrzuId), marshal, 50*time.Minute); err != nil { //todo 50 MINUTTTT!!!!
		return nil, err
	}
	return resp, nil
}

func (u *useCase) SendOtp(otpReq *models.OTP) error {
	ok := utils.CheckPhoneNum(&otpReq.Account)
	if !ok {
		return errors.New("invalid phone")
	}
	return u.i.SendOtp(otpReq)
}

func (u *useCase) ConfirmOtp(otp *models.OTP) (string, error) {
	if err := u.i.ConfirmOtp(otp); err != nil {
		return "", err
	}
	return utils.GenerateJWT()
}

func (u *useCase) GetServices() ([]models.Service, error) {
	return u.i.GetServices()
}

func (u *useCase) GetCondition(sSrvId, sOrzuId string) ([]models.Condition, error) {
	if _, err := strconv.Atoi(sOrzuId); err != nil {
		return nil, err
	}

	if _, err := strconv.Atoi(sSrvId); err != nil {
		return nil, err
	}

	return u.i.GetCondition(sSrvId, sOrzuId)
}

func (u *useCase) CreateTrnash(req *models.CreateTranshReq) (*models.Resp, error) {
	credit := models.OrzuCredit{
		ConditionId: req.PCredConditions,
		TerminalId:  req.TerminalId,
		Sum:         req.PSum,
		Recipient:   req.PRecipient,
		StatusId:    Created,
	}

	client, notFound, err := u.r.GetClientByOrzuId(req.POrzuId)
	if err != nil {
		return nil, err
	}

	if notFound {
		if redisErr, err := u.r.GetRCache(fmt.Sprintf("orzu_%d", req.POrzuId), &client); redisErr {
			return nil, errors.New("redis error {get}")
		} else if err != nil {
			return nil, err
		}
		if err := u.r.CreateClient(client); err != nil {
			return nil, err
		}
	}
	client.Pan = req.Pan
	credit.ClientId = client.Id
	if err := u.r.CreateCredit(&credit); err != nil {
		return nil, err
	}
	resp, err := u.i.CreateTranshe(req)
	if err != nil {
		return nil, err
	}

	if err := u.r.UpdateCreditTranshId(credit.Id, resp.Payload.(float64)); err != nil {
		return nil, err
	}

	return resp, nil
}

func (u *useCase) CheckCard(orzyId, pan string) error {
	return u.i.CheckCard(orzyId, pan)
}

func (u *useCase) CheckToken(token string) error {
	if token == "" || !strings.HasPrefix(token, "Bearer") {
		return errors.New("missing token or prefix")
	}
	return utils.JWTConfirm(strings.Replace(token, "Bearer ", "", 1))
}
