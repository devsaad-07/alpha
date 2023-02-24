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
	Name     string `gorm:"column:name"`
	IsActive bool   `gorm:"column:is_active"`
}

type UserConsolidatedMetrics struct {
	gorm.Model
	UserId                   int             `gorm:"column:user_id"`
	UserCustodyWalletId      int             `gorm:"column:user_custody_wallet_id"`
	UserUuid                 string          `gorm:"column:user_uuid"`
	TotalTradeCount          int             `gorm:"column:total_trade_count"`
	TotalBuyCount            int             `gorm:"column:total_buy_count"`
	TotalSellCount           int             `gorm:"column:total_sell_count"`
	TotalTradeGmv            sql.NullFloat64 `gorm:"column:total_trade_gmv"`
	TotalBuyGmv              sql.NullFloat64 `gorm:"column:total_buy_gmv"`
	TotalSellGmv             sql.NullFloat64 `gorm:"column:total_sell_gmv"`
	FirstTradeDate           string          `gorm:"column:first_trade_date"`
	LastTradeDate            string          `gorm:"column:last_trade_date"`
	TotalInrDepositCount     sql.NullInt64   `gorm:"column:total_inr_deposit_count"`
	TotalInrDepositAmount    sql.NullFloat64 `gorm:"column:total_inr_deposit_amount"`
	TotalInrWithdrawalCount  sql.NullInt64   `gorm:"column:total_inr_withdrawal_count"`
	TotalInrWithdrawalAmount sql.NullFloat64 `gorm:"column:total_inr_withdrawal_amount"`
}

func GetUser(userId int64) UserConsolidatedMetrics {
	var resp UserConsolidatedMetrics

	db := GetDb()
	tx := db.Where("user_id = ?", userId).First(&resp)
	if tx.Error != nil {
		return resp
	}
	return resp
}

func GetAllRules(ruleType string) (rule []Rules) {
	db := GetDb()
	tx := db.Where("type = ?", ruleType).Find(&rule)
	if tx.Error != nil {
		return
	}
	return rule
}

func SaveRule(rule Rules) (err error) {
	db := GetDb()
	tx := db.Create(&rule)
	return tx.Error
}

func UpdateRuleStatus(id int, status bool) {
	db := GetDb()
	tx := db.Model(&Rules{}).Where("id = ?", id).Update("is_active", status)
	if tx.Error != nil {
		return
	}
}

func GetRuleById(id int) Rules {
	var rule Rules
	db := GetDb()
	tx := db.Where("id = ?", id).Find(&rule)
	if tx.Error != nil {
		return rule
	}
	return rule
}
