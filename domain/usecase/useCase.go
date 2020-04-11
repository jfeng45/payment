package usecase

import "github.com/jfeng45/payment/domain/model"

type GetPaymentUseCaseInterface interface {
	GetPayment(id int) (*model.Payment, error)
}
type MakePaymentUseCaseInterface interface {
	MakePayment(payment *model.Payment) (*model.Payment, error)
}
