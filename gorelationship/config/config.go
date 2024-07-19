package config

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

const (
	Username     = "neo4j"
	Password     = "12341234"
	DatabaseName = "neo4j"
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	URI = "neo4j://localhost:17687"
)

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
