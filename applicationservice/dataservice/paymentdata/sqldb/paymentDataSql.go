// Package sql represents SQL database implementation of the user data persistence layer
package sqldb

import (
"database/sql" 
_ "github.com/go-sql-driver/mysql"
"github.com/jfeng45/payment/app/logger"
"github.com/jfeng45/payment/domain/model"
"github.com/jfeng45/payment/tool"
"github.com/jfeng45/payment/tool/gdbc"
"github.com/pkg/errors"
"time"
)

const (
	//QUERY_PAYMENT_BY_ID   string = "SELECT * FROM payment where id =?"
	QUERY_PAYMENT_BY_ID   string = "SELECT id, payment_method, status, created_time FROM payment where id =?"
	QUERY_PAYMENT_BY_ORDER_NUMBER        = "SELECT * FROM payment where order_number =?"
	QUERY_PAYMENT                = "SELECT * FROM payment "
	//UPDATE_PAYMENT               = "update userinfo set username=?, department=?, created=? where uid=?"
	INSERT_PAYMENT               = "INSERT payment SET source_account=?,target_account=?, amount=?, payment_method=?," +
		"status=?, order_number =?,created_time=?, completion_time=?"
)

// PaymentDataSql is the SQL implementation of PaymentDataInterface
type PaymentDataSql struct {
	DB gdbc.SqlGdbc
}

func (uds *PaymentDataSql) Find(id int) (*model.Payment, error) {
	rows, err := uds.DB.Query(QUERY_PAYMENT_BY_ID, id)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	return retrievePayment(rows)
}
func retrievePayment(rows *sql.Rows) (*model.Payment, error) {
	if rows.Next() {
		return rowsToPayment(rows)
	}
	return nil, nil
}
func rowsToPayment(rows *sql.Rows) (*model.Payment, error) {
	var ds string
	payment := &model.Payment{}
	err := rows.Scan(&payment.Id, &payment.PaymentMethod, &payment.Status, &ds)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	created, err := time.Parse(tool.FORMAT_ISO8601_DATE, ds)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	payment.CreatedTime = created

	//logger.Log.Debug("rows to Payment:", *payment)
	return payment, nil
}
func (uds *PaymentDataSql) FindByOrderNumber(orderNumber string) (*model.Payment, error) {
	rows, err := uds.DB.Query(QUERY_PAYMENT_BY_ORDER_NUMBER, orderNumber)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	return retrievePayment(rows)
}

func (uds *PaymentDataSql) FindAll() ([]model.Payment, error) {

	rows, err := uds.DB.Query(QUERY_PAYMENT)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	users := []model.Payment{}

	//var ds string
	for rows.Next() {
		user, err := rowsToPayment(rows)
		if err != nil {
			return users, errors.Wrap(err, "")
		}
		users = append(users, *user)

	}
	//need to check error for rows.Next()
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "")
	}
	logger.Log.Debug("find user list:", users)
	return users, nil
}

func (uds *PaymentDataSql) Insert(p *model.Payment) (*model.Payment, error) {

	stmt, err := uds.DB.Prepare(INSERT_PAYMENT)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer stmt.Close()
	res, err := stmt.Exec(p.SourceAccount,p.TargetAccount,p.Amount, p.PaymentMethod, p.Status,
		p.OrderNumber, p.CreatedTime, p.CompletionTime)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	p.Id = int(id)
	//logger.Log.Debug("user inserted:", p)
	return p, nil
}

//func (uds *PaymentDataSql) EnableTx(tx dataservice.TxDataInterface) {
//	uds.DB = tx.GetTx()
//}
