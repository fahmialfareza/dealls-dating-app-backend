package newrelic

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type NewRelic struct {
	App *newrelic.Application
	Log *logrus.Logger
}
