package app

import (
	"database/sql"
	ycq "github.com/jetbasrawi/go.cqrs"
	glogger "github.com/jfeng45/glogger/logfactory"
	"github.com/jfeng45/gmessaging"
	"github.com/jfeng45/gmessaging/nat"
	"github.com/jfeng45/payment/app/config"
	"github.com/jfeng45/payment/app/container"
	"github.com/jfeng45/payment/app/container/containerhelper"
	"github.com/jfeng45/payment/app/container/servicecontainer"
	"github.com/jfeng45/payment/app/logger"
	"github.com/jfeng45/payment/domain/command"
	"github.com/jfeng45/payment/tool/gdbc"
	"github.com/jfeng45/payment/tool/gdbc/databasehandler"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"log"
)

// InitApp loads the application configurations from a file and saved it in appConfig and initialize the logger
// The appConfig is cached in container, so it only loads the configuration file once.
// InitApp only needs to be called once. If the configuration changes, you can call it again to reinitialize the app.
func InitApp(filename...string) (container.Container, error) {
	config, err := config.BuildConfig(filename...)
	if err != nil {
		return nil, errors.Wrap(err, "loadConfig")
	}
	err = initLogger(config)
	if err != nil {
		return nil, err
	}
	return initContainer(config)
}

func initContainer(config *config.AppConfig) (container.Container, error) {
	factoryMap := make(map[string]interface{})
	c := servicecontainer.ServiceContainer{factoryMap,config}
	gdbc, err :=initGdbc(&c)
	if err != nil {
		return nil,err
	}
	c.Put(container.DATABASE, gdbc)
	ec, err := initMessagingService()
	if err != nil {
		return nil, err
	}
	c.Put(container.MESSAGING_SERVER, ec)
	d := initDispatcher()
	c.Put(container.DISPATCHER, d)
	loadDispatcher(c)
	return &c, nil
}

func initLogger (config *config.AppConfig) error{
	log, err := glogger.InitLogger(config.Log)
	if err != nil {
		return errors.Wrap(err, "loadLogger")
	}
	logger.SetLogger(log)
	return nil
}

func initGdbc(sc *servicecontainer.ServiceContainer) (gdbc.SqlGdbc,error) {

	db, err := sql.Open(sc.AppConfig.SQLConfig.DriverName, sc.AppConfig.SQLConfig.UrlAddress)
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

func initMessagingService() (gmessaging.MessagingInterface, error) {
	url := config.MESSAGING_SERVER_URL
	nc, err :=nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	//defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	nat := nat.Nat{ec}
	return nat, nil
	//defer ec.Close()
}
func initDispatcher() ycq.Dispatcher {
	// Create the dispatcher
	dispatcher := ycq.NewInMemoryDispatcher()
	return dispatcher
}

func loadDispatcher(c servicecontainer.ServiceContainer) error {
	var value interface{}
	var found bool
	if value, found = c.Get(container.DISPATCHER); !found {
		//logger.Log.Debug("find CacheGrpc key=%v \n", key)
		message := "can't find key=" + container.DISPATCHER + " in containier "
		return errors.New(message)
	}
	d := value.(ycq.Dispatcher)
	mpuc, err := containerhelper.BuildMakePaymentUseCase(&c)
	if err != nil {
		return err
	}
	cpch := command.MakePaymentCommandHandler{mpuc}
	d.RegisterHandler(&cpch,&command.MakePaymentCommand{})
	return nil
}



