package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

const LokiEndpoint = "http://loki:3100"

func BuildLokiHook(conf *Config) *lokirus.LokiHook {
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(getFormatter(conf.LogFormat)).
		WithStaticLabels(lokirus.Labels{
			"app":         conf.ServiceName,
			"environment": conf.EnvironmentName,
		})

	return lokirus.NewLokiHookWithOpts(
		LokiEndpoint,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	)
}
