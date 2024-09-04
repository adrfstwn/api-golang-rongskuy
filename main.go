package main

import (
	"backend-rongskuy/middlewares"
	_ "backend-rongskuy/routers"
	"log"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	// Baca konfigurasi dari app.conf
	dbDriver := beego.AppConfig.DefaultString("db_driver", "mysql")
	dbUser, err := beego.AppConfig.String("db_user")
	if err != nil {
		log.Fatalf("Failed to get db_user: %v", err)
	}
	dbPassword, err := beego.AppConfig.String("db_password")
	if err != nil {
		log.Fatalf("Failed to get db_password: %v", err)
	}
	dbName, err := beego.AppConfig.String("db_name")
	if err != nil {
		log.Fatalf("Failed to get db_name: %v", err)
	}
	dbHost := beego.AppConfig.DefaultString("db_host", "127.0.0.1")
	dbPort := beego.AppConfig.DefaultString("db_port", "3306")

	log.Printf("DB Config: driver=%s, user=%s, password=%s, name=%s, host=%s, port=%s", dbDriver, dbUser, dbPassword, dbName, dbHost, dbPort)

	// Buat string koneksi
	connStr := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	// Daftarkan database dengan alias default
	orm.RegisterDataBase("default", dbDriver, connStr)

}

func main() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	// Daftarkan middleware JWT
	beego.InsertFilter("/admin/*", beego.BeforeRouter, middlewares.JWTMiddleware)
	beego.InsertFilter("/user/*", beego.BeforeRouter, middlewares.JWTMiddleware)

	// Run the Beego server
	beego.Run()
}
