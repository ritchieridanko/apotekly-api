package broker

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	producer *kafka.Writer
}

func NewProducer(producer *kafka.Writer) *Producer {
	return &Producer{producer}
}

func (p *Producer) Publish(ctx context.Context, topic, key string, event proto.Message) error {
	bytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	requestID := fmt.Sprintf("%s", ctx.Value(constants.CtxKeyRequestID))

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: bytes,
		Headers: []kafka.Header{
			{Key: "trace_id", Value: []byte(traceID)},
			{Key: "correlation_id", Value: []byte(requestID)},
			{Key: "content_type", Value: []byte("application/x-protobuf")},
		},
	}

	return p.producer.WriteMessages(ctx, message)
}
