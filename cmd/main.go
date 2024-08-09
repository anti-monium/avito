package main

import (
	api "avito_bootcamp/pkg/apartment_sale_api"
	database "avito_bootcamp/pkg/database"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {

	var db database.IApartmentStorage
	if _, ok := os.LookupEnv("POSTGRES_HOST"); !ok {
		// db = api.NewApartmentStorage()
		log.Info("kek")
	} else {
		var err error
		db, err = database.NewApartmentDatabase()
		if err != nil {
			log.WithError(err).Error("smth goes wrong")
			return
		}
	}
	log.Info(fmt.Sprint(db))

	router := gin.Default()
	server := api.NewSaleServer(db)

	// noAuth
	//router.GET("/dummyLogin/:user_type", server.GetDummyLogin)
	//router.POST("/login", server.PostLogin)
	//router.POST("/register", server.PostRegister)
	// authOnly
	router.GET("/house/:id", server.GetHouseById)
	router.POST("/house/:id/subscribe", server.PostHouseSubscribe)
	router.POST("/flat/create", server.PostFlatCreate)
	// moderationsOnly
	router.POST("/house/create", server.PostHouseCreate)
	router.POST("/flat/update", server.PostFlatUpdate)

	router.Run(os.Getenv("SERVER_PORT"))
}
