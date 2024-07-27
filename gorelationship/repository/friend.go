package repository

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type Friender interface {
	Accept(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (int, error)
	Count(ctx context.Context, name string) (int64, error)
	List(ctx context.Context, userID string) ([]ListResult, error)
	ListMutuals(ctx context.Context, userID, friendUserID string) ([]ListMutualsResult, error)
	ListRequests(ctx context.Context, tx neo4j.ManagedTransaction, userID string) ([]ListRequestsResult, error)
	MutualCount(ctx context.Context, userID1, userID2 string) (int64, error)
	Request(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) error
	RequestCount(ctx context.Context, userID string) (int64, error)
	Validate(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) error
	RelationshipStatus(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (*RelationshipStatusResult, error)
}

type ListRequestsResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) ListRequests(ctx context.Context, tx neo4j.ManagedTransaction, userID string) ([]ListRequestsResult, error) {
	result, err := tx.Run(ctx, `
		MATCH (:User {id: $userID})-[:FRIEND]->(requests) 
		RETURN requests.id AS id, requests.name AS name, requests.createdDate AS createdDate
	`, map[string]any{
		"userID": userID,
	})
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
	if result.Err() != nil {
		return nil, err
	}

	return rz, nil

	/*requests, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {id: $userID})-[:FRIEND]->(requests)
			RETURN requests.id AS id, requests.name AS name, requests.createdDate AS createdDate
		`, map[string]any{
			"userID": userID,
		})
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

	return requests.([]ListRequestsResult), nil*/
}

func (f *Friend) Request(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) error {
	_, err := tx.Run(ctx, `
		MATCH (r:User {id: $request})
		MATCH (a:User {id: $approve})
		MERGE (r)-[:FRIEND {status: 'request'}]->(a)
	`, map[string]any{
		"request": requestUserID,
		"approve": approveUserID,
	})
	if err != nil {
		return err
	}

	return nil

	/*_, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		rc, err := tx.Run(ctx, `
			MATCH (r:User) WHERE r.id = $request
			MATCH (a:User) WHERE a.id = $approve
			RETURN r.id AS rID, a.id AS aID
		`, map[string]any{
			"request": requestUserID,
			"approve": approveUserID,
		})
		if err != nil {
			return nil, err
		}

		r, err := rc.Single(ctx)
		if err != nil {
			return nil, err
		}

		rID := r.AsMap()["rID"]
		if rID == nil {
			return nil, errors.New("request user not found")
		}

		aID := r.AsMap()["aID"]
		if aID == nil {
			return nil, errors.New("approve user not found")
		}

		_, err = tx.Run(ctx, `
			MATCH (r:User {id: $request})
			MATCH (a:User {id: $approve})
			MERGE (r)-[:FRIEND {status: ['request']}]->(a)
		`, map[string]any{
			"request": requestUserID,
			"approve": approveUserID,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil*/
}

func (f *Friend) Accept(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (int, error) {
	result, err := tx.Run(ctx, `
			MATCH (r:User) WHERE r.id = $request
			MATCH (a:User) WHERE a.id = $approve
			MERGE (a)-[:FRIEND {status: ['accept']}]->(r)
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

	/*_, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		rc, err := tx.Run(ctx, `
			MATCH (r:User) WHERE r.id = $request
			MATCH (a:User) WHERE a.id = $approve
			RETURN r.id AS rID, a.id AS aID
		`, map[string]any{
			"request": requestUserID,
			"approve": approveUserID,
		})
		if err != nil {
			return nil, err
		}

		r, err := rc.Single(ctx)
		if err != nil {
			return nil, err
		}

		rID := r.AsMap()["rID"]
		if rID == nil {
			return nil, errors.New("request user not found")
		}

		aID := r.AsMap()["aID"]
		if aID == nil {
			return nil, errors.New("approve user not found")
		}

		rc, err = tx.Run(ctx, `
			MATCH (r:User {id: $request})
			MATCH (a:User {id: $approve})
			MATCH (r)-[f1:FRIEND {status: ['request']}]->(a)
			OPTIONAL MATCH (a)-[f2:FRIEND {status: ['accept']}]->(r)
			RETURN f1.status AS request, f2.status AS accept
		`, map[string]any{
			"request": requestUserID,
			"approve": approveUserID,
		})
		if err != nil {
			return nil, err
		}

		r, err = rc.Single(ctx)
		if err != nil {
			return nil, errors.New("You must first receive a friend request from " + requestUserID)
		}

		request := r.AsMap()["request"]
		if request == nil {
			return nil, errors.New("impossible case: " + requestUserID)
		}

		accept := r.AsMap()["accept"]
		if accept != nil {
			return nil, errors.New("already friend")
		}

		_, err = tx.Run(ctx, `
			MATCH (r:User) WHERE r.id = $request
			MATCH (a:User) WHERE a.id = $approve
			MERGE (a)-[:FRIEND {status: ['accept']}]->(r)
		`, map[string]any{
			"request": requestUserID,
			"approve": approveUserID,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil*/
}

type ListResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (f *Friend) List(ctx context.Context, userID string) ([]ListResult, error) {
	fs, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {id: $userID})-[:FRIEND {status: ['accept']}]->(fs)
			RETURN fs.id AS id, fs.name AS name, fs.createdDate AS createdDate
		`, map[string]any{
			"userID": userID,
		})
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

func (f *Friend) ListMutuals(ctx context.Context, userID, friendUserID string) ([]ListMutualsResult, error) {
	mf, err := f.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {id: $userID})-[:FRIEND]->(mf)<-[:FRIEND]-(:User {id: $friendUserID})
			RETURN mf.id AS id, mf.name AS name, mf.createdDate AS createdDate
		`, map[string]any{
			"userID":       userID,
			"friendUserID": friendUserID,
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
	})
	if err != nil {
		return nil, err
	}

	return mf.([]ListMutualsResult), nil
}

func (f *Friend) MutualCount(ctx context.Context, userID1, userID2 string) (int64, error) {
	cnt, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {id: $userID1})-[:FRIEND]->(mf)<-[:FRIEND]-(:User {userID: $userID2})
			RETURN COUNT(mf) AS count
		`, map[string]any{
			"userID1": userID1,
			"userID2": userID2,
		})
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

func (f *Friend) Count(ctx context.Context, userID string) (int64, error) {
	cnt, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {name: $userID})-[:FRIEND {status: ['accept']}]->(fs)
			RETURN COUNT(fs) AS count
		`, map[string]any{
			"userID": userID,
		})
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

func (f *Friend) RequestCount(ctx context.Context, userID string) (int64, error) {
	cnt, err := f.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (:User {id: $userID})-[:FRIEND {status: 'request'}]->(fs)
			RETURN COUNT(fs) AS count
		`, map[string]any{
			"userID": userID,
		})
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

type RelationshipStatusResult struct {
	RequestUserID *string
	AcceptUserID  *string
}

func (f *Friend) RelationshipStatus(ctx context.Context, tx neo4j.ManagedTransaction, requestUserID, approveUserID string) (*RelationshipStatusResult, error) {
	result, err := tx.Run(ctx, `
		MATCH (r:User {id: $request})
		MATCH (a:User {id: $approve})
		OPTIONAL MATCH (r)-[f1:FRIEND {status: 'request'}]->(a)
		OPTIONAL MATCH (a)-[f2:FRIEND {status: 'accept'}]->(r)
		RETURN f1.status AS request, f2.status AS accept
	`, map[string]any{
		"request": requestUserID,
		"approve": approveUserID,
	})
	if err != nil {
		return nil, err
	}

	rs, err := result.Single(ctx)
	if err != nil {
		return nil, err
		//return nil, errors.New("You must first receive a friend request from " + requestUserID)
	}

	status := &RelationshipStatusResult{
		RequestUserID: nil,
		AcceptUserID:  nil,
	}

	rID := rs.AsMap()["request"]
	if rID != nil {
		p := rID.(string)
		status.RequestUserID = &p
	}

	aID := rs.AsMap()["accept"]
	if aID != nil {
		p := aID.(string)
		status.AcceptUserID = &p
	}

	return status, nil
}

type Friend struct {
	sess neo4j.SessionWithContext
}

func NewFriend(sess neo4j.SessionWithContext) *Friend {
	return &Friend{
		sess: sess,
	}
}
