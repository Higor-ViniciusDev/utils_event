package events

import (
	"sync"
	"time"
)

// EventoInterface representa um evento do sistema.
// Um evento é qualquer ação ou ocorrência que pode ser tratada por outros componentes.
// Exemplos: usuário fez login, pedido foi criado, arquivo foi enviado, etc.
type EventoInterface interface {
	// GetNome retorna o nome do evento, por exemplo: "UsuarioLogado".
	GetNome() string
	// GetDateTime retorna a data e hora em que o evento ocorreu.
	GetDateTime() time.Time
	// GetValues retorna os dados associados ao evento.
	// Pode ser qualquer informação relevante sobre o evento.
	GetValues() any
}

// EventoHandlerInterface representa um manipulador de eventos.
// Um handler é responsável por executar alguma ação quando um evento ocorre.
// Exemplo: enviar um e-mail quando um pedido é criado.
type EventoHandlerInterface interface {
	// Handle executa a ação desejada ao receber um evento.
	Handle(event EventoInterface, wg *sync.WaitGroup)
}

// EventDispachtInterface gerencia o registro e disparo de eventos e seus handlers.
// Ele permite adicionar, remover e acionar handlers para eventos específicos.
type EventDispachtInterface interface {
	// RegistrarHandler associa um handler a um evento específico pelo nome.
	// Assim, quando esse evento ocorrer, o handler será chamado.
	RegistrarHandler(eventoNome string, handler EventoHandlerInterface) error
	// Dispatch dispara um evento, chamando todos os handlers registrados para ele.
	Dispatch(evento EventoInterface) error
	// Remove remove um handler específico de um evento.
	Remove(eventoNome string, handler EventoHandlerInterface) error
	// HasHandlers verifica se existem handlers registrados para um evento.
	HasHandlers(eventoNome string, handle EventoHandlerInterface) bool
	// Clear remove todos os handlers de todos os eventos.
	Clear() error
}
