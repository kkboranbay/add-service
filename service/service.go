package service

import "context"

type Service interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

func (basicService) Sum(ctx context.Context, a, b int) (int, error) {
	return a + b, nil
}

func (basicService) Concat(ctx context.Context, a, b string) (string, error) {
	return a + b, nil
}
