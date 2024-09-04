package models

import (
	"fmt"
	"log"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Trash model
type Trash struct {
	Id          int64     `orm:"auto"`
	UserID      int64     `orm:"column(user_id)"`
	Gambar      string    `orm:"type(text)"`
	Description string    `orm:"type(text)"`
	Weight      int64     `orm:"type(int)"`
	Latitude    float64   `orm:"type(decimal(10,8))"`
	Longitude   float64   `orm:"type(decimal(11,8))"`
	Whatsapp    string    `orm:"size(20)"`
	Status      string    `orm:"size(50);default('available')"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
}

func (u *Trash) TableName() string {
	return "trash"
}

func init() {
	// Register model
	orm.RegisterModel(new(Trash))
}

// AddTrash insert a new Trash into database and returns
// last inserted Id on success.
func AddTrash(m *Trash) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTrashById retrieves Trash by Id. Returns error if
// Id doesn't exist

func GetTrashById(id int64) (v *Trash, err error) {
	o := orm.NewOrm()
	v = &Trash{Id: id}
	if err = o.QueryTable(new(Trash)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetTrashesByUserID(userID int64) ([]Trash, error) {
	o := orm.NewOrm()
	var trashes []Trash
	_, err := o.QueryTable("trash").Filter("user_id", userID).All(&trashes)
	if err != nil {
		return nil, err
	}
	return trashes, nil
}

// GetAllTrashes retrieves all trash posts from the database
func GetAllTrash() ([]Trash, error) {
	var trashes []Trash
	_, err := orm.NewOrm().QueryTable(new(Trash)).All(&trashes)
	return trashes, err
}

// UpdateTrash updates Trash by Id and returns error if
// the record to be updated doesn't exist
func UpdateTrashById(trash *Trash) error {
	o := orm.NewOrm()
	// Log updated values
	log.Printf("Updating trash with values: %+v\n", trash)
	// Use transaction if needed
	_, err := o.Update(trash, "Description", "Weight", "Latitude", "Longitude", "Whatsapp", "Gambar", "UpdatedAt")
	if err != nil {
		log.Printf("Update error: %v\n", err)
	}
	return err
}

// UpdateTrashStatus updates the status of a Trash record by its ID
func UpdateTrashStatus(id int64, status string) error {
	o := orm.NewOrm()
	// Update status and timestamp
	_, err := o.QueryTable(new(Trash)).Filter("Id", id).Update(orm.Params{
		"Status":    status,
		"UpdatedAt": time.Now(),
	})
	if err != nil {
		log.Printf("UpdateTrashStatus error: %v\n", err)
	}
	return err
}

// DeleteTrash deletes Trash by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTrash(id int64) (err error) {
	o := orm.NewOrm()
	v := Trash{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Trash{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
