package main

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Устанавливаем соединение с RabbitMQ
	// URL формата: amqp://[username][:password]@host[:port][/vhost]
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("Соединение с RabbitMQ установлено")

	// 2. Создаем канал (channel)
	// Большинство операций выполняются через каналы
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %v", err)
	}
	defer ch.Close()

	// 3. Объявляем очередь, из которой будем читать
	// Если очередь не существует, она будет создана. Это идемпотентная операция.
	queueName := "hello"
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable (сохранять ли очередь при перезапуске RabbitMQ)
		false,     // delete when unused (автоудаление)
		false,     // exclusive (использоваться только этим соединением)
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %v", err)
	}
	log.Printf("Очередь '%s' объявлена", q.Name)

	// 4. Регистрируем потребителя (consumer)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer (уникальный тег, пустая строка = сгенерировать)
		false,  // auto-ack (автоматически подтверждать получение)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать потребителя: %v", err)
	}

	// 5. Запускаем горутину для обработки входящих сообщений
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Получено сообщение: %s", d.Body)

			// Имитация обработки сообщения
			time.Sleep(1 * time.Second)

			// Ручное подтверждение обработки сообщения
			// Благодаря auto-ack: false, мы контролируем, когда сообщение
			// будет удалено из очереди.
			d.Ack(false)
		}
	}()

	log.Println(" [*] Ожидание сообщений. Для выхода нажмите CTRL+C")
	<-forever
}
