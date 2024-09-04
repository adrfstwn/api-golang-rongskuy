package models

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Transaction struct {
	Id            int64     `orm:"auto"`                        // Primary key
	UserID        int64     `orm:"column(user_id)"`             // Foreign key to User
	TrashID       int64     `orm:"column(trash_id)"`            // Foreign key to Trash
	CoinsEarned   int64     `orm:"default(0)"`                  // Coins earned from transaction
	Status        string    `orm:"size(50);default(pending)"`   // Status of the transaction
	KodeTransaksi string    `orm:"size(255)"`                   // Transaction code
	Keterangan    string    `orm:"size(255)"`                   // Description
	CreatedAt     time.Time `orm:"auto_now_add;type(datetime)"` // Creation timestamp
	UpdatedAt     time.Time `orm:"auto_now;type(datetime)"`     // Update timestamp
}

func (t *Transaction) TableName() string {
	return "transaction"
}

func init() {
	orm.RegisterModel(new(Transaction))
}

// AddTransaction inserts a new Transaction into database and returns the last inserted Id on success.
func AddTransaction(m *Transaction) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetTransactionsByTrashID(trashID int64) (ml []Transaction, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transaction))

	// Filter transactions by UserID
	qs = qs.Filter("trash_id", trashID)

	// Retrieve all transactions
	var transactions []Transaction
	if _, err = qs.All(&transactions); err == nil {
		return transactions, nil
	}
	return nil, err
}
func GetTransactionsByUserID(trashID int64) (ml []Transaction, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transaction))

	// Filter transactions by UserID
	qs = qs.Filter("user_id", trashID)

	// Retrieve all transactions
	var transactions []Transaction
	if _, err = qs.All(&transactions); err == nil {
		return transactions, nil
	}
	return nil, err
}

// GetTransactionById retrieves Transaction by Id. Returns error if Id doesn't exist.
func GetTransactionById(id int64) (v *Transaction, err error) {
	o := orm.NewOrm()
	v = &Transaction{Id: id}
	if err = o.QueryTable(new(Transaction)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTransaction retrieves all Transactions matching certain conditions. Returns an empty list if no records exist.
func GetAllTransaction(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transaction))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("invalid order: must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("invalid order: must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("error: sortby, order sizes mismatch or order size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("error: unused 'order' fields")
		}
	}

	var l []Transaction
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
