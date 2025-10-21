package publishers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/services/broker"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/events"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

const authErrorTracer string = "publisher.auth"

type AuthEventPublisher interface {
	PublishAuthRegistered(ctx context.Context, authID int64, email, token string) (err error)
}

type authEventPublisher struct {
	producer *broker.Producer
	appName  string
}

func NewAuthEventPublisher(producer *broker.Producer, appName string) AuthEventPublisher {
	return &authEventPublisher{producer, appName}
}

func (e *authEventPublisher) PublishAuthRegistered(ctx context.Context, authID int64, email, token string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "PublishAuthRegistered")
	defer span.End()

	et := constants.EventTypeAuthRegistered
	data := events.AuthRegistered{
		Recipient: email,
		Token:     token,
	}

	bytes, err := proto.Marshal(&data)
	if err != nil {
		wErr := fmt.Errorf("failed to publish event %s: %w", et, err)
		return ce.NewError(span, ce.CodeEventPublishingFailed, ce.MsgInternalServer, wErr)
	}

	key := fmt.Sprintf("auth-%d", authID)
	event := events.Event{
		EventId:       utils.NewUUID().String(),
		EventType:     et,
		SourceService: e.appName,
		Timestamp:     time.Now().UTC().UnixMilli(),
		Data:          bytes,
	}

	if err := e.producer.Publish(ctx, "auth-events", key, &event); err != nil {
		wErr := fmt.Errorf("failed to publish event %s: %w", et, err)
		return ce.NewError(span, ce.CodeEventPublishingFailed, ce.MsgInternalServer, wErr)
	}

	log.Printf("event %s published successfully", et)
	return nil
}
