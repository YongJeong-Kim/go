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
	Accept(ctx context.Context, requestUserID, approveUserID string) error
	Count(ctx context.Context, userID string) (int64, error)
	List(ctx context.Context, userID string) ([]repository.ListResult, error)
	ListMutuals(ctx context.Context, userID, friendUserID string) ([]repository.ListMutualsResult, error)
	ListRequests(ctx context.Context, userID string) ([]repository.ListRequestsResult, error)
	MutualCount(ctx context.Context, userID1, userID2 string) (int64, error)
	Request(ctx context.Context, requestUserID, approveUserID string) error
	RequestCount(ctx context.Context, userID string) (int64, error)
}

func (f *Friend) Accept(ctx context.Context, requestUserID, approveUserID string) error {
	if requestUserID == approveUserID {
		return errors.New("cannot accept yourself")
	}

	err := uuid.Validate(requestUserID)
	if err != nil {
		return errors.New("invalid request user uuid")
	}

	err = uuid.Validate(approveUserID)
	if err != nil {
		return errors.New("invalid approve user uuid")
	}

	_, err = neo4j.ExecuteWrite(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (struct{}, error) {
		err := f.Friend.Validate(ctx, tx, requestUserID, approveUserID)
		if err != nil {
			return struct{}{}, err
		}

		rs, err := f.Friend.RelationshipStatus(ctx, tx, requestUserID, approveUserID)
		if err != nil {
			return struct{}{}, err
		}

		switch {
		case rs.RequestUserID != nil && rs.AcceptUserID == nil:
			// ok accept
		case rs.RequestUserID == nil && rs.AcceptUserID != nil:
			return struct{}{}, fmt.Errorf("impossible case. check relationship req: %s, acc: %s", requestUserID, approveUserID)
		case rs.RequestUserID == nil && rs.AcceptUserID == nil:
			return struct{}{}, errors.New("You must first receive a friend request from " + requestUserID)
		case rs.RequestUserID != nil && rs.AcceptUserID != nil:
			return struct{}{}, errors.New("already friend")
		default:
			return struct{}{}, fmt.Errorf("impossible default case. check relationship req: %s, acc: %s", requestUserID, approveUserID)
		}

		created, err := f.Friend.Accept(ctx, tx, requestUserID, approveUserID)
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
	return f.Friend.Count(ctx, userID)
}

func (f *Friend) List(ctx context.Context, userID string) ([]repository.ListResult, error) {
	return f.Friend.List(ctx, userID)
}

func (f *Friend) ListMutuals(ctx context.Context, userID, friendUserID string) ([]repository.ListMutualsResult, error) {
	return f.Friend.ListMutuals(ctx, userID, friendUserID)
}

func (f *Friend) ListRequests(ctx context.Context, userID string) ([]repository.ListRequestsResult, error) {
	reqs, err := neo4j.ExecuteRead(ctx, f.Sess, func(tx neo4j.ManagedTransaction) ([]repository.ListRequestsResult, error) {
		reqs, err := f.Friend.ListRequests(ctx, tx, userID)
		if err != nil {
			return nil, err
		}

		return reqs, nil
	})
	if err != nil {
		return nil, err
	}
	return reqs, nil
	//return f.Friend.ListRequests(ctx, userID)
}

func (f *Friend) MutualCount(ctx context.Context, userID1, userID2 string) (int64, error) {
	return f.Friend.MutualCount(ctx, userID1, userID2)
}

func (f *Friend) Request(ctx context.Context, requestUserID, approveUserID string) error {
	fz, err := neo4j.ExecuteWrite(ctx, f.Sess, func(tx neo4j.ManagedTransaction) (int64, error) {
		err := f.Friend.Validate(ctx, tx, requestUserID, approveUserID)
		if err != nil {
			return 0, err
		}
		// TODO check RelationshipStatus
		err = f.Friend.Request(ctx, tx, requestUserID, approveUserID)
		if err != nil {
			return 0, err
		}
		return 0, nil
	})
	if err != nil {
		return err
	}
	_ = fz
	return nil
	//return f.Friend.Request(ctx, requestUserID, approveUserID)
}

func (f *Friend) RequestCount(ctx context.Context, userID string) (int64, error) {
	return f.Friend.RequestCount(ctx, userID)
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
