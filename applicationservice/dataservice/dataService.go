package dataservice

import "github.com/jfeng45/payment/domain/model"

type PaymentDataInterface interface {
	Insert(user *model.Payment) (resultUser *model.Payment, err error)
	Find(id int) (*model.Payment, error)
}
