package containerhelper

import (
	"github.com/jfeng45/payment/app/container"
	"github.com/jfeng45/payment/domain/usecase"
	"github.com/pkg/errors"
)

func BuildMakePaymentUseCase(c container.Container) (usecase.MakePaymentUseCaseInterface, error) {
	key := container.USECASE_CREATE_PAYMENT
	value, err := c.BuildUseCase(key)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return value.(usecase.MakePaymentUseCaseInterface), nil
}

func BuildGetPaymentUseCase(c container.Container) (usecase.GetPaymentUseCaseInterface, error) {
	key := container.USECASE_GET_PAYMENT
	value, err := c.BuildUseCase(key)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return value.(usecase.GetPaymentUseCaseInterface), nil
}

