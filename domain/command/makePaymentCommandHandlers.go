package command

import (
	ycq "github.com/jetbasrawi/go.cqrs"
	"github.com/jfeng45/payment/app/logger"
	"github.com/jfeng45/payment/domain/usecase"
)

type MakePaymentCommandHandler struct {
	Mpuci usecase.MakePaymentUseCaseInterface
}

func (mpch *MakePaymentCommandHandler) Handle (message ycq.CommandMessage) error {
	switch cmd := message.Command().(type) {
	case *MakePaymentCommand:
		payment := cmd.NewPayment()
		_, err := mpch.Mpuci.MakePayment(&payment)
		if err != nil {
			logger.Log.Errorf("error in MakePaymentCommandHandler:", err)
		}
	default:
		logger.Log.Errorf("command type mismatch in MakePaymentCommandHandler:")
	}
	return nil
}