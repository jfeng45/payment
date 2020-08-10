package app

import (
	"database/sql"
	ycq "github.com/jetbasrawi/go.cqrs"
	logConfig "github.com/jfeng45/glogger/config"
	logFactory "github.com/jfeng45/glogger/factory"
	"github.com/jfeng45/gmessaging"
	gmessagingConfig "github.com/jfeng45/gmessaging/config"
	gmessagingFactory "github.com/jfeng45/gmessaging/factory"
	"github.com/jfeng45/payment/app/config"
	"github.com/jfeng45/payment/app/container"
	"github.com/jfeng45/payment/app/container/containerhelper"
	"github.com/jfeng45/payment/app/container/servicecontainer"
	"github.com/jfeng45/payment/app/logger"
	"github.com/jfeng45/payment/domain/command"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// InitApp loads the application configurations from a file and saved it in appConfig and initialize the logger
// The appConfig is cached in container, so it only loads the configuration file once.
// InitApp only needs to be called once. If the configuration changes, you can call it again to reinitialize the app.
func InitApp(filename...string) (container.Container, error) {
	config, err := config.BuildConfig(filename...)
	if err != nil {
		return nil, errors.Wrap(err, "loadConfig")
	}
	err = initLogger(&config.LogConfig)
	if err != nil {
		return nil, err
	}
	return initContainer(config)
}

func initContainer(config *config.AppConfig) (container.Container, error) {
	factoryMap := make(map[string]interface{})
	c := servicecontainer.ServiceContainer{factoryMap,config}
	gdbc, err :=initGdbc(&c.AppConfig.SQLConfig)
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
	loadDispatcher(&c)
	return &c, nil
}

func initLogger (lc *logConfig.Logging) error{
	log, err := logFactory.Build(lc)
	if err != nil {
		return errors.Wrap(err, "loadLogger")
	}
	logger.SetLogger(log)
	return nil
}

func initGdbc(dsc *config.DataStoreConfig) (*sql.DB,error) {

	db, err := sql.Open(dsc.DriverName, dsc.UrlAddress)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	// check the connection
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	//dt := databasehandler.SqlDBTx{DB: db}
	return db, nil
}

func initMessagingService() (gmessaging.MessagingEncodedInterface, error) {
	config := gmessagingConfig.Messaging{gmessagingConfig.CODE_NATS, config.MESSAGING_SERVER_URL, nats.JSON_ENCODER}
	return gmessagingFactory.BuildEncoded(&config)
}
func initDispatcher() ycq.Dispatcher {
	// Create the dispatcher
	dispatcher := ycq.NewInMemoryDispatcher()
	return dispatcher
}

func loadDispatcher(c container.Container) error {
	var value interface{}
	var found bool
	if value, found = c.Get(container.DISPATCHER); !found {
		//logger.Log.Debug("find CacheGrpc key=%v \n", key)
		message := "can't find key=" + container.DISPATCHER + " in containier "
		return errors.New(message)
	}
	d := value.(ycq.Dispatcher)
	mpuc, err := containerhelper.BuildMakePaymentUseCase(c)
	if err != nil {
		return err
	}
	cpch := command.MakePaymentCommandHandler{mpuc}
	d.RegisterHandler(&cpch,&command.MakePaymentCommand{})
	return nil
}



