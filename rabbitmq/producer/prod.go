package producer

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	// Объявляем обменник
	exchangeName := "logs_topic"
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Отправляем сообщения с разными routing keys
	messages := []struct {
		routingKey string
		body       string
	}{
		{"user.created", "Новый пользователь: John"},
		{"user.updated", "Обновлен пользователь: John"},
		{"email.sent", "Письмо отправлено: john@example.com"},
	}

	for _, msg := range messages {
		err = ch.PublishWithContext(ctx,
			exchangeName,
			msg.routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.body),
			})

		if err != nil {
			log.Printf("Ошибка: %v", err)
		} else {
			log.Printf("Отправлено [%s]: %s", msg.routingKey, msg.body)
		}
	}
}
