package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventoDisparador struct {
	handlers map[string][]EventoHandlerInterface
}

func NewEventoDisparador() *EventoDisparador {
	return &EventoDisparador{
		handlers: make(map[string][]EventoHandlerInterface),
	}
}

func (ev *EventoDisparador) Disparador(event EventoInterface) error {
	if handlers, ok := ev.handlers[event.GetNome()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

// Register adiciona um novo handler para um evento específico no EventoDisparador.
// Se o handler já estiver registrado para o evento, retorna ErrHandlerAlreadyRegistered.
//
// Por exemplo, para registrar dois handlers diferentes para o evento "venda_finalizada",
// como disparar um e-mail e disparar uma cobrança:
//
//	disparador := &EventoDisparador{handlers: make(map[string][]EventoHandlerInterface)}
//	disparador.Register("venda_finalizada", EmailHandler{})
//	disparador.Register("venda_finalizada", CobrancaHandler{})
func (ed *EventoDisparador) Register(eventName string, handler EventoHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}

	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventoDisparador) Has(eventName string, handler EventoHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventoDisparador) Remove(eventName string, handler EventoHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (ed *EventoDisparador) Clear() {
	ed.handlers = make(map[string][]EventoHandlerInterface)
}
