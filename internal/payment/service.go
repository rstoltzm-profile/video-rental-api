package payment

import "context"

type Service interface {
	MakePayment(ctx context.Context, req Payment) (int, error)
}

type service struct {
	writer PaymentWriter
}

func NewService(writer PaymentWriter) Service {
	return &service{
		writer: writer,
	}
}

func (s *service) MakePayment(ctx context.Context, req Payment) (int, error) {
	return s.writer.InsertPayment(ctx, req)
}
