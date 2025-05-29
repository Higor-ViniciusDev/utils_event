package main

import (
	"fmt"

	"github.com/Higor-ViniciusDev/utils/pkg/rabbitMQ"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch, err := rabbitMQ.OpenChannel()

	if err != nil {
		panic(err)
	}

	defer ch.Close()

	// Cria um canal para receber mensagens
	// Este canal será usado para receber mensagens do RabbitMQ
	out := make(chan amqp.Delivery)

	// Inicia o consumidor para escutar mensagens na fila
	// e envia as mensagens recebidas para o canal de saída
	go rabbitMQ.Consumer(ch, out)

	for msg := range out {
		// Processa a mensagem recebida
		// Aqui você pode adicionar a lógica para manipular a mensagem
		// Por exemplo, imprimir o conteúdo da mensagem
		fmt.Println("Mensagem recebida:", string(msg.Body))

		// Confirma o recebimento da mensagem
		msg.Ack(false)
	}
}
