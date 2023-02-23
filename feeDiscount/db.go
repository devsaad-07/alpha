package feediscount

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

type UserMetrics struct {
	user_id                     int
	user_custody_wallet_id      int
	user_uuid                   string
	total_trade_count           int
	total_buy_count             int
	total_sell_count            int
	total_trade_gmv             sql.NullFloat64
	total_buy_gmv               sql.NullFloat64
	total_sell_gmv              sql.NullFloat64
	first_trade_date            string
	last_trade_date             string
	total_inr_deposit_count     sql.NullInt64
	total_inr_deposit_amount    sql.NullFloat64
	total_inr_withdrawal_count  sql.NullInt64
	total_inr_withdrawal_amount sql.NullFloat64
	created_at                  string
	updated_at                  string
}

func Init_db() {
	// Connect to the database
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5433/cs_india?sslmode=disable")
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database!")
}

func GetDb() *sql.DB {
	return db
}

func GetUser(userId int64) UserMetrics {
	var resp UserMetrics

	stmt, err := db.Prepare("SELECT * FROM user_consolidated_metrics WHERE user_id = $1")
	if err != nil {
		panic(err)
	}

	// Execute the SQL statement and retrieve the row

	err = stmt.QueryRow(userId).Scan(
		&resp.user_id,
		&resp.user_custody_wallet_id,
		&resp.user_uuid,
		&resp.total_trade_count,
		&resp.total_buy_count,
		&resp.total_sell_count,
		&resp.total_trade_gmv,
		&resp.total_buy_gmv,
		&resp.total_sell_gmv,
		&resp.first_trade_date,
		&resp.last_trade_date,
		&resp.total_inr_deposit_count,
		&resp.total_inr_deposit_amount,
		&resp.total_inr_withdrawal_count,
		&resp.total_inr_withdrawal_amount,
		&resp.created_at,
		&resp.updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Print(resp)
	return resp
}
