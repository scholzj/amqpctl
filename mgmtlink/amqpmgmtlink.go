package mgmtlink

import (
	"qpid.apache.org/electron"
	"time"
	"qpid.apache.org/amqp"
	"crypto/x509"
	"crypto/tls"
	"net"
	"fmt"
	"io/ioutil"
	"os"
)

type AmqpConfiguration struct {
	AmqpHostname string
	AmqpPort int32
	AmqpUsername string
	AmqpPassword string
	SaslMechanism string
	SslCaFile string
	SslCertFile string
	SslKeyFile string
	SslSkipHostnameVerification bool
}

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

func (l *AmqpMgmtLink) ConfigureConnection(amqpCfg AmqpConfiguration) (err error) {
	err = nil

	// URL (Hostname / port)
	l.Url = fmt.Sprintf("%v:%v", amqpCfg.AmqpHostname, amqpCfg.AmqpPort)

	// Username
	if amqpCfg.AmqpUsername != "" {
		l.Username = amqpCfg.AmqpUsername
	}

	// Password
	if amqpCfg.AmqpPassword != "" {
		l.Password = amqpCfg.AmqpPassword
	}

	// SASL Mechanism
	if amqpCfg.SaslMechanism != "" {
		l.SaslMechanism = amqpCfg.SaslMechanism
	}

	// CA certificate
	if amqpCfg.SslCaFile != "" {
		brokerCert, err := ioutil.ReadFile(amqpCfg.SslCaFile)
		if err != nil {
			fmt.Printf("Ups, something went wrong while loading CA certificate %s ... %s", amqpCfg.SslCaFile, err)
			os.Exit(1)
		}

		brokerCertPool := x509.NewCertPool()
		brokerCertPool.AppendCertsFromPEM(brokerCert)

		l.BrokerCertificate = brokerCertPool

		// Hostname verification
		l.SslSkipVerify = amqpCfg.SslSkipHostnameVerification

		// Client certificate and key
		if amqpCfg.SslCertFile != "" && amqpCfg.SslKeyFile != "" {
			memberKey, err := tls.LoadX509KeyPair(amqpCfg.SslCertFile, amqpCfg.SslKeyFile)
			if err != nil {
				fmt.Printf("Ups, something went wrong while loading client certificate (%s) / key (%s) ... %s", amqpCfg.SslCertFile, amqpCfg.SslKeyFile, err)
				os.Exit(1)
			}

			l.ClientCertificate = &memberKey
		}
	}

	return
}