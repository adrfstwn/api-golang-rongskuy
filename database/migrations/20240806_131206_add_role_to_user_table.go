package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AddRoleToUserTable_20240806_131206 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddRoleToUserTable_20240806_131206{}
	m.Created = "20240806_131206"

	migration.Register("AddRoleToUserTable_20240806_131206", m)
}

// Run the migrations
func (m *AddRoleToUserTable_20240806_131206) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE user ADD COLUMN role VARCHAR(255) DEFAULT 'user'")
}

// Reverse the migrations
func (m *AddRoleToUserTable_20240806_131206) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE user DROP COLUMN role")

}
