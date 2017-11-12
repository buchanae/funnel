package events

import "golang.org/x/net/context"

type Service struct {
	Writer
}

func (s *Service) WriteEvent(ctx context.Context, e *Event) (*WriteEventResponse, error) {
	return &WriteEventResponse{}, s.Writer.WriteEvent(ctx, e)
}
