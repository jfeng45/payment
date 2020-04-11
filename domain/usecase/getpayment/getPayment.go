package getpayment

import (
	"github.com/jfeng45/payment/applicationservice/dataservice"
	"github.com/jfeng45/payment/domain/model"
)

type GetPaymentUseCase struct {
	PaymentDataInterface dataservice.PaymentDataInterface
	//EventBus              ycq.EventBus
}

func (gpuc *GetPaymentUseCase) GetPayment(id int) (*model.Payment, error) {
	return gpuc.PaymentDataInterface.Find(id)
}

