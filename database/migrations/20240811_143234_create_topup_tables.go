package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateTopupTables_20240811_143234 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTopupTables_20240811_143234{}
	m.Created = "20240811_143234"

	migration.Register("CreateTopupTables_20240811_143234", m)
}

// Run the migrations
func (m *CreateTopupTables_20240811_143234) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
	CREATE TABLE topup (
		id varchar(100) PRIMARY KEY,
		user_id INT,
		status VARCHAR(150) NOT NULL,
		amount DECIMAL(15, 2) NOT NULL,
		snap_url VARCHAR(255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id)
	)
`)
}

// Reverse the migrations
func (m *CreateTopupTables_20240811_143234) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE topup")
}
