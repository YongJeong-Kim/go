package main

import (
	"context"
	"gorelationship/api"
	"gorelationship/config"
	"gorelationship/repository"
	"gorelationship/service"
)

func main() {
	ctx := context.Background()
	driver := config.NewDriver(ctx, config.URI, config.Username, config.Password)
	sess := config.NewSession(ctx, driver, config.DatabaseName)
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
