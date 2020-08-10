package servicecontainer

import (
	"database/sql"
	"github.com/jfeng45/payment/app/config"
	"github.com/jfeng45/payment/app/container"
	"github.com/jfeng45/payment/app/logger"
	"github.com/jfeng45/payment/applicationservice/dataservice/paymentdata/sqldb"
	"github.com/jfeng45/gmessaging"
	"github.com/jfeng45/payment/domain/usecase/getpayment"
	"github.com/jfeng45/payment/domain/usecase/makepayment"
	"github.com/jfeng45/payment/tool/gdbc"
	"github.com/jfeng45/payment/tool/gdbc/databasehandler"
	"github.com/pkg/errors"
)

type ServiceContainer struct {
	FactoryMap map[string]interface{}
	AppConfig  *config.AppConfig
}

func (sc *ServiceContainer) BuildUseCase(code string) (interface{}, error) {
	dt, err := buildGdbc(sc)
	if err != nil {
		return nil, err
	}
	pds := sqldb.PaymentDataSql{dt}
	var value interface{}
	var found bool
	if value, found = sc.Get(container.MESSAGING_SERVER); !found {
		message := "can't find key= in containier " + container.EVENT_BUS
		logger.Log.Errorf(message)
	}
	ms := value.(gmessaging.MessagingEncodedInterface)
	switch code {
		case container.USECASE_GET_PAYMENT:
			uc := getpayment.GetPaymentUseCase{&pds}
			logger.Log.Debug("found db in container for key:", code)
			return &uc, nil
		case container.USECASE_CREATE_PAYMENT:
			uc := makepayment.MakePaymentUseCase{&pds, ms}
			logger.Log.Debug("found db in container for key:", code)
			return &uc, nil
		}
	return nil, nil
}

func buildGdbc(sc *ServiceContainer) (gdbc.SqlGdbc,error) {
	sqlc := sc.AppConfig.SQLConfig
	key := sqlc.Code
	//if it is already in container, return
	if value, found := sc.Get(key); found {
		//logger.Log.Debugf("found db value %+v\n:", value)
		sdb := value.(*sql.DB)
		sdt := databasehandler.SqlDBTx{DB: sdb}
		logger.Log.Debug("found db in container for key:", key)
		return &sdt, nil
	}

	db, err := sql.Open(sqlc.DriverName, sqlc.UrlAddress)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	// check the connection
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	dt := databasehandler.SqlDBTx{DB: db}
	return &dt, nil
}

func (sc *ServiceContainer) Get(code string) (interface{}, bool) {
	value, found := sc.FactoryMap[code]
	return value, found
}

func (sc *ServiceContainer) Put(code string, value interface{}) {
	sc.FactoryMap[code] = value
}



