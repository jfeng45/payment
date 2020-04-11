package event

import (
	"github.com/jfeng45/payment/domain/model"
	"time"
)

type PaymentCreatedEvent struct {
	Id int
	SourceAccount int
	TargetAccount int
	Amount float32
	Status model.PaymentStatus
	PaymentMethod model.PaymentMethod
	OrderNumber string
	CreatedTime time.Time
	CompletionTime time.Time
}

func NewPaymentCreatedEvent(p model.Payment) PaymentCreatedEvent{
	pce := PaymentCreatedEvent{p.Id,p.SourceAccount,p.TargetAccount, p.Amount,
		p.Status, p.PaymentMethod, p.OrderNumber, p.CreatedTime,
		p.CompletionTime }
	return pce
}


