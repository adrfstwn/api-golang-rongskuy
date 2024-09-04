package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateUserTable_20240801_201240 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateUserTable_20240801_201240{}
	m.Created = "20240801_201240"

	migration.Register("CreateUserTable_20240801_201240", m)
}

// Run the migrations
func (m *CreateUserTable_20240801_201240) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
		CREATE TABLE user(
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`)

}

// Reverse the migrations
func (m *CreateUserTable_20240801_201240) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL(`DROP TABLE user`)
}
