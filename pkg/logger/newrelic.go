package logger

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func StartTransaction(app *newrelic.Application, log *logrus.Logger, ctx context.Context, name string) (*newrelic.Transaction, context.Context) {
	txn := app.StartTransaction(name)
	ctx = newrelic.NewContext(ctx, txn)
	txnLogger := log.WithContext(ctx)
	ctx = context.WithValue(ctx, "log", txnLogger)

	return txn, ctx
}

func StartSegment(ctx context.Context, name string) *newrelic.Segment {
	var segment *newrelic.Segment

	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment = txn.StartSegment(name)
	}

	return segment
}
