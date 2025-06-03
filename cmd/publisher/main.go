package main

import "github.com/Higor-ViniciusDev/utils/pkg/rabbitmq"

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbitmq.Publish(ch, "nova queue fila para tratamento", "amq.direct")
}
