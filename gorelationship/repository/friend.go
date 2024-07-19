package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type Friender interface {
	Accept(ctx context.Context, param map[string]any) error
	Count(ctx context.Context, param map[string]any) (int64, error)
	List(ctx context.Context, param map[string]any) ([]ListResult, error)
	ListMutuals(ctx context.Context, param map[string]any) ([]ListMutualsResult, error)
	ListRequests(ctx context.Context, param map[string]any) ([]ListRequestsResult, error)
	MutualCount(ctx context.Context, param map[string]any) (int64, error)
	Request(ctx context.Context, param map[string]any) error
	RequestCount(ctx context.Context, param map[string]any) (int64, error)
}

type ListRequestsResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) ListRequests(ctx context.Context, param map[string]any) ([]ListRequestsResult, error) {
	requests, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name})-[:FRIEND]->(requests) 
			RETURN requests.id AS id, requests.name AS name, requests.createdDate AS createdDate
		`, param)
		if err != nil {
			return nil, err
		}

		var rz []ListRequestsResult
		for result.Next(ctx) {
			var rr ListRequestsResult
			rr.ID = result.Record().AsMap()["id"].(string)
			rr.Name = result.Record().AsMap()["name"].(string)
			rr.CreatedDate = result.Record().AsMap()["createdDate"].(time.Time)
			rz = append(rz, rr)
		}

		return rz, nil
	})
	if err != nil {
		return nil, err
	}

	return requests.([]ListRequestsResult), nil
}

func (f *Friend) Request(ctx context.Context, param map[string]any) error {
	_, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (f *Friend) Accept(ctx context.Context, param map[string]any) error {
	_, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (request:User {name: $request})
			MATCH (approve:User {name: $approve})
			MERGE (approve)-[:FRIEND {status: ['accept']}]->(request)
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

type ListResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) List(ctx context.Context, param map[string]any) ([]ListResult, error) {
	fs, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name})-[:FRIEND {status: ['accept']}]->(fs)
			RETURN fs.id AS id, fs.name AS name, fs.createdDate AS createdDate
		`, param)
		if err != nil {
			return nil, err
		}

		fs, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		var rz []ListResult
		for _, f := range fs {
			var rr ListResult
			rr.ID = f.AsMap()["id"].(string)
			rr.Name = f.AsMap()["name"].(string)
			rr.CreatedDate = f.AsMap()["createdDate"].(time.Time)
			rz = append(rz, rr)
		}
		return rz, nil
	})
	if err != nil {
		return nil, err
	}

	return fs.([]ListResult), nil
}

type ListMutualsResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) ListMutuals(ctx context.Context, param map[string]any) ([]ListMutualsResult, error) {
	mf, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name1})-[:FRIEND]->(mf)<-[:FRIEND]-(:User {name: $name2})
			RETURN mf.id AS id, mf.name AS name, mf.createdDate AS createdDate
		`, param)
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
	})
	if err != nil {
		return nil, err
	}

	return mf.([]ListMutualsResult), nil
}

func (f *Friend) MutualCount(ctx context.Context, param map[string]any) (int64, error) {
	cnt, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name1})-[:FRIEND]->(mf)<-[:FRIEND]-(:User {name: $name2})
			RETURN COUNT(mf) AS count
		`, param)
		if err != nil {
			return nil, err
		}

		cnt, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return cnt.AsMap()["count"].(int64), nil
	})
	if err != nil {
		return 0, err
	}
	return cnt.(int64), nil
}

func (f *Friend) Count(ctx context.Context, param map[string]any) (int64, error) {
	cnt, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name})-[:FRIEND {status: ['accept']}]->(fs)
			RETURN COUNT(fs) AS count
		`, param)
		if err != nil {
			return nil, err
		}

		cnt, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		return cnt.AsMap()["count"].(int64), nil
	})
	if err != nil {
		return 0, err
	}
	return cnt.(int64), nil
}

func (f *Friend) RequestCount(ctx context.Context, param map[string]any) (int64, error) {
	cnt, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $name})-[:FRIEND {status: ['request']}]->(fs)
			RETURN COUNT(fs) AS count
		`, param)
		if err != nil {
			return nil, err
		}

		cnt, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		return cnt.AsMap()["count"].(int64), nil
	})
	if err != nil {
		return 0, err
	}
	return cnt.(int64), nil
}

type Friend struct {
	sess neo4j.SessionWithContext
}

func NewFriend(sess neo4j.SessionWithContext) *Friend {
	return &Friend{
		sess: sess,
	}
}
