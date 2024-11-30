package payment

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentService struct {
	payRepo *PaymentRepository
}

func NewPaymentService(payRepo *PaymentRepository, logger *zap.Logger) *PaymentService {
	return &PaymentService{payRepo: payRepo}
}

var ErrRecordNotFound = errors.New("payment not found")

func (service *PaymentService) CreatePaymentRequest(requestId uint, amount int, callbackURL string) (*string, error) {
	status := "pending"
	paymentRequest := PaymentRequest{
		RequestId:     requestId,
		Amount:        amount,
		PaymentStatus: status,
		CallbackURL:   callbackURL,
		CreatedAt:     time.Now(),
	}
	err := service.payRepo.CreatePaymentRequest(&paymentRequest)
	if err != nil {
		return nil, err
	}

	redirectURL := fmt.Sprintf("http://localhost:8083/payment/payment-page?paymentId=%v", paymentRequest.ID)

	return &redirectURL, nil
}

func (service *PaymentService) UpdatePaymentRequest(paymentIdStr, status string) (*string, error) {
	paymentId, err := strconv.ParseUint(paymentIdStr, 10, 32)
	if err != nil {
		zap.L().Error("invalid rentRequestId", zap.Error(err))
		return nil, err
	}

	paymentRequest, err := service.payRepo.GetPaymentRequestByID(uint(paymentId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	paymentRequest.PaymentStatus = status
	paymentRequest.UpdatedAt = time.Now()

	if err := service.payRepo.UpdatePaymentRequest(paymentRequest); err != nil {
		return nil, err
	}

	callbackURL := fmt.Sprintf("%v&status=%v", paymentRequest.CallbackURL, status)

	return &callbackURL, nil
}
