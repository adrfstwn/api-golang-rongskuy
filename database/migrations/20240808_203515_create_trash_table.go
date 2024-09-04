package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateTrashTable_20240808_203515 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTrashTable_20240808_203515{}
	m.Created = "20240808_203515"

	migration.Register("CreateTrashTable_20240808_203515", m)
}

// Run the migrations
func (m *CreateTrashTable_20240808_203515) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
		CREATE TABLE trash (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			gambar TEXT,
			description TEXT,
			weight FLOAT,
			latitude DECIMAL(10,8),
			longitude DECIMAL(11,8),
			whatsapp VARCHAR(20),
			status VARCHAR(50) DEFAULT 'available',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id)
		)
	`)

}

// Reverse the migrations
func (m *CreateTrashTable_20240808_203515) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE trash")

}
