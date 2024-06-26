package embedded

import (
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"log"
	"net/url"
)

var Servers = []string{
	"nats://localhost:14222",
	"nats://localhost:24222",
	"nats://localhost:34222",
}
var Opts1 = &server.Options{
	Host:     "0.0.0.0",
	Port:     14222,
	HTTPPort: 18222,
	Routes: []*url.URL{
		{
			Scheme: "nats",
			Host:   "localhost:16222",
		},
	},
	Cluster: server.ClusterOpts{
		Name: "NATS",
		Host: "localhost",
		Port: 16222,
	},
	Username: "aaa",
	Password: "1234",
	//ConfigFile: "",
}

var Opts2 = &server.Options{
	Host: "0.0.0.0",
	Port: 24222,
	//RoutesStr: "nats://localhost:16222",
	HTTPPort: 28222,
	Routes: []*url.URL{
		{
			Scheme: "nats",
			Host:   "localhost:16222",
		},
	},
	Cluster: server.ClusterOpts{
		Name: "NATS",
		Host: "localhost",
		Port: 26222,
	},
	Username: "aaa",
	Password: "1234",

	//ConfigFile: "",
}

var Opts3 = &server.Options{
	Host: "0.0.0.0",
	Port: 34222,
	//RoutesStr: "nats://localhost:16222",
	HTTPPort: 38222,
	Routes: []*url.URL{
		{
			Scheme: "nats",
			Host:   "localhost:16222",
		},
	},
	Cluster: server.ClusterOpts{
		Name: "NATS",
		Host: "localhost",
		Port: 36222,
	},
	Username: "aaa",
	Password: "1234",
	//ConfigFile: "",
}

func NewServer(opts *server.Options) *server.Server {
	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return ns
}

func (n *Notify) Publish(subject string, data []byte) error {
	return n.Nats.Publish(subject, data)
}

func (n *Notify) Drain() error {
	return n.Drain()
}

type Notifier interface {
	Publish(subject string, data []byte) error
	Drain() error
}

type Notify struct {
	Nats *nats.Conn
}

func NewNotify(conn *nats.Conn) *Notify {
	return &Notify{
		Nats: conn,
	}
}
