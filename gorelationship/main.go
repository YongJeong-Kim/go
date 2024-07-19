package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorelationship/api"
	"gorelationship/repository"
	"gorelationship/service"
	"log"
)

const (
	username     = "neo4j"
	password     = "12341234"
	databaseName = "neo4j"
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	uri = "neo4j://localhost:17687"
)

func main() {
	ctx := context.Background()
	driver := NewDriver(ctx, uri, username, password)
	sess := NewSession(ctx, driver, databaseName)
	defer sess.Close(ctx)

	var rf repository.Friender = repository.NewFriend(sess)
	var ru repository.UserManager = repository.NewUser(sess)
	//repo := repository.NewRepository(rf, ru)

	var sf service.Friender = service.NewFriend(rf)
	var su service.UserManager = service.NewUser(ru)
	svc := service.NewService(sf, su)

	server := api.NewServer(svc)
	server.SetupRouter()
	server.Router.Run(":8080")
}

func NewDriver(ctx context.Context, uri, username, password string) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	//defer driver.Close(ctx)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Connectivity successful")

	return driver
}

func NewSession(ctx context.Context, driver neo4j.DriverWithContext, databaseName string) neo4j.SessionWithContext {
	return driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	//defer session.Close(ctx)
}

/*func (s *Server) DeleteAll(ctx context.Context) error {
	_, err := s.Session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (n) DETACH DELETE n
		`, nil)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}*/
