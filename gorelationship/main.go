package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"time"
)

const username = "neo4j"
const password = "12341234"
const databaseName = "neo4j"

// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
const uri = "neo4j://localhost:17687"

type Server struct {
	Session neo4j.SessionWithContext
}

func main() {
	ctx := context.Background()
	driver := NewDriver(ctx, uri, username, password)
	s := NewServer(NewSession(ctx, driver, databaseName))
	defer s.Session.Close(ctx)

	/*uID, err := s.CreateUser(ctx, map[string]any{
		"name": "dddname",
	})
	if err != nil {
		panic(err)
	}
	log.Println(uID)*/

	s.FriendRequest(ctx, map[string]any{
		"from": "dddname",
		"to":   "aaaname",
	})

	/*s.FriendRequest(ctx, map[string]any{
		"from": "aaaname",
		"to":   "bbbname",
	})*/

	s.ListFriendRequests(ctx, map[string]any{
		"name": "aaaname",
	})
}

func NewServer(s neo4j.SessionWithContext) *Server {
	return &Server{
		Session: s,
	}
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

func (s *Server) CreateUser(ctx context.Context, param map[string]any) (string, error) {
	uID, err := s.Session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			CREATE (u:User {id: randomuuid(), name: $name, createdDate: datetime()})
			RETURN u.id AS id
		`, param)
		if err != nil {
			return "", err
		}

		u, err := result.Single(ctx)
		if err != nil {
			return "", err
		}
		uID, _ := u.AsMap()["id"]
		return uID.(string), nil
	})
	if err != nil {
		return "", err
	}

	return uID.(string), nil
}

func (s *Server) FriendRequest(ctx context.Context, param map[string]any) error {
	_, err := s.Session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (from:User {name: $from})
			MATCH (to:User {name: $to})
			MERGE (from)-[:FRIEND {status: ['request']}]->(to)
		`, param)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) AcceptFriendRequest(ctx context.Context, param map[string]any) error {
	_, err := s.Session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (from:User {name: $from})
			MATCH (to:User {name: $to})
			MERGE (to)-[:FRIEND {status: ['accept']}]->(from)
		`, param)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

type ListFriendRequestsResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (s *Server) ListFriendRequests(ctx context.Context, param map[string]any) ([]ListFriendRequestsResult, error) {
	requests, err := s.Session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name})-[:FRIEND]->(requests) 
			RETURN requests.id AS id, requests.name AS name, requests.createdDate AS createdDate
		`, param)
		if err != nil {
			return nil, err
		}

		requests, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		var rz []ListFriendRequestsResult
		for _, r := range requests {
			var rr ListFriendRequestsResult
			rr.ID = r.AsMap()["id"].(string)
			rr.Name = r.AsMap()["name"].(string)
			rr.CreatedDate = r.AsMap()["createdDate"].(time.Time)
			rz = append(rz, rr)
		}

		return rz, nil
	})
	if err != nil {
		//return nil, err
		log.Println(err)
	}

	return requests.([]ListFriendRequestsResult), nil
}

func (s *Server) DeleteAll(ctx context.Context) error {
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
}
