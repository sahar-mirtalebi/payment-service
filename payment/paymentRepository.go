package payment

import (
	"time"

	"gorm.io/gorm"
)

type PaymentRequest struct {
	ID            uint
	RequestId     uint
	Amount        int
	PaymentStatus string
	CallbackURL   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (payRepo *PaymentRepository) CreatePaymentRequest(paymentRequest *PaymentRequest) error {
	return payRepo.db.Create(paymentRequest).Error
}

func (payRepo *PaymentRepository) GetPaymentRequestByID(paymentId uint) (*PaymentRequest, error) {
	var paymentRequest PaymentRequest
	err := payRepo.db.First(&paymentRequest, paymentId).Error
	if err != nil {
		return nil, err
	}
	return &paymentRequest, nil
}

func (payRepo *PaymentRepository) UpdatePaymentRequest(paymentRequest *PaymentRequest) error {
	return payRepo.db.Save(paymentRequest).Error
}
