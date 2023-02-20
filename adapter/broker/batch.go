package broker

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
)

const (
	defaultBufferSize = 100
)

type (
	errorHandler func(
		ctx context.Context,
		err error,
		msgs []*kafka.Message,
		handler func(ctx context.Context, msg [][]byte, db rep.Storage) error,
	)

	batch struct {
		log          *zerolog.Logger
		handler      func(ctx context.Context, msg [][]byte, db rep.Storage) error
		errorHandler errorHandler
		buf          [][]byte
		msgsBuf      []*kafka.Message
		maxBufSize   int
	}
)

func newBatch(log zerolog.Logger, topic string, bufSize int,
	handler func(ctx context.Context, msg [][]byte, db rep.Storage) error) *batch {

	log = log.With().Str("cmp", "batch").Str(keyTopic, topic).Logger()

	if bufSize <= 0 {
		bufSize = defaultBufferSize
	}

	return &batch{
		log:        &log,
		handler:    handler,
		maxBufSize: bufSize,
		buf:        make([][]byte, 0, bufSize),
	}
}

func (b *batch) setErrorHandler(errorHandler errorHandler) {
	b.errorHandler = errorHandler
}

func (b *batch) insertMessage(ctx context.Context, msg *kafka.Message, db rep.Storage) {
	b.msgsBuf = append(b.msgsBuf, msg)
	b.buf = append(b.buf, msg.Value)
	if len(b.buf) >= b.maxBufSize {
		b.flushBuffer(ctx, db)
	}
}

func (b *batch) flushBuffer(ctx context.Context, db rep.Storage) {
	if len(b.buf) > 0 {
		start := time.Now()
		err := b.handler(ctx, b.buf, db)
		handleDur := time.Since(start)
		if b.errorHandler != nil {
			b.errorHandler(ctx, err, b.msgsBuf, b.handler)
		}
		b.resetBuf()
		b.log.Info().Dur("duration", handleDur).Int("count", len(b.buf)).Msg("buffer reset")
	}
}

func (b *batch) resetBuf() {
	b.msgsBuf = make([]*kafka.Message, 0, b.maxBufSize)
	b.buf = make([][]byte, 0, b.maxBufSize)
}
