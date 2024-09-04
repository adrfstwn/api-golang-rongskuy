package routers

import (
	"backend-rongskuy/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	// Authentication
	beego.Router("/register", &controllers.AuthController{}, "post:Register")
	beego.Router("/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/logout", &controllers.AuthController{}, "post:Logout")

	// Forgot Password
	beego.Router("/forgot-password", &controllers.AuthController{}, "post:ForgotPassword")
	beego.Router("/reset-password", &controllers.AuthController{}, "post:ResetPassword")

	// Admin Page Dashboard

	// Admin Page Postingan
	beego.Router("/admin/trash", &controllers.TrashController{}, "get:ReadAll")
	beego.Router("/admin/buy", &controllers.TransactionController{}, "post:BuyTrash")

	// Admin Page Riwayat Transaksi
	beego.Router("/admin/history", &controllers.TransactionController{}, "get:ReadAdmin")

	// Admin Page Users Edit Delete
	beego.Router("/admin/users/:id", &controllers.AdminController{}, "put:UserEdit")
	beego.Router("/admin/users/:id", &controllers.AdminController{}, "delete:UserDelete")
	beego.Router("/admin/users/:id", &controllers.AdminController{}, "get:UserGetOne")
	beego.Router("/admin/users", &controllers.AdminController{}, "get:UserGetAll")

	// User Page Dashboard

	// User Page Posting Trash
	beego.Router("/user/trash", &controllers.TrashController{}, "post:Create")
	beego.Router("/user/trash/:id", &controllers.TrashController{}, "put:Edit")
	beego.Router("/user/trash/:id", &controllers.TrashController{}, "delete:Delete")
	beego.Router("/user/trash", &controllers.TrashController{}, "get:Read")

	// User Page Riwayat Transaksi
	beego.Router("/user/history", &controllers.TransactionController{}, "get:Read")
}
