package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/config"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/logger"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/utils"
	"io"
	"log"
	"net/http"
	"time"
)

type adapter struct {
	orzu config.OrzuParams
	otp  config.OtpParams
	c    http.Client
	l    logger.Logger
}

func New(l logger.Logger, p config.OrzuParams, o config.OtpParams) IAdapter {
	return &adapter{c: http.Client{
		Timeout: 30 * time.Second,
	}, l: l, orzu: p, otp: o}
}

func (a *adapter) Pong() (r *models.Resp, err error) {
	defer a.l.Println(&r)
	request, err := http.NewRequest(http.MethodGet, a.orzu.Url+"/ping", nil)
	if err != nil {
		return
	}

	a.l.Printf("Sending {PING} request to %s", a.orzu.Url)
	response, err := a.c.Do(request)
	if err != nil {
		return
	}
	return r, json.NewDecoder(response.Body).Decode(&r)
}

func (a *adapter) GetUserByInn(inn string) (r *models.PreCheckResponse, err error) {
	defer a.l.Println(&r)
	request, err := http.NewRequest(http.MethodGet, a.orzu.Url+"/getClientByInnNew/"+inn, nil)
	if err != nil {
		return
	}
	request.Header.Add("token", a.orzu.Token)
	a.l.Printf("Sending {GET_CLIENT_BY_INN} %s", request.URL.String())
	response, err := a.c.Do(request)
	if err != nil {
		return
	}

	return r, json.NewDecoder(response.Body).Decode(&r)
}

func (a *adapter) SendOtp(otp *models.OTP) error {
	otp.Lifetime = a.otp.LifeTime
	otp.ConfirmLimit = a.otp.ConfirmLimit
	indent, err := json.Marshal(&otp)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, a.otp.Url, bytes.NewBuffer(indent))
	if err != nil {
		return err
	}
	a.l.Printf("Sending {OTP} request to %s with body %s", a.otp.Url, string(indent))
	response, err := a.c.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return errors.New(response.Status)
	}
	defer a.l.Println(otp)
	return json.NewDecoder(response.Body).Decode(otp)
}

func (a *adapter) ConfirmOtp(otp *models.OTP) error {
	marshal, err := json.Marshal(&otp)
	if err != nil {
		return err
	}
	a.l.Printf("sending request to %s with body %s", a.otp.Url+"/"+otp.ID, string(marshal))
	request, err := http.NewRequest(http.MethodPatch, a.otp.Url+"/"+otp.ID, bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	response, err := a.c.Do(request)
	if err != nil {
		return err
	}
	a.l.Printf("response from otp service %v", response)
	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (a *adapter) GetServices() ([]models.Service, error) {
	request, err := http.NewRequest(http.MethodGet, a.orzu.Url+"/getServices/"+a.orzu.TToken, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("token", a.orzu.Token)
	a.l.Printf("Sending {GET_SERVICES} %s", request.URL.String())
	response, err := a.c.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var s models.SrvResponse
	if err := json.NewDecoder(response.Body).Decode(&s); err != nil {
		return nil, err
	}
	if s.Code != http.StatusOK {
		return nil, errors.New(s.Message)
	}
	return s.Payload, nil
}

func (a *adapter) GetCondition(srvId, orzuId string) ([]models.Condition, error) {
	request, err := http.NewRequest(http.MethodGet, a.orzu.Url+"/getServiceConditions", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("token", a.orzu.Token)
	query := request.URL.Query()
	query.Add("service_id", srvId)
	query.Add("orzu_id", orzuId)
	query.Add("token", a.orzu.TToken)
	request.URL.RawQuery = query.Encode()
	a.l.Printf("Sending {GET_SERVICE_CONDITIONS} to %s", request.URL.String())
	response, err := a.c.Do(request)
	if err != nil {
		return nil, err
	}
	var resp models.CndResponse
	all, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	a.l.Printf("Response from Orzu service: %s", string(all))
	err = json.Unmarshal(all, &resp)
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	return resp.Payload, nil
}

func (a *adapter) CreateTranshe(req *models.CreateTranshReq) (*models.Resp, error) {
	req.Hash = utils.GetSha256Hash(a.orzu.SecretKey, req.POrzuId, a.orzu.ServiceId, a.orzu.Token)
	req.PToken = a.orzu.Token //todo
	req.PServiceId = a.orzu.ServiceId
	req.PPhoneUUID = "2333-33232-4434-3433" //todo delete after tests
	marshal, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, a.orzu.Url+"/orzupay/pay", bytes.NewBuffer(marshal))
	if err != nil {
		return nil, err
	}
	request.Header.Add("token", a.orzu.Token)
	a.l.Printf("Sending {CREATE_TRANSHE} to %s with body: %s", a.orzu.Url+"/orzupay/createTranch", string(marshal))
	response, err := a.c.Do(request)
	if err != nil {
		return nil, err
	}
	var resp models.Resp
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, err
	}
	log.Println(resp)
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	return &resp, nil
}

func (a *adapter) CheckCard(orzuId, pan string) error {
	request, err := http.NewRequest(http.MethodGet, a.orzu.Url+"/preCheckCard", nil)
	if err != nil {
		return err
	}
	request.Header.Add("token", a.orzu.Token)
	query := request.URL.Query()
	query.Add("pan", pan)
	query.Add("orzu_id", orzuId)
	request.URL.RawQuery = query.Encode()
	a.l.Printf("Sending {CHECK_CARD} to %s", request.URL.String())
	response, err := a.c.Do(request)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		a.l.Warn(err)
	}
	a.l.Printf("response from Orzu service: %s", string(body))
	if response.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
