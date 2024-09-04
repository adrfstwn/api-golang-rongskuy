package controllers

import (
	"backend-rongskuy/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// AdminPageController operations for AdminPage
type AdminController struct {
	beego.Controller
}

func (c *AdminController) UserGetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	v, err := models.GetUserModelById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

func (c *AdminController) UserGetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) == 2 {
				query[kv[0]] = kv[1]
			}
		}
	}

	l, err := models.GetAllUserModel(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

func (c *AdminController) UserEdit() {
	// Get ID from URL
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid or missing ID",
			"details": "The URL must contain a valid ID",
		}
		c.ServeJSON()
		return
	}

	// Read JSON
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to read request body",
		}
		c.ServeJSON()
		return
	}

	// Convert to string for logging
	bodyStr := string(body)
	fmt.Println("Request Body:", bodyStr)

	// Parse JSON
	var updateData struct {
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}
	if err := json.Unmarshal(body, &updateData); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Fetch existing user from database
	user, err := models.GetUserModelById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound) // 404 Not Found
		c.Data["json"] = map[string]string{
			"error":   "User not found",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Merge data
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Password != "" {
		// Validate the password
		if err := validatePassword(updateData.Password); err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
			c.Data["json"] = map[string]string{
				"error":   "Invalid password",
				"details": err.Error(),
			}
			c.ServeJSON()
			return
		}

		// Hash password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	// Update the UserModel in the database
	if err := models.UpdateUserModelById(user); err == nil {
		c.Ctx.Output.SetStatus(http.StatusOK) // 200 OK
		c.Data["json"] = map[string]interface{}{
			"id":      user.Id,
			"success": "User successfully updated",
		}
	} else {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]interface{}{
			"error":   "Failed to update user",
			"details": err.Error(),
		}
	}

	c.ServeJSON()
}

func (c *AdminController) UserDelete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := models.DeleteUserModel(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
