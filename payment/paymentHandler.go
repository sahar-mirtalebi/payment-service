package payment

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	service  *PaymentService
	validate *validator.Validate
}

func NewPaymentHandler(service *PaymentService, logger *zap.Logger, validate *validator.Validate) *PaymentHandler {
	return &PaymentHandler{service: service, validate: validate}
}

func (handler *PaymentHandler) CreatePaymentRequest(c echo.Context) error {

	var paymentRequest struct {
		RequestId   uint   `json:"requestId" validate:"required"`
		Amount      int    `json:"amount" validate:"required"`
		CallbackURL string `json:"callbackURL" validate:"required"`
	}

	if err := c.Bind(&paymentRequest); err != nil {
		zap.L().Error("failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := handler.validate.Struct(paymentRequest); err != nil {
		zap.L().Error("provided data is invalid", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	redirectURL, err := handler.service.CreatePaymentRequest(paymentRequest.RequestId, paymentRequest.Amount, paymentRequest.CallbackURL)
	if err != nil {
		zap.L().Error("Error creating payment request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create payment request")
	}

	return c.JSON(http.StatusCreated, map[string]string{"redirectURL": *redirectURL})
}

func (handler *PaymentHandler) UpdatePaymentRequest(c echo.Context) error {
	var paymentInfo struct {
		PaymentIdSTR string `json:"paymentId"`
		Status       string `json:"status"`
	}

	if err := c.Bind(&paymentInfo); err != nil {
		zap.L().Error("failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	callbackURL, err := handler.service.UpdatePaymentRequest(paymentInfo.PaymentIdSTR, paymentInfo.Status)
	if err != nil {
		zap.L().Error("error updating status", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update status")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"callbackUrl": *callbackURL,
	})
}
