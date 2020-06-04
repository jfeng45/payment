// Package config reasd configurations from a YAML file and load them into a AppConfig type to save the configuration
// information for the application.
// Configuration for different environment can be saved in files with different suffix, for example [Dev], [Prod]
package config
import (
	"fmt"
	logConfig "github.com/jfeng45/glogger/config"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)
// use case code. Need to map to the use case code (UseCaseConfig) in the configuration yaml file.
// Client app use those to retrieve use case from the container
const (
	LOG_ENABLE_CALLER bool   = true

	DB_CONFIG_CODE           = "sql"
	DB_DRIVER_CODE    string = "sqldb"
	DB_DRIVER_NAME    string = "mysql"
	DB_NAME                  = ""
	DB_SOURCE_NAME    string ="root:@tcp(localhost:4333)/service_config?charset=utf8"

	MESSAGING_SERVER_URL string= nats.DefaultURL

	SUBJECT_PAYMENT_CREATED string ="payment.paymentCreated"
	SUBJECT_PAYMENT_CANCELED string ="Payment.paymentCanceled"

	SUBJECT_MAKE_PAYMENT string ="payment.makePayment"
)

// AppConfig represents the application config
type AppConfig struct {
	SQLConfig     DataStoreConfig   `yaml:"sqlConfig"`
	ZapConfig     logConfig.Logging `yaml:"zapConfig"`
	LorusConfig   logConfig.Logging `yaml:"logrusConfig"`
	LogConfig    logConfig.Logging `yaml:"logConfig"`
	UseCaseConfig UseCaseConfig     `yaml:"useCaseConfig"`
}

// UseCaseConfig represents different use cases
type UseCaseConfig struct {
	Payment PaymentConfig  `yaml:"getPayment"`
}

// RegistrationConfig represents registration use case
type PaymentConfig struct {
	//Code           string     `yaml:"code"`
	PaymentDataConfig DataConfig `yaml:"userDataConfig"`
}
//// DataConfig represents data service
type DataConfig struct {
	Code            string          `yaml:"code"`
	DataStoreConfig DataStoreConfig `yaml:"dataStoreConfig"`
}

// DataConfig represents handlers for data store. It can be a database or a gRPC connection
type DataStoreConfig struct {
	Code string `yaml:"code"`
	// Only database has a driver name, for grpc it is "tcp" ( network) for server
	DriverName string `yaml:"driverName"`
	// For database, this is datasource name; for grpc, it is target url
	UrlAddress string `yaml:"urlAddress"`
	// Only some databases need this database name
	DbName string `yaml:"dbName"`
}
// BuildConfig build the AppConfig
// if the filaname is not empty, then it reads the file of the filename (in the same folder) and put it into the AppConfig
func BuildConfig(filename ...string) (*AppConfig, error) {
	if len(filename) == 1 {
		return buildConfigFromFile(filename[0])
	} else {
		return BuildConfigWithoutFile()
	}
}

// BuildConfigWithoutFile create AppConfig with adhoc value
func BuildConfigWithoutFile() (*AppConfig, error) {

	dsc := DataStoreConfig{DB_DRIVER_CODE, DB_DRIVER_NAME,
		DB_SOURCE_NAME, DB_NAME}
	lc := logConfig.Logging{logConfig.ZAP, logConfig.DEBUG, LOG_ENABLE_CALLER}
	dc := DataConfig{DB_CONFIG_CODE, dsc}
	pc := PaymentConfig{dc}
	ucc := UseCaseConfig{pc}
	//ac := AppConfig{dsc, lc,lc,lc, "GetPaymentUseCase"}
	ac := AppConfig{dsc, lc,lc,lc, ucc}
	fmt.Printf("appConfig:", ac)
	return &ac, nil

}

// buildConfigFromFile reads the file of the filename (in the same folder) and put it into the AppConfig
func buildConfigFromFile(filename string) (*AppConfig, error) {

	var ac AppConfig
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read error")
	}
	err = yaml.Unmarshal(file, &ac)

	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	//err = validateConfig(ac)
	//if err != nil {
	//	return nil, errors.Wrap(err, "validate config")
	//}
	fmt.Println("appConfig:", ac)
	return &ac, nil
}

