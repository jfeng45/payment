package command

import (
	ycq "github.com/jetbasrawi/go.cqrs"
	"github.com/jfeng45/payment/domain/model"
	"strconv"
	"time"
)

type MakePaymentCommand struct {
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

func (pc *MakePaymentCommand) NewPayment() *model.Payment{
	return &model.Payment{pc.Id, pc.SourceAccount, pc.TargetAccount, pc.Amount,
		pc.Status, pc.PaymentMethod, pc.OrderNumber, pc.CreatedTime,
		pc.CompletionTime}

}

func (pc *MakePaymentCommand) NewMakePaymentDescriptor() *ycq.CommandDescriptor {
	aggregateId := strconv.Itoa(pc.Id)
	return ycq.NewCommandMessage(aggregateId, pc)
}