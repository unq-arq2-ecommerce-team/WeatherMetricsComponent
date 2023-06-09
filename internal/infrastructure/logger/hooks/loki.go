package hooks

import (
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

const AppName = "weather-metrics-component"
const LokiEndpoint = "http://loki:3100"

func BuildLokiHook() *lokirus.LokiHook {
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(&logrus.JSONFormatter{}).
		WithStaticLabels(lokirus.Labels{
			"app":         AppName,
			"environment": "development",
		})

	return lokirus.NewLokiHookWithOpts(
		LokiEndpoint,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel)
}
