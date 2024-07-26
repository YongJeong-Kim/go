package service

import (
	"context"
	"errors"
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

	return f.Friend.Accept(ctx, requestUserID, approveUserID)
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
	return f.Friend.ListRequests(ctx, userID)
}

func (f *Friend) MutualCount(ctx context.Context, userID1, userID2 string) (int64, error) {
	return f.Friend.MutualCount(ctx, userID1, userID2)
}

func (f *Friend) Request(ctx context.Context, requestUserID, approveUserID string) error {
	return f.Friend.Request(ctx, requestUserID, approveUserID)
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
