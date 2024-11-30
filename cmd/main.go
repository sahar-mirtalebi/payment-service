package main

import (
	"log"
	"payment-service/payment"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func NewLogger() (*zap.Logger, error) {
// 	logger, err := zap.NewProduction()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return logger, nil
// }

func NewValidator() *validator.Validate {
	return validator.New()
}

func NewDB() (*gorm.DB, error) {
	dsn := "host=localhost user=admin password=sahar223010 dbname=rental_service_db search_path=payment-service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RegisteredRoutes(e *echo.Echo, handler *payment.PaymentHandler) {

	e.GET("/payment/payment-page", func(c echo.Context) error {
		return c.File("../static/payment.html")
	})

	e.POST("/payment/request", handler.CreatePaymentRequest)
	e.POST("/payment/mock", handler.UpdatePaymentRequest)

}

func main() {

	e := echo.New()

	app := fx.New(
		fx.Provide(
			// NewLogger,
			NewValidator,
			NewDB,
			payment.NewPaymentRepository,
			payment.NewPaymentService,
			payment.NewPaymentHandler,
			func() *echo.Echo {
				return e
			},
		),
		fx.Invoke(
			func(e *echo.Echo, handler *payment.PaymentHandler) {
				RegisteredRoutes(e, handler)
			},
			func() {
				if err := e.Start(":8083"); err != nil {
					log.Fatal("Echo server failed to start", zap.Error(err))
				}
			},
		),
	)
	app.Run()
}
