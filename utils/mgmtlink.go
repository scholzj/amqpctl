package utils

import (
	"qpid.apache.org/electron"
	"time"
	"qpid.apache.org/amqp"
)

type MgmtLink struct {
	Url string
	container electron.Container
	connection electron.Connection
	session electron.Session
	receiver electron.Receiver
	sender electron.Sender
	replyTo string
}

func (l *MgmtLink) Connect() (err error) {
	err = nil

	l.container = electron.NewContainer("amqpctl")

	l.connection, err = l.container.Dial("tcp", l.Url, electron.Heartbeat(time.Duration(10 * time.Second)))
	if err != nil {
		return
	}

	l.session, err = l.connection.Session(electron.IncomingCapacity(1024000), electron.OutgoingWindow(1024))
	if err != nil {
		return
	}

	l.replyTo = "tempAddress"

	l.receiver, err = l.session.Receiver(electron.Source(l.replyTo), electron.Capacity(100), electron.Prefetch(true))
	if err != nil {
		return
	}

	l.sender, err = l.session.Sender(electron.Target("$management"))
	if err != nil {
		return
	}

	return
}

func (l *MgmtLink) Close() {
	l.connection.Close(nil)
}

func (l *MgmtLink) Operation(reqProperties map[string]interface{}, reqBody map[interface{}]interface{}) (respProperties map[string]interface{}, respBody interface{}, err error) {
	err = nil
	respProperties = nil
	respBody = nil

	reqMsg := amqp.NewMessage()
	reqMsg.SetReplyTo(l.replyTo)

	if reqProperties != nil {
		reqMsg.SetProperties(reqProperties)
	}

	if reqBody != nil {
		reqMsg.Marshal(reqBody)
	}

	res := l.sender.SendSync(reqMsg)
	if res.Error != nil {
		err = res.Error
		return
	}

	respMsg, err := l.receiver.ReceiveTimeout(time.Duration(10 * time.Second))
	if err == nil {
		respProperties, respBody = l.parseResponse(respMsg.Message)
		respMsg.Accept()
		return
	} else {
		return
	}
}

func (l *MgmtLink) parseResponse(msg amqp.Message) (properties map[string]interface{}, body interface{}) {
	properties = msg.Properties()
	body = msg.Body()
	return
}