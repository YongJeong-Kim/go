package service

import (
	"context"
	"gorelationship/repository"
)

type Friender interface {
	Accept(ctx context.Context, param map[string]any) error
	Count(ctx context.Context, param map[string]any) (int64, error)
	List(ctx context.Context, param map[string]any) ([]repository.ListResult, error)
	ListMutuals(ctx context.Context, param map[string]any) ([]repository.ListMutualsResult, error)
	ListRequests(ctx context.Context, param map[string]any) ([]repository.ListRequestsResult, error)
	MutualCount(ctx context.Context, param map[string]any) (int64, error)
	Request(ctx context.Context, param map[string]any) error
	RequestCount(ctx context.Context, param map[string]any) (int64, error)
}

func (f *Friend) Accept(ctx context.Context, param map[string]any) error {
	return f.Friend.Accept(ctx, param)
}

func (f *Friend) Count(ctx context.Context, param map[string]any) (int64, error) {
	return f.Friend.Count(ctx, param)
}

func (f *Friend) List(ctx context.Context, param map[string]any) ([]repository.ListResult, error) {
	return f.Friend.List(ctx, param)
}

func (f *Friend) ListMutuals(ctx context.Context, param map[string]any) ([]repository.ListMutualsResult, error) {
	return f.Friend.ListMutuals(ctx, param)
}

func (f *Friend) ListRequests(ctx context.Context, param map[string]any) ([]repository.ListRequestsResult, error) {
	return f.Friend.ListRequests(ctx, param)
}

func (f *Friend) MutualCount(ctx context.Context, param map[string]any) (int64, error) {
	return f.Friend.MutualCount(ctx, param)
}

func (f *Friend) Request(ctx context.Context, param map[string]any) error {
	return f.Friend.Request(ctx, param)
}

func (f *Friend) RequestCount(ctx context.Context, param map[string]any) (int64, error) {
	return f.Friend.RequestCount(ctx, param)
}

type Friend struct {
	Friend repository.Friender
}

func NewFriend(friender repository.Friender) *Friend {
	return &Friend{
		Friend: friender,
	}
}
