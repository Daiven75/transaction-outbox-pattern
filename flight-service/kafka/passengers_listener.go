package kafka

import (
	"encoding/json"
	"errors"
	"flight-service/model"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"log"
)

type Listener struct {
	Consumer  sarama.Consumer
	Topic     string
	Partition int32
}

func NewKafkaListener(broker, topic string, partition int32) (*Listener, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		return nil, errors.New("error creating Kafka consumer" + err.Error())
	}

	return &Listener{
		consumer,
		topic,
		partition,
	}, nil
}

func (l *Listener) Listen(db *gorm.DB) {
	partitionConsumer, err := l.Consumer.ConsumePartition(l.Topic, l.Partition, sarama.OffsetNewest)

	if err != nil {
		log.Println("Failed to start partition consumer", err)
		panic(err)
	}

	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		log.Printf("Message received: %s\n", string(msg.Value))

		var flight model.Flight
		var passenger model.Passenger

		err := json.Unmarshal(msg.Value, &passenger)
		if err != nil {
			log.Printf("Error converting JSON to struct: %v", err)
			return
		}

		result := db.Preload("Passengers").First(&flight, passenger.FlightID)
		if result.Error != nil {
			log.Printf("Flight with id = %d not found!", passenger.FlightID)
			return
		}

		flight.Passengers = append(flight.Passengers, passenger)
		db.Create(&passenger)
	}
}

func (l *Listener) Close() error {
	return l.Consumer.Close()
}
