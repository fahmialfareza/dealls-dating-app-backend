package newrelic

import (
	"os"
	"time"

	"github.com/cenkalti/backoff"
	nrConfig "github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func NewNewRelic(appName string, license string) (*NewRelic, error) {
	var (
		app *newrelic.Application
		log *logrus.Logger
	)

	conn := func() (err error) {
		log = logrus.New()
		log.SetLevel(logrus.InfoLevel)

		app, err = newrelic.NewApplication(
			newrelic.ConfigAppName(appName),
			newrelic.ConfigLicense(license),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigAppLogEnabled(true),
			nrlogrus.ConfigLogger(log),
		)
		if err != nil {
			return err
		}

		nrlogrusFormatter := nrConfig.NewFormatter(app, &logrus.JSONFormatter{})
		log.SetFormatter(nrlogrusFormatter)
		log.Out = os.Stdout

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(conn, expBackoff)
	if err != nil {
		return &NewRelic{}, err
	}

	return &NewRelic{
		App: app,
		Log: log,
	}, nil
}
