package rabbitMQ

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp.Channel, error) {
	//Conecta ao RabbitMQ server na porta 5672
	// com o usuário guest e senha guest
	// Certifique-se de que o RabbitMQ está rodando e acessível
	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}

	return channel, nil
}

func Consumer(ch *amqp.Channel, out chan amqp.Delivery) error {
	// Inicia a escuta na fila e retorna um canal de mensagens recebidas
	msgs, err := ch.Consume(
		"minhaFila", // Nome da fila
		"go-chan",   // Consumer tag (pode ser vazio)
		false,       // Auto-acknowledge - Dar baixa na mensagem automaticamente e já pode remover da fila
		false,       // Exclusivo - Não exclusivo para este consumidor
		false,       // Não compartilhar com outros consumidores
		false,       // Não esperar confirmação de entrega
		nil,         // Argumentos adicionais - pode ser nil
	)

	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg     // Envia a mensagem recebida para o canal de saída
		msg.Ack(false) // Confirma o recebimento da mensagem
	}

	return nil
}
