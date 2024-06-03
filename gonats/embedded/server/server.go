package main

import (
	"context"
	"github.com/nats-io/nats-server/v2/server"
	"golang.org/x/sync/errgroup"
	_ "golang.org/x/sync/errgroup"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
	os.Interrupt,
}

var opts1 = &server.Options{
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

var opts2 = &server.Options{
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

var opts3 = &server.Options{
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

func main() {
	s1 := NewServer(opts1)
	s2 := NewServer(opts2)
	s3 := NewServer(opts3)

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		log.Println("start server 1")
		s1.Start()
		log.Println("start server 2")
		s2.Start()
		log.Println("start server 3")
		s3.Start()
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		log.Println("shutdown server 1")
		s1.WaitForShutdown()
		log.Println("shutdown server 2")
		s2.WaitForShutdown()
		log.Println("shutdown server 3")
		s3.WaitForShutdown()
		//ns.Shutdown()
		return nil
	})

	/*ticker := time.NewTicker(60 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				log.Println("in done")
				return
			case t := <-ticker.C:
				log.Println("ticker in.", t)
				ticker.Stop()
				//done <- true
				log.Println("shutdown server 1")
				s1.Shutdown()
				return
			}
		}
	}()*/

	err := group.Wait()
	if err != nil {
		log.Fatal("group wait error.", err)
	}
}

func NewServer(opts *server.Options) *server.Server {
	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return ns
}
