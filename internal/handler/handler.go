package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/usecase"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/logger"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/models"
	"math"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	u usecase.IUseCase
	l logger.Logger
}

func New(u usecase.IUseCase, l logger.Logger) Handler {
	return Handler{u: u, l: l}
}

// ping godoc
// @Summary Проверка доступности сервиса
// @Description Возвращает статус доступности сервиса
// @Tags health
// @Success 200 {object} models.Resp "Сервис доступен"
// @Failure 503 {object} map[string]interface{} "Сервис недоступен"
// @Router /ping [get]
func (h *Handler) ping(c *gin.Context) {
	response, err := h.u.Ping()
	if err != nil {
		h.l.Error(err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}
	c.JSON(response.Code, response)
}

// getUserByInn godoc
// @Summary Получение пользователя по ИНН
// @Description Возвращает данные пользователя по его ИНН
// @Tags users
// @Param inn path string true "inn" example(015306859)
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Success 200 {object} models.PreCheckResponse "Пользователь найден"
// @Example
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /user/{inn} [get]
func (h *Handler) getUserByInn(c *gin.Context) {
	inn := c.Param("inn")
	response, err := h.u.GetUserByInn(inn)
	if err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(response.Code, response)
}

// sendOtp godoc
// @Summary Отправка одноразового пароля (OTP)
// @Description Отправляет OTP на номер телефона пользователя
// @Tags otp
// @Accept json
// @Produce json
// @Param otp body models.OTP true "Данные OTP"
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Success 200 {object} models.OTP "OTP успешно отправлен"
// @Failure 400 {object} map[string]interface{} "Некорректные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /otp/send [post]
func (h *Handler) sendOtp(c *gin.Context) {
	var otp models.OTP
	if err := c.ShouldBindJSON(&otp); err != nil {
		h.l.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.u.SendOtp(&otp); err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, otp)
}

// confirmOtp godoc
// @Summary Подтверждение одноразового пароля (OTP)
// @Description Подтверждает введенный пользователем OTP
// @Tags otp
// @Accept json
// @Produce json
// @Param otp body models.OTP true "Данные OTP"
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Success 200 {object} map[string]interface{} "OTP успешно подтвержден" example({"token" : "JWTToken"})
// @Failure 400 {object} map[string]interface{} "Некорректные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /otp/confirm [post]
func (h *Handler) confirmOtp(c *gin.Context) {
	var otp models.OTP
	if err := c.ShouldBindJSON(&otp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		h.l.Warn(err.Error())
		return
	}
	token, err := h.u.ConfirmOtp(&otp)
	if err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// getServices godoc
// @Summary Получение списка сервисов
// @Description Возвращает список доступных сервисов
// @Tags services
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Success 200 {object} []models.Service "Список сервисов"
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /services [get]
func (h *Handler) getServices(c *gin.Context) {
	services, err := h.u.GetServices()
	if err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// getConditions godoc
// @Summary Получение условий по сервису
// @Description Возвращает условия по указанным service_id и orzu_id
// @Tags conditions
// @Param orzu_id query string true "orzu_id" example(9995)
// @Param service_id query string true "service_id" example(118)
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Success 200 {array} models.Condition "Условия по сервису"
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /conditions [get]
func (h *Handler) getConditions(c *gin.Context) {
	orzuId := c.Query("orzu_id")
	srvId := c.Query("service_id")
	condition, err := h.u.GetCondition(srvId, orzuId)
	if err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, condition)
}

// createTransh godoc
// @Summary Создание транша
// @Description Создает транш на основе данных, переданных в запросе
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTranshReq true "Данные для создания транша"
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Success 200 {object} models.Resp "Транш успешно создан"
// @Failure 400 {object} map[string]interface{} "Некорректные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /transactions/create [post]
func (h *Handler) createTransh(c *gin.Context) {
	var req models.CreateTranshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	trnash, err := h.u.CreateTrnash(&req)
	if err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trnash)
}

// checkUserPan godoc
// @Summary Проверка PAN пользователя
// @Description Проверяет PAN пользователя по его Orzu ID и PAN
// @Tags users
// @Param orzu_id query string true "orzu_id" example(9995)
// @Param Authorization header string true "Токен авторизации" example("Bearer your_token")
// @Param pan query string true "pan" example(4890844032435600)
// @Success 200 {object} map[string]interface{} "ok"
// @Failure 500 {object} map[string]interface{} "Ошибка на сервере"
// @Router /user/check-pan [get]
func (h *Handler) checkUserPan(c *gin.Context) {
	orzuId := c.Query("orzu_id")
	pan := c.Query("pan")
	if err := h.u.CheckCard(orzuId, pan); err != nil {
		h.l.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *Handler) ServiceAuth(hash string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token != hash {
			h.l.Warn("invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			c.Abort()
		}
		c.Next()
	}
}

func (h *Handler) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if err := h.u.CheckToken(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *Handler) Logger() gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		h.l.KVLog("hostname", hostname)
		h.l.KVLog("latency", latency)
		h.l.KVLog("clientIP", clientIP)
		h.l.KVLog("method", c.Request.Method)
		h.l.KVLog("referrer", referer)
		h.l.KVLog("dataLength", dataLength)
		h.l.KVLog("statusCode", statusCode)
		h.l.KVLog("clientUserAgent", clientUserAgent)
		if len(c.Errors) > 0 {
			h.l.Error(errors.New(c.Errors.ByType(gin.ErrorTypePrivate).String()), "ERROR")
		} else {
			format := "[ClientIP]: %s | [HostName]: %s | [Time]: %s | [Method]: %s | [Path]: %s | [StatusCode]: %d | [DataLength]: %d | [ClientUserAgent]: %s | [Latency]: %d"
			msg := fmt.Sprintf(format, clientIP, hostname, time.Now().Format("02/Jan/2006:15:04:05 +5"), c.Request.Method, path, statusCode, dataLength, clientUserAgent, latency)
			if statusCode > 499 {
				h.l.Error(errors.New(msg), "ERROR")
			} else if statusCode > 399 {
				h.l.Warn(msg)
			} else {
				h.l.Info(msg)
			}
		}
	}
}
