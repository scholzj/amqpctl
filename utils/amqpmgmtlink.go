package utils

import (
	"qpid.apache.org/electron"
	"time"
	"qpid.apache.org/amqp"
	"crypto/x509"
	"crypto/tls"
	"net"
)

type AmqpMgmtLink struct {
	Url string
	Username string
	Password string
	SaslMechanism string
	BrokerCertificate *x509.CertPool
	SslSkipVerify bool
	ClientCertificate *tls.Certificate
	container electron.Container
	connection electron.Connection
	session electron.Session
	receiver electron.Receiver
	sender electron.Sender
	replyTo string
}

func (l *AmqpMgmtLink) Connect() (err error) {
	err = nil

	l.container = electron.NewContainer("amqpctl")

	err = l.connectAmqp()
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

func (l *AmqpMgmtLink) connectAmqp() (err error) {
	err = nil
	var conn net.Conn

	if l.BrokerCertificate != nil {
		conn, err = l.connectSsl()
		if err != nil {
			return
		}
	} else {
		conn, err = l.connectTcp()
		if err != nil {
			return
		}
	}

	var options []electron.ConnectionOption
	options = append(options, electron.Heartbeat(time.Duration(10 * time.Second)))

	if l.Username != "" && l.Password != "" {
		options = append(options, electron.User(l.Username))
		options = append(options, electron.Password([]byte(l.Password)))
	}

	if l.SaslMechanism != "" {
		options = append(options, electron.SASLAllowedMechs(l.SaslMechanism))
	}

	l.connection, err = l.container.Connection(conn, options...)

	if err != nil {
		return
	}

	return
}

func (l *AmqpMgmtLink) connectSsl() (conn net.Conn, err error) {
	err = nil

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = l.SslSkipVerify

	if l.BrokerCertificate != nil {
	tlsConfig.ClientCAs = l.BrokerCertificate
	}

	if l.ClientCertificate != nil {
	tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return l.ClientCertificate, nil
	}
	}

	conn, err = tls.Dial("tcp", l.Url, tlsConfig)

	return
}

func (l *AmqpMgmtLink) connectTcp() (conn net.Conn, err error) {
	err = nil

	conn, err = net.Dial("tcp", l.Url)

	return
}

func (l *AmqpMgmtLink) Close() {
	l.connection.Close(nil)
}

func (l *AmqpMgmtLink) Operation(reqProperties map[string]interface{}, reqBody map[interface{}]interface{}) (respProperties map[string]interface{}, respBody interface{}, err error) {
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

func (l *AmqpMgmtLink) parseResponse(msg amqp.Message) (properties map[string]interface{}, body interface{}) {
	properties = msg.Properties()
	body = msg.Body()
	return
}