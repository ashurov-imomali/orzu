package models

type Resp struct {
	HeaderResp
	Payload interface{} `json:"payload"`
}

type HeaderResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PreCheckResponse struct {
	HeaderResp
	Payload Payload `json:"payload"`
}

type Payload struct {
	OrzuId            int    `json:"orzu_id"`
	ClientName        string `json:"client_name"`
	SetDate           string `json:"set_date"`
	PhoneNumber       string `json:"phone_number"`
	PassportId        string `json:"passport_id"`
	PassportIssueDate string `json:"passport_issue_date"`
}

type Request struct {
	Inn string `json:"inn"`
}

type OTP struct {
	ID           string `json:"id"`
	Account      string `json:"account"`
	Value        string `json:"value"`
	Lifetime     int64  `json:"lifetime"`
	ConfirmLimit int64  `json:"validate_limit"`
	State        int64  `json:"state"`
	CreatedAt    string `json:"created_at"`
	ExpiredAt    string `json:"expired_at"`
}

type SrvResponse struct {
	HeaderResp
	Payload []Service `json:"payload"`
}

type Service struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	ExternalID string `json:"externalID"`
}

type CndResponse struct {
	HeaderResp
	Payload []Condition `json:"payload"`
}

type Condition struct {
	ConditionId   int64   `json:"condition_id"`
	Name          string  `json:"name"`
	MinSumma      float64 `json:"min_summa"`
	MaxSumma      float64 `json:"max_summa"`
	IntervalUnits string  `json:"interval_units"`
	Term          int     `json:"term"`
	ClComission   float32 `json:"cl_comission"`
	PrcRate       float32 `json:"prc_rate"`
}

type OrzuCredit struct {
	Id          int64   `json:"id" gorm:"column:id"`
	TerminalId  int64   `json:"terminal_id" gorm:"column:terminal_id"`
	Sum         float64 `json:"sum" gorm:"column:sum"`
	Recipient   string  `json:"recipient" gorm:"recipient"`
	ConditionId int64   `json:"condition_id" gorm:"column:condition_id"`
	StatusId    int64   `json:"status_id" gorm:"column:status_id"` //todo
	TranshId    string  `json:"transh_id" gorm:"column:transh_id"`
	ClientId    int64   `json:"client_id" gorm:"column:client_id"`
}

type OrzuClient struct {
	Id                int64  `json:"id" gorm:"column:id"`
	OrzuId            int    `json:"orzu_id" gorm:"column:orzu_id"`
	Pan               string `json:"pan" gorm:"column:pan"`
	Name              string `json:"client_name" gorm:"column:name"`
	PhoneNumber       string `json:"phone_number" gorm:"column:phone_number"`
	SetDate           string `json:"set_date" gorm:"column:set_date"` //todo
	PassportId        string `json:"passport_id" gorm:"column:passport_id"`
	PassportIssueDate string `json:"passport_issue_date" gorm:"column:passport_issue_date"`
}

type CreateTranshReq struct {
	TerminalId      int64   `json:"terminal_id" gorm:"column:terminal_id"`
	POrzuId         int     `json:"pOrzuId"`
	PInn            string  `json:"pInn"`
	Pan             string  `json:"pan"`
	PSum            float64 `json:"pSum"`
	PServiceId      string  `json:"pServiceId"`
	PToken          string  `json:"pToken"`
	PCredConditions int64   `json:"pCredConditions"`
	PRecipient      string  `json:"pRecipient"`
	Hash            string  `json:"hash"`
	PPhoneUUID      string  `json:"pPhoneUUID"`
}
