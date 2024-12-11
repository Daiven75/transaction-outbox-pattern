package main

import (
	"flight-service/kafka"
	"flight-service/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func DatabaseConnection(host, port, user, password, dbName string) *gorm.DB {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})

	if err != nil {
		fmt.Println("failed to connect database", err)
		log.Fatal(err)
	}

	return db
}

func main() {
	viper.SetDefault("DB.HOST", "localhost")
	viper.SetDefault("DB.PORT", "5432")
	viper.SetDefault("DB.USERNAME", "postgres")
	viper.SetDefault("DB.PASSWORD", "postgres")
	viper.SetDefault("DB.DATABASE", "flight_db")
	viper.SetDefault("KAFKA.BROKER", "localhost:9092")

	env := os.Getenv("ENV")

	if env == "docker" {
		viper.BindEnv("DB.HOST")
		viper.BindEnv("DB.PORT")
		viper.BindEnv("DB.USERNAME")
		viper.BindEnv("DB.PASSWORD")
		viper.BindEnv("DB.DATABASE")
		viper.BindEnv("KAFKA.BROKER")
	}

	host := viper.GetString("DB.HOST")
	port := viper.GetString("DB.PORT")
	user := viper.GetString("DB.USERNAME")
	password := viper.GetString("DB.PASSWORD")
	dbName := viper.GetString("DB.DATABASE")
	broker := viper.GetString("KAFKA.BROKER")

	fmt.Printf("Host: %s, Port: %s, Username: %s, Password: %s, Database: %s, kafka: %s\n", host, port, user, password, dbName, broker)

	db := DatabaseConnection(host, port, user, password, dbName)

	db.Migrator().DropTable(&model.Flight{}, &model.Passenger{})
	db.AutoMigrate(&model.Flight{}, &model.Passenger{})

	listener, err := kafka.NewKafkaListener(broker, "mysql.passenger_db.passenger_outbox", 0)
	if err != nil {
		log.Fatalf("Error creating Kafka: %v", err)
	}

	defer listener.Close()

	go func() {
		fmt.Println("Starting Kafka Consumer...")
		listener.Listen(db)
	}()

	router := gin.Default()

	router.GET("/flights", func(ctx *gin.Context) {
		var flights []model.Flight

		result := db.Preload("Passengers").Find(&flights)
		if result.Error != nil {
			panic(result.Error)
		}

		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusOK, flights)
	})

	router.POST("/flights", func(ctx *gin.Context) {
		var flight model.Flight
		if err := ctx.ShouldBind(&flight); err != nil {
			panic(err)
		}

		db.Create(&flight)

		ctx.Status(http.StatusCreated)
	})

	server := http.Server{
		Addr:    ":8888",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		return
	}
}
