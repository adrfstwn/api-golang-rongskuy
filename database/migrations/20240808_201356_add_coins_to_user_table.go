package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddCoinsToUserTable_20240808_201356 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddCoinsToUserTable_20240808_201356{}
	m.Created = "20240808_201356"

	migration.Register("AddCoinsToUserTable_20240808_201356", m)
}

// Run the migrations
func (m *AddCoinsToUserTable_20240808_201356) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE user ADD COLUMN coins INT(255) DEFAULT '0'")
	m.SQL("ALTER TABLE user ADD COLUMN email_verified_at DATETIME NULL")

}

// Reverse the migrations
func (m *AddCoinsToUserTable_20240808_201356) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE user DROP COLUMN coins")
	m.SQL("ALTER TABLE user DROP COLUMN email_verified_at")

}
