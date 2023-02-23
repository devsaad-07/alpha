package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Rules struct {
	gorm.Model
	Type     string `gorm:"column:type"`
	Rule     string `gorm:"column:rule"`
	IsActive bool   `gorm:"column:is_active"`
}

type UserMetrics struct {
	gorm.Model
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

func GetUser(userId int64) UserMetrics {
	var resp UserMetrics

	db := GetDb()
	tx := db.Where("user_id = ?", userId).First(&resp)
	if tx.Error != nil {
		return resp
	}
	return resp
}
