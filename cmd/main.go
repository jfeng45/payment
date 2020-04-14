package main

import (
	ycq "github.com/jetbasrawi/go.cqrs"
	"github.com/jfeng45/gmessaging"
	"github.com/jfeng45/payment/app"
	"github.com/jfeng45/payment/app/config"
	"github.com/jfeng45/payment/app/container"
	"github.com/jfeng45/payment/app/container/containerhelper"
	"github.com/jfeng45/payment/app/logger"
	"github.com/jfeng45/payment/domain/command"
	"github.com/jfeng45/payment/domain/model"
	"log"
	"runtime"
	"time"
)

func main() {
	c, err := app.InitApp()
	if err != nil {
		log.Println("err:", err)
	}
	//testMySql(c)
	testSubscribe(c)
}
func testMySql(c container.Container) {
	//testGetPayment(c)
	testMakePayment(c)
}
func testSubscribe(c container.Container) {
	var value interface{}
	var found bool
	if value, found = c.Get(container.MESSAGING_SERVER); !found {
		message := "can't find key=" + container.MESSAGING_SERVER + " in container "
		logger.Log.Errorf(message)
	}
	ms := value.(gmessaging.MessagingInterface)
	if value, found = c.Get(container.DISPATCHER); !found {
		message := "can't find key=" + container.DISPATCHER + " in container "
		logger.Log.Errorf("err:",message)
	}
	d := value.(ycq.Dispatcher)
	subject := config.SUBJECT_MAKE_PAYMENT
	_, err := ms.Subscribe(subject, func(mpc *command.MakePaymentCommand) {
		mpd := mpc.NewMakePaymentDescriptor()
		logger.Log.Debugf("payload:",mpc)
		err := d.Dispatch(mpd)
		if err != nil {
			logger.Log.Errorf("err:",err)
		}
	})
	if err != nil {
		logger.Log.Errorf("err:",err)
	}
	log.Printf("Listening on [%s]", subject)
	runtime.Goexit()
}

func testGetPayment(c container.Container) {
	//It is uid in database. Make sure you have it in database, otherwise it won't find it.
	id := 2
	gpuc, err := containerhelper.BuildGetPaymentUseCase(c)
	if err != nil {
		logger.Log.Fatalf("getPaymentUseCase interface build failed:%+v\n", err)
	}
	user, err := gpuc.GetPayment(id)
	if err != nil {
		logger.Log.Errorf("getPayment failed failed:%+v\n", err)
	}
	logger.Log.Info("find user:", user)
}

func testMakePayment(c container.Container) {

	mpuc, err := containerhelper.BuildMakePaymentUseCase(c)
	if err != nil {
		logger.Log.Fatalf("makePaymentUseCase interface build failed:%+v\n", err)
	}
	//id is auto generated in database
	var id int
	var completionTime time.Time
	p := model.Payment{id,1,32,20,
		model.PAYMENT_STATUS_COMPLETED,model.PAYMENT_METHOD_ALIPAY,"3",
		time.Now(), completionTime}
	payment, err := mpuc.MakePayment(&p)
	if err != nil {
		logger.Log.Errorf("makePaymentUseCase failed failed:%+v\n", err)
	}
	logger.Log.Info("payment created:", payment)
}

