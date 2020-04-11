package makepayment

import (
	"github.com/jfeng45/payment/app/config"
	"github.com/jfeng45/payment/applicationservice/dataservice"
	"github.com/jfeng45/gmessaging"
	"github.com/jfeng45/payment/domain/event"
	"github.com/jfeng45/payment/domain/model"
	"github.com/pkg/errors"
)

type MakePaymentUseCase struct {
	PaymentDataInterface dataservice.PaymentDataInterface
	Mi                   gmessaging.MessagingInterface
}

func (mpu *MakePaymentUseCase) MakePayment(payment *model.Payment) (*model.Payment, error) {
	payment, err := mpu.PaymentDataInterface.Insert(payment)
	if err!= nil {
		return nil, errors.Wrap(err, "")
	}
	pce := event.NewPaymentCreatedEvent(*payment)
	err = mpu.Mi.Publish(config.SUBJECT_PAYMENT_CREATED, pce)
	if err != nil {
		return nil, err
	}
	return payment, nil
}


