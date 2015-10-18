package consumer

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log"
	"testing"

	"github.com/ricbra/rabbitmq-cli-consumer/command"
	"github.com/ricbra/rabbitmq-cli-consumer/config"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseAndEscapesParamsInURI(t *testing.T) {
	uri := ParseURI("richard", "my@:secr%t", "localhost", "123", "/vhost")

	assert.Equal(t, "amqp://richard:my%40%3Asecr%25t@localhost:123/vhost", uri)
}

func TestAddsSlashWhenMissingInVhost(t *testing.T) {
	uri := ParseURI("richard", "secret", "localhost", "123", "vhost")

	assert.Equal(t, "amqp://richard:secret@localhost:123/vhost", uri)
}

func TestSetQosFails(t *testing.T) {
	config := createConfig()
	ch := new(TestChannel)

	var b bytes.Buffer
	errLogger := log.New(&b, "", 0)
	infLogger := log.New(&b, "", 0)

	ch.On("Qos", 3, 0, true).Return(errors.New("Error occured")).Once()

	err := Initialize(&config, ch, errLogger, infLogger)

	ch.AssertExpectations(t)
	ch.AssertNotCalled(t, "QueueDeclare", "worker", true, false, false, false, amqp.Table{})
	assert.NotNil(t, err)
}

func TestSetQosSucceeds(t *testing.T) {
	config := createConfig()
	ch := new(TestChannel)

	var b bytes.Buffer
	errLogger := log.New(&b, "", 0)
	infLogger := log.New(&b, "", 0)

	ch.On("Qos", 3, 0, true).Return(nil).Once()
	ch.On("QueueDeclare", "worker", true, false, false, false, amqp.Table{}).Return(amqp.Queue{}, errors.New("error")).Once()

	Initialize(&config, ch, errLogger, infLogger)

	ch.AssertExpectations(t)
}

func TestDeclareQueueFails(t *testing.T) {
	config := createConfig()
	ch := new(TestChannel)

	var b bytes.Buffer
	errLogger := log.New(&b, "", 0)
	infLogger := log.New(&b, "", 0)

	ch.On("Qos", 3, 0, true).Return(nil).Once()
	ch.On("QueueDeclare", "worker", true, false, false, false, amqp.Table{}).Return(amqp.Queue{}, errors.New("error")).Once()

	err := Initialize(&config, ch, errLogger, infLogger)

	ch.AssertExpectations(t)
	ch.AssertNotCalled(t, "ExchangeDeclare", "worker", "test", true, false, false, false, amqp.Table{})
	assert.NotNil(t, err)
}

func TestDeclareQueueSucceeds(t *testing.T) {
	config := createConfig()
	ch := new(TestChannel)

	var b bytes.Buffer
	errLogger := log.New(&b, "", 0)
	infLogger := log.New(&b, "", 0)

	ch.On("Qos", 3, 0, true).Return(nil).Once()
	ch.On("QueueDeclare", "worker", true, false, false, false, amqp.Table{}).Return(amqp.Queue{}, nil).Once()
	ch.On("ExchangeDeclare", "worker", "test", true, false, false, false, amqp.Table{}).Return(errors.New("error")).Once()

	Initialize(&config, ch, errLogger, infLogger)

	ch.AssertExpectations(t)
}

func TestBindQueueFails(t *testing.T) {
	config := createConfig()
	ch := new(TestChannel)

	var b bytes.Buffer
	errLogger := log.New(&b, "", 0)
	infLogger := log.New(&b, "", 0)

	ch.On("Qos", 3, 0, true).Return(nil).Once()
	ch.On("QueueDeclare", "worker", true, false, false, false, amqp.Table{}).Return(amqp.Queue{}, nil).Once()
	ch.On("ExchangeDeclare", "worker", "test", true, false, false, false, amqp.Table{}).Return(nil).Once()
	ch.On("QueueBind", "worker", "", "worker", false, amqp.Table{}).Return(errors.New("error")).Once()

	err := Initialize(&config, ch, errLogger, infLogger)

	ch.AssertExpectations(t)
	assert.NotNil(t, err)
}

func TestBindQueueSucceeds(t *testing.T) {
	config := createConfig()
	ch := new(TestChannel)

	var b bytes.Buffer
	errLogger := log.New(&b, "", 0)
	infLogger := log.New(&b, "", 0)

	ch.On("Qos", 3, 0, true).Return(nil).Once()
	ch.On("QueueDeclare", "worker", true, false, false, false, amqp.Table{}).Return(amqp.Queue{}, nil).Once()
	ch.On("ExchangeDeclare", "worker", "test", true, false, false, false, amqp.Table{}).Return(nil).Once()
	ch.On("QueueBind", "worker", "", "worker", false, amqp.Table{}).Return(nil).Once()

	err := Initialize(&config, ch, errLogger, infLogger)

	ch.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestProcessingMessageWithSuccess(t *testing.T) {
	msg := new(TestDelivery)
	executer := new(TestExecuter)
	factory := &command.CommandFactory{
		Cmd:  "test",
		Args: []string{"aa"},
	}
	consumer := Consumer{
		Executer:    executer,
		Factory:     factory,
		Compression: false,
	}
	body := []byte("the_body")
	args := base64.StdEncoding.EncodeToString(body)
	cmd := factory.Create(args)
	executer.On("Execute", cmd).Return(true).Once()
	msg.On("Body").Return(body).Once()
	msg.On("Ack", true).Return(nil).Once()

	consumer.ProcessMessage(msg)

	executer.AssertExpectations(t)
	msg.AssertExpectations(t)
}

func TestProcessingMessageWithFailure(t *testing.T) {
	msg := new(TestDelivery)
	executer := new(TestExecuter)
	factory := &command.CommandFactory{
		Cmd:  "test",
		Args: []string{"aa"},
	}
	consumer := Consumer{
		Executer:    executer,
		Factory:     factory,
		Compression: false,
	}
	body := []byte("the_body")
	args := base64.StdEncoding.EncodeToString(body)
	cmd := factory.Create(args)
	executer.On("Execute", cmd).Return(false).Once()
	msg.On("Body").Return(body).Once()
	msg.On("Nack", true, true).Return(nil).Once()

	consumer.ProcessMessage(msg)

	executer.AssertExpectations(t)
	msg.AssertExpectations(t)
}

type TestCommand struct {
	mock.Mock
}

func (t *TestCommand) CombinedOutput() (out []byte, err error) {
	argsT := t.Called()

	return argsT.Get(0).([]byte), argsT.Error(1)
}

type TestExecuter struct {
	mock.Mock
}

func (t *TestExecuter) Execute(cmd command.Command) bool {
	argsT := t.Called(cmd)

	return argsT.Get(0).(bool)
}

type TestDelivery struct {
	mock.Mock
	body []byte
}

func (t *TestDelivery) Ack(multiple bool) error {
	argstT := t.Called(multiple)

	return argstT.Error(0)
}

func (t *TestDelivery) Nack(multiple, requeue bool) error {
	argsT := t.Called(multiple, requeue)

	return argsT.Error(0)
}

func (t *TestDelivery) Body() []byte {
	argsT := t.Called()

	return argsT.Get(0).([]byte)
}

type TestChannel struct {
	mock.Mock
}

func (t *TestChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	argsT := t.Called(name, kind, durable, autoDelete, internal, noWait, args)

	return argsT.Error(0)
}

func (t *TestChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	argsT := t.Called(name, durable, autoDelete, exclusive, noWait, args)

	return argsT.Get(0).(amqp.Queue), argsT.Error(1)
}

func (t *TestChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	argsT := t.Called(prefetchCount, prefetchSize, global)

	return argsT.Error(0)
}

func (t *TestChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	argsT := t.Called(name, key, exchange, noWait, args)

	return argsT.Error(0)
}

func createConfig() config.Config {
	return config.CreateFromString(`[rabbitmq]
  host=localhost
  username=ricbra
  password=t3st
  vhost=staging
  queue=worker
  port=123

  [prefetch]
  count=3
  global=On

  [exchange]
  name=worker
  autodelete=Off
  type=test
  durable=On

  [logs]
  error=a
  info=b
  `)
}