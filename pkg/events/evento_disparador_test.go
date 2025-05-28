package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload any
}

func (e *TestEvent) GetNome() string {
	return e.Name
}

func (e *TestEvent) GetValues() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventoInterface, wg *sync.WaitGroup) {
}

type EventoDisparadorTestSuite struct {
	suite.Suite
	event            TestEvent
	event2           TestEvent
	handler          TestEventHandler
	handler2         TestEventHandler
	handler3         TestEventHandler
	EventoDisparador *EventoDisparador
}

func (suite *EventoDisparadorTestSuite) SetupTest() {
	suite.EventoDisparador = NewEventoDisparador()
	suite.handler = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler3 = TestEventHandler{
		ID: 3,
	}
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
}

func (suite *EventoDisparadorTestSuite) TestEventoDisparador_Register() {
	err := suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	err = suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	assert.Equal(suite.T(), &suite.handler, suite.EventoDisparador.handlers[suite.event.GetNome()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.EventoDisparador.handlers[suite.event.GetNome()][1])
}

func (suite *EventoDisparadorTestSuite) TestEventoDisparador_Register_WithSameHandler() {
	err := suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	err = suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))
}

func (suite *EventoDisparadorTestSuite) TestEventoDisparador_Clear() {
	// Event 1
	err := suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	err = suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	// Event 2
	err = suite.EventoDisparador.Register(suite.event2.GetNome(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event2.GetNome()]))

	suite.EventoDisparador.Clear()
	suite.Equal(0, len(suite.EventoDisparador.handlers))
}

func (suite *EventoDisparadorTestSuite) TestEventoDisparador_Has() {
	// Event 1
	err := suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	err = suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	assert.True(suite.T(), suite.EventoDisparador.Has(suite.event.GetNome(), &suite.handler))
	assert.True(suite.T(), suite.EventoDisparador.Has(suite.event.GetNome(), &suite.handler2))
	assert.False(suite.T(), suite.EventoDisparador.Has(suite.event.GetNome(), &suite.handler3))
}

func (suite *EventoDisparadorTestSuite) TestEventoDisparador_Remove() {
	// Event 1
	err := suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	err = suite.EventoDisparador.Register(suite.event.GetNome(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	// Event 2
	err = suite.EventoDisparador.Register(suite.event2.GetNome(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event2.GetNome()]))

	suite.EventoDisparador.Remove(suite.event.GetNome(), &suite.handler)
	suite.Equal(1, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))
	assert.Equal(suite.T(), &suite.handler2, suite.EventoDisparador.handlers[suite.event.GetNome()][0])

	suite.EventoDisparador.Remove(suite.event.GetNome(), &suite.handler2)
	suite.Equal(0, len(suite.EventoDisparador.handlers[suite.event.GetNome()]))

	suite.EventoDisparador.Remove(suite.event2.GetNome(), &suite.handler3)
	suite.Equal(0, len(suite.EventoDisparador.handlers[suite.event2.GetNome()]))

}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventoInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventoDisparadorTestSuite) TestEventDispatch_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event)

	eh2 := &MockHandler{}
	eh2.On("Handle", &suite.event)

	//"Registrar" o manipulador de eventos
	suite.EventoDisparador.Register(suite.event.GetNome(), eh)
	suite.EventoDisparador.Register(suite.event.GetNome(), eh2)

	// Disparar o evento
	suite.EventoDisparador.Disparador(&suite.event)

	// Verificar se o manipulador foi chamado
	eh.AssertExpectations(suite.T())
	eh2.AssertExpectations(suite.T())

	// Verificar se o número de chamadas é o esperado
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eh2.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventoDisparadorTestSuite))
}
