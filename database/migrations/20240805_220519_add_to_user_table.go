package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddToUserTable_20240805_220519 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddToUserTable_20240805_220519{}
	m.Created = "20240805_220519"

	migration.Register("AddToUserTable_20240805_220519", m)
}

func (m *AddToUserTable_20240805_220519) Up() {
	// Menambahkan kolom reset_token dan reset_expiry
	m.SQL("ALTER TABLE user ADD COLUMN reset_token VARCHAR(255) DEFAULT NULL")
	m.SQL("ALTER TABLE user ADD COLUMN reset_expiry DATETIME DEFAULT NULL")
}

// Reverse the migrations
func (m *AddToUserTable_20240805_220519) Down() {
	// Menghapus kolom reset_token dan reset_expiry
	m.SQL("ALTER TABLE user DROP COLUMN reset_token")
	m.SQL("ALTER TABLE user DROP COLUMN reset_expiry")
}
