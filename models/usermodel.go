package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/beego/beego/v2/client/orm"
)

type UserModel struct {
	Id              int64     `orm:"auto"`
	Name            string    `orm:"size(128)"`
	Email           string    `orm:"size(128)"`
	Password        string    `orm:"size(128)"`
	EmailVerifiedAt time.Time `orm:"type(datetime);null"`
	Role            string    `orm:"size(255)"`
	Coins           int64     `orm:"default(0)"`
	CreatedAt       time.Time `orm:"type(datetime);auto_now_add"`
	UpdatedAt       time.Time `orm:"type(datetime);auto_now"`
	ResetToken      string    `json:"reset_token"`
	ResetExpiry     time.Time `json:"reset_expiry"`
}

func (u *UserModel) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(UserModel))
}

// UpdateUserCoins updates the coins for a user. Function signature updated to accept int64.
func UpdateUserCoins(userID int64, coins int64) error {
	o := orm.NewOrm()
	user := &UserModel{Id: userID}
	if err := o.Read(user); err != nil {
		return err
	}
	user.Coins += coins
	if _, err := o.Update(user); err != nil {
		return err
	}
	return nil
}

// generateRandomToken generates a random token using UUID
func generateRandomToken() string {
	return uuid.NewString()
}

// GetUserModelByName retrieves a user by their name
func GetUserModelByName(name string) (*UserModel, error) {
	o := orm.NewOrm()
	var user UserModel
	err := o.QueryTable("user").Filter("Name", name).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserModelByEmail retrieves a user by their email
func GetUserModelByEmail(email string) (*UserModel, error) {
	o := orm.NewOrm()
	var user UserModel
	err := o.QueryTable("user").Filter("Email", email).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AddUserModel insert a new UserModel into database and returns
// last inserted Id on success.
func AddUserModel(m *UserModel) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserModelById retrieves UserModel by Id.
func GetUserModelById(id int64) (v *UserModel, err error) {
	o := orm.NewOrm()
	v = &UserModel{Id: id}
	if err = o.QueryTable(new(UserModel)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserModel retrieves all UserModel matches certain condition. Returns empty list if
// no records exist
func GetAllUserModel(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserModel))
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
					return nil, errors.New("error: Invalid order. Must be either [asc|desc]")
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
					return nil, errors.New("reror: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("error: unused 'order' fields")
		}
	}

	var l []UserModel
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

// UpdateUserModel updates UserModel by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserModelById(m *UserModel) (err error) {
	o := orm.NewOrm()
	v := UserModel{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserModel deletes UserModel by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserModel(id int64) (err error) {
	o := orm.NewOrm()
	v := UserModel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserModel{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// GenerateResetToken generates a reset token and saves it in the UserModel
func (u *UserModel) GenerateResetToken() (string, error) {
	token := generateRandomToken()
	expiry := time.Now().Add(1 * time.Hour) // Token valid for 1 hour

	u.ResetToken = token
	u.ResetExpiry = expiry

	// Create a copy of the user model to update in the database
	_, err := orm.NewOrm().Update(&UserModel{
		Id:          u.Id,
		ResetToken:  u.ResetToken,
		ResetExpiry: u.ResetExpiry,
	}, "ResetToken", "ResetExpiry")
	if err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword validates the reset token and updates the password
func (u *UserModel) ResetPassword(newPassword string) error {
	if time.Now().After(u.ResetExpiry) {
		return errors.New("reset token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.ResetToken = ""
	u.ResetExpiry = time.Time{}

	_, err = orm.NewOrm().Update(&UserModel{
		Id:          u.Id,
		Password:    u.Password,
		ResetToken:  u.ResetToken,
		ResetExpiry: u.ResetExpiry,
	}, "Password", "ResetToken", "ResetExpiry")
	if err != nil {
		return err
	}

	return nil
}

// GetUserModelByResetToken fetches a user by reset token from the database
func GetUserModelByResetToken(token string) (*UserModel, error) {
	o := orm.NewOrm()
	user := &UserModel{}
	err := o.QueryTable("user").Filter("ResetToken", token).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
