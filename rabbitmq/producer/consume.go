package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	exchangeName := "logs_topic"

	// Объявляем обменник (должен совпадать с producer)
	ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)

	// Создаем временную очередь с уникальным именем
	q, _ := ch.QueueDeclare(
		"",    // пустое имя - сервер сгенерирует уникальное
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	// Привязываем очередь к обменнику с нужными routing keys
	routingKeys := []string{"user.*", "email.#"}
	for _, key := range routingKeys {
		ch.QueueBind(
			q.Name,
			key,
			exchangeName,
			false,
			nil,
		)
		log.Printf("Привязано: %s -> %s", key, q.Name)
	}

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	log.Println("Ожидание сообщений...")
	for d := range msgs {
		log.Printf("Получено [%s]: %s", d.RoutingKey, d.Body)
	}
}
