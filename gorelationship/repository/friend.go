package repository

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type Friender interface {
	Accept(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, acceptUserID string) (int, error)
	Count(ctx context.Context, tx neo4j.ManagedTransaction, userID string) (int64, error)
	List(ctx context.Context, tx neo4j.ManagedTransaction, userID string) ([]ListResult, error)
	ListMutuals(ctx context.Context, tx neo4j.ManagedTransaction, userID1, userID2 string) ([]ListMutualsResult, error)
	ListFromRequests(ctx context.Context, tx neo4j.ManagedTransaction, userID string) ([]ListFromRequestsResult, error)
	MutualCount(ctx context.Context, tx neo4j.ManagedTransaction, userID1, userID2 string) (int64, error)
	Request(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (int, error)
	FromRequestCount(ctx context.Context, tx neo4j.ManagedTransaction, userID string) (int64, error)
	Validate(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) error
	RelationshipStatus(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (string, error)
}

type ListFromRequestsResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) ListFromRequests(ctx context.Context, tx neo4j.ManagedTransaction, userID string) ([]ListFromRequestsResult, error) {
	result, err := tx.Run(ctx, `
		MATCH (u:User)-[:REQUEST]->(:User {id: $userID}) 
		RETURN u.id AS id, u.name AS name, u.createdDate AS createdDate
	`, map[string]any{
		"userID": userID,
	})
	if err != nil {
		return nil, err
	}

	var rz []ListFromRequestsResult
	for result.Next(ctx) {
		var rr ListFromRequestsResult
		rr.ID = result.Record().AsMap()["id"].(string)
		rr.Name = result.Record().AsMap()["name"].(string)
		rr.CreatedDate = result.Record().AsMap()["createdDate"].(time.Time)
		rz = append(rz, rr)
	}
	if result.Err() != nil {
		return nil, err
	}

	return rz, nil
}

func (f *Friend) Request(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (int, error) {
	result, err := tx.Run(ctx, `
		MATCH (r:User {id: $request})
		MATCH (a:User {id: $approve})
		MERGE (r)-[:REQUEST]->(a)
	`, map[string]any{
		"request": requestUserID,
		"approve": approveUserID,
	})
	if err != nil {
		return 0, err
	}

	c, err := result.Consume(ctx)
	if err != nil {
		return 0, err
	}

	return c.Counters().RelationshipsCreated(), nil
}

func (f *Friend) Accept(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, acceptUserID string) (int, error) {
	result, err := tx.Run(ctx, `
		MATCH (r:User {id: $request})-[req:REQUEST]->(a:User {id: $accept})
		DELETE req
		CREATE (a)-[:FRIEND]->(r)
	`, map[string]any{
		"request": requestUserID,
		"accept":  acceptUserID,
	})
	if err != nil {
		return 0, err
	}

	c, err := result.Consume(ctx)
	if err != nil {
		return 0, err
	}

	return c.Counters().RelationshipsCreated(), nil
}

type ListResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) List(ctx context.Context, tx neo4j.ManagedTransaction, userID string) ([]ListResult, error) {
	result, err := tx.Run(ctx, `
		MATCH (:User {id: $userID})-[:FRIEND]->(u:User)
		RETURN u.id AS id, u.name AS name, u.createdDate AS createdDate
	`, map[string]any{
		"userID": userID,
	})
	if err != nil {
		return nil, err
	}

	frs, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}

	var rs []ListResult
	for _, fr := range frs {
		var rr ListResult
		rr.ID = fr.AsMap()["id"].(string)
		rr.Name = fr.AsMap()["name"].(string)
		rr.CreatedDate = fr.AsMap()["createdDate"].(time.Time)
		rs = append(rs, rr)
	}

	return rs, nil
}

type ListMutualsResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) ListMutuals(ctx context.Context, tx neo4j.ManagedTransaction, userID1, userID2 string) ([]ListMutualsResult, error) {
	result, err := tx.Run(ctx, `
		MATCH (:User {id: $userID})-[:FRIEND]-(fs)-[:FRIEND]-(:User {id: $userID2})
		RETURN fs.id AS id, fs.name AS name, fs.createdDate AS createdDate
	`, map[string]any{
		"userID1": userID1,
		"userID2": userID2,
	})
	if err != nil {
		return nil, err
	}

	mf, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}

	var mz []ListMutualsResult
	for _, m := range mf {
		var f ListMutualsResult
		f.ID = m.AsMap()["id"].(string)
		f.Name = m.AsMap()["name"].(string)
		f.CreatedDate = m.AsMap()["createdDate"].(time.Time)
		mz = append(mz, f)
	}

	return mz, nil
}

func (f *Friend) MutualCount(ctx context.Context, tx neo4j.ManagedTransaction, userID1, userID2 string) (int64, error) {
	result, err := tx.Run(ctx, `
		MATCH (:User {id: $userID1})-[:FRIEND]-(fs)-[:FRIEND]-(:User {id: $userID2})
		RETURN COUNT(fs) AS count
	`, map[string]any{
		"userID1": userID1,
		"userID2": userID2,
	})
	if err != nil {
		return 0, err
	}

	cnt, err := result.Single(ctx)
	if err != nil {
		return 0, err
	}

	return cnt.AsMap()["count"].(int64), nil
}

func (f *Friend) Count(ctx context.Context, tx neo4j.ManagedTransaction, userID string) (int64, error) {
	result, err := tx.Run(ctx, `
		MATCH (:User {id: $userID})-[f:FRIEND]->(:User)
		RETURN COUNT(f) AS count
	`, map[string]any{
		"userID": userID,
	})
	if err != nil {
		return 0, err
	}

	cnt, err := result.Single(ctx)
	if err != nil {
		return 0, err
	}
	return cnt.AsMap()["count"].(int64), nil
}

func (f *Friend) FromRequestCount(ctx context.Context, tx neo4j.ManagedTransaction, userID string) (int64, error) {
	result, err := tx.Run(ctx, `
		MATCH (:User)-[r:REQUEST]->(:User {id: $userID})
		RETURN COUNT(r) AS count
	`, map[string]any{
		"userID": userID,
	})
	if err != nil {
		return 0, err
	}

	cnt, err := result.Single(ctx)
	if err != nil {
		return 0, err
	}

	return cnt.AsMap()["count"].(int64), nil
}

func (f *Friend) Validate(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) error {
	rc, err := tx.Run(ctx, `
		MATCH (r:User) WHERE r.id = $request
		MATCH (a:User) WHERE a.id = $approve
		RETURN r.id AS rID, a.id AS aID
	`, map[string]any{
		"request": requestUserID,
		"approve": approveUserID,
	})
	if err != nil {
		return err
	}

	r, err := rc.Single(ctx)
	if err != nil {
		return err
	}

	if r.AsMap()["rID"] == nil {
		return errors.New("request user not found")
	}

	if r.AsMap()["aID"] == nil {
		return errors.New("approve user not found")
	}
	return nil
}

func (f *Friend) RelationshipStatus(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, acceptUserID string) (string, error) {
	result, err := tx.Run(ctx, `
		OPTIONAL MATCH (:User {id: $request})-[r:REQUEST]->(:User {id: $accept})
		OPTIONAL MATCH (:User {id: $accept})-[f:FRIEND]->(:User {id: $request})
		RETURN 
			CASE r 
				WHEN IS NULL THEN 
					CASE f 
						WHEN IS NULL THEN 'NO'
						ELSE type(f)
					END
				ELSE type(r) 
			END AS type
	`, map[string]any{
		"request": requestUserID,
		"accept":  acceptUserID,
	})
	if err != nil {
		return "", err
	}

	rs, err := result.Single(ctx)
	if err != nil {
		return "", err
	}

	return rs.AsMap()["type"].(string), nil
}

type Friend struct {
}

func NewFriend() *Friend {
	return &Friend{}
}
