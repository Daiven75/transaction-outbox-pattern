package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"passenger-service/model"
)

func DatabaseConnection(host, port, username, password, database string) *gorm.DB {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dbInfo), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return db
}

func main() {
	viper.SetDefault("DB.HOST", "localhost")
	viper.SetDefault("DB.PORT", "3306")
	viper.SetDefault("DB.USERNAME", "root")
	viper.SetDefault("DB.PASSWORD", "mysql")
	viper.SetDefault("DB.DATABASE", "passenger_db")

	env := os.Getenv("ENV")

	if env == "docker" {
		viper.BindEnv("DB.HOST")
		viper.BindEnv("DB.PORT")
		viper.BindEnv("DB.USERNAME")
		viper.BindEnv("DB.PASSWORD")
		viper.BindEnv("DB.DATABASE")
	}

	host := viper.GetString("DB.HOST")
	port := viper.GetString("DB.PORT")
	user := viper.GetString("DB.USERNAME")
	password := viper.GetString("DB.PASSWORD")
	database := viper.GetString("DB.DATABASE")

	fmt.Printf("Host: %s, Port: %s, Username: %s, Password: %s, Database: %s\n", host, port, user, password, database)

	db := DatabaseConnection(host, port, user, password, database)

	db.Migrator().DropTable(&model.Passenger{}, &model.PassengerOutbox{})
	db.AutoMigrate(&model.Passenger{}, &model.PassengerOutbox{})

	router := gin.Default()

	router.GET("/passengers", func(ctx *gin.Context) {
		var passenger []model.Passenger

		result := db.Find(&passenger)
		if result.Error != nil {
			panic(result.Error)
		}

		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusOK, &passenger)
	})

	router.POST("/passengers", func(ctx *gin.Context) {
		var passenger model.Passenger

		if err := ctx.ShouldBind(&passenger); err != nil {
			panic(err)
		}

		err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Create(&passenger).Error; err != nil {
				return err
			}

			var passengerOutbox model.PassengerOutbox
			passengerOutbox.FromPassenger(passenger)

			if err := tx.Create(&passengerOutbox).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Fatal("Transaction error:", err)
		} else {
			fmt.Println("Transaction completed successfully!")
		}

		ctx.Status(http.StatusCreated)
	})

	server := http.Server{
		Addr:    ":8887",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return
	}

}
