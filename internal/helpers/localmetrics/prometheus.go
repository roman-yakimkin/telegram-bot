package localmetrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"go.uber.org/zap"
)

var (
	CntMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tg_bot",
			Name:      "user_msg_total",
			Help:      "The total messages of user by message type",
		},
		[]string{"user_id", "message_type"},
	)

	PerformDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "tg_bot",
			Name:      "msg_performing_duration",
			Help:      "Duration of messages performing",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.95: 0.005,
				0.99: 0.001,
			},
		},
		[]string{"user_id", "message_type"},
	)
)

func HandleMetrics(cfg config.Config, logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(cfg.PrometheusMetricsURL, nil); err != nil {
		logger.Error("error listening prometheus handler", zap.Error(err))
		return
	}
}
