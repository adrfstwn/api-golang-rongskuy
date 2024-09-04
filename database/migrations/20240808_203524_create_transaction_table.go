package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateTransactionTable_20240808_203524 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTransactionTable_20240808_203524{}
	m.Created = "20240808_203524"

	migration.Register("CreateTransactionTable_20240808_203524", m)
}

// Run the migrations
func (m *CreateTransactionTable_20240808_203524) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
		CREATE TABLE transaction (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			trash_id INT,
			coins_earned INT,
			status VARCHAR(50) DEFAULT 'pending',
			kode_transaksi VARCHAR(255)
			keterangan VARCHAR(255)
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id),
			FOREIGN KEY (trash_id) REFERENCES trash(id)
		)
	`)

}

// Reverse the migrations
func (m *CreateTransactionTable_20240808_203524) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE transaction")

}
