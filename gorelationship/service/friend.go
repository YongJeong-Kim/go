package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorelationship/repository"
)

type Friender interface {
	Accept(ctx context.Context, requestUserID, acceptUserID string) error
	Count(ctx context.Context, userID string) (int64, error)
	List(ctx context.Context, userID string) ([]repository.ListResult, error)
	ListMutuals(ctx context.Context, userID1, userID2 string) ([]repository.ListMutualsResult, error)
	ListFromRequests(ctx context.Context, userID string) ([]repository.ListFromRequestsResult, error)
	MutualCount(ctx context.Context, userID1, userID2 string) (int64, error)
	Request(ctx context.Context, requestUserID, acceptUserID string) error
	FromRequestCount(ctx context.Context, userID string) (int64, error)
}

func (f *Friend) Accept(ctx context.Context, requestUserID, acceptUserID string) error {
	if requestUserID == acceptUserID {
		return errors.New("cannot accept yourself")
	}

	err := uuid.Validate(requestUserID)
	if err != nil {
		return errors.New("invalid request user uuid")
	}

	err = uuid.Validate(acceptUserID)
	if err != nil {
		return errors.New("invalid approve user uuid")
	}

	_, err = neo4j.ExecuteWrite(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (struct{}, error) {
		err := f.Friend.Validate(ctx, tx, requestUserID, acceptUserID)
		if err != nil {
			return struct{}{}, err
		}

		rs, err := f.Friend.RelationshipStatus(ctx, tx, requestUserID, acceptUserID)
		if err != nil {
			return struct{}{}, err
		}

		switch rs {
		case "NO":
			return struct{}{}, errors.New("You must first receive a friend request from " + requestUserID)
		case "REQUEST":
			// ok
		case "FRIEND":
			return struct{}{}, errors.New("already friend")
		default:
			return struct{}{}, fmt.Errorf("impossible case. check relationship req: %s, acc: %s", requestUserID, acceptUserID)
		}

		created, err := f.Friend.Accept(ctx, tx, requestUserID, acceptUserID)
		if err != nil {
			return struct{}{}, err
		}
		if created == 0 {
			return struct{}{}, errors.New("already friend")
		}

		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (f *Friend) Count(ctx context.Context, userID string) (int64, error) {
	err := uuid.Validate(userID)
	if err != nil {
		return 0, errors.New("invalid user uuid")
	}

	cnt, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (int64, error) {
		cnt, err := f.Friend.Count(ctx, tx, userID)
		if err != nil {
			return 0, err
		}

		return cnt, nil
	})
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (f *Friend) List(ctx context.Context, userID string) ([]repository.ListResult, error) {
	err := uuid.Validate(userID)
	if err != nil {
		return nil, errors.New("invalid user uuid")
	}

	fs, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) ([]repository.ListResult, error) {
		fs, err := f.Friend.List(ctx, tx, userID)
		if err != nil {
			return nil, err
		}

		return fs, nil
	})
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func (f *Friend) ListMutuals(ctx context.Context, userID1, userID2 string) ([]repository.ListMutualsResult, error) {
	mutuals, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) ([]repository.ListMutualsResult, error) {
		mutuals, err := f.Friend.ListMutuals(ctx, tx, userID1, userID2)
		if err != nil {
			return nil, err
		}

		return mutuals, nil
	})
	if err != nil {
		return nil, err
	}

	return mutuals, nil
}

func (f *Friend) ListFromRequests(ctx context.Context, userID string) ([]repository.ListFromRequestsResult, error) {
	if err := uuid.Validate(userID); err != nil {
		return nil, errors.New("invalid user uuid")
	}

	reqs, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) ([]repository.ListFromRequestsResult, error) {
		reqs, err := f.Friend.ListFromRequests(ctx, tx, userID)
		if err != nil {
			return nil, err
		}

		return reqs, nil
	})
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (f *Friend) MutualCount(ctx context.Context, userID1, userID2 string) (int64, error) {
	if err := uuid.Validate(userID1); err != nil {
		return 0, errors.New("invalid user1 uuid")
	}
	if err := uuid.Validate(userID2); err != nil {
		return 0, errors.New("invalid user2 uuid")
	}

	cnt, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (int64, error) {
		cnt, err := f.Friend.MutualCount(ctx, tx, userID1, userID2)
		if err != nil {
			return 0, err
		}

		return cnt, nil
	})
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (f *Friend) Request(ctx context.Context, requestUserID, acceptUserID string) error {
	if requestUserID == acceptUserID {
		return errors.New("cannot request yourself")
	}

	err := uuid.Validate(requestUserID)
	if err != nil {
		return errors.New("invalid request user uuid")
	}

	err = uuid.Validate(acceptUserID)
	if err != nil {
		return errors.New("invalid accept user uuid")
	}

	_, err = neo4j.ExecuteWrite(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (struct{}, error) {
		err := f.Friend.Validate(ctx, tx, requestUserID, acceptUserID)
		if err != nil {
			return struct{}{}, err
		}

		rs, err := f.Friend.RelationshipStatus(ctx, tx, requestUserID, acceptUserID)
		if err != nil {
			return struct{}{}, err
		}

		switch rs {
		case "NO":
			// ok
		case "REQUEST":
			return struct{}{}, errors.New("already send request")
		case "FRIEND":
			return struct{}{}, errors.New("already friend")
		default:
			return struct{}{}, fmt.Errorf("impossible case. check relationship req: %s, acc: %s", requestUserID, acceptUserID)
		}

		created, err := f.Friend.Request(ctx, tx, requestUserID, acceptUserID)
		if err != nil {
			return struct{}{}, err
		}
		if created == 0 {
			return struct{}{}, errors.New("send request fail")
		}

		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (f *Friend) FromRequestCount(ctx context.Context, userID string) (int64, error) {
	err := uuid.Validate(userID)
	if err != nil {
		return 0, errors.New("invalid user uuid")
	}

	count, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (int64, error) {
		count, err := f.Friend.FromRequestCount(ctx, tx, userID)
		if err != nil {
			return 0, err
		}

		return count, nil
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

type Friend struct {
	Sess   neo4j.SessionWithContext
	Friend repository.Friender
}

func NewFriend(sess neo4j.SessionWithContext, friender repository.Friender) *Friend {
	return &Friend{
		Sess:   sess,
		Friend: friender,
	}
}
