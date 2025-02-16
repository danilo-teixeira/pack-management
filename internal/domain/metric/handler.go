package metric

import (
	"pack-management/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/bun"
)

type (
	HandlerParams struct {
		App *fiber.App `validate:"required"`
		DB  *bun.DB    `validate:"required"`
	}

	handler struct {
		app *fiber.App
		db  *bun.DB
	}
)

// // Pool Status
// OpenConnections int // The number of established connections both in use and idle.
// InUse           int // The number of connections currently in use.
// Idle            int // The number of idle connections.

// // Counters
// WaitCount         int64         // The total number of connections waited for.
// WaitDuration      time.Duration // The total time blocked waiting for a new connection.
// MaxIdleClosed     int64         // The total number of connections closed due to SetMaxIdleConns.
// MaxIdleTimeClosed int64         // The total number of connections closed due to SetConnMaxIdleTime.
// MaxLifetimeClosed int64         // The total number of connections closed due to SetConnMaxLifetime.

var (
	dbMaxOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_max_open_connections",
		Help: "Maximum number of open connections to the database.",
	})
	dbOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_open_connections",
		Help: "The number of established connections both in use and idle.",
	})
	dbInUseConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_in_use_connections",
		Help: "The number of connections currently in use.",
	})
	dbIdleConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_idle_connections",
		Help: "The number of idle connections.",
	})
	dbWaitCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_wait_count",
		Help: "The total number of connections waited for.",
	})
	dbWaitDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_wait_duration",
		Help: "The total time blocked waiting for a new connection.",
	})
	dbMaxIdleClosed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_max_idle_closed",
		Help: "The total number of connections closed due to SetMaxIdleConns.",
	})
	dbMaxIdleTimeClosed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_max_idle_time_closed",
		Help: "The total number of connections closed due to SetConnMaxIdleTime.",
	})
	dbMaxLifetimeClosed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_stats_max_lifetime_closed",
		Help: "The total number of connections closed due to SetConnMaxLifetime.",
	})
)

func NewHTPPHandler(params *HandlerParams) *handler {
	params.validate()

	h := &handler{
		app: params.App,
		db:  params.DB,
	}

	prometheus.MustRegister(
		dbMaxOpenConnections,
		dbOpenConnections,
		dbInUseConnections,
		dbIdleConnections,
		dbWaitCount,
		dbWaitDuration,
		dbMaxIdleClosed,
		dbMaxIdleTimeClosed,
		dbMaxLifetimeClosed,
	)

	h.app.Get("/metrics", h.listMetrics)

	return h
}

func (p *HandlerParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (h *handler) listMetrics(ctx *fiber.Ctx) error {
	dbStats := h.db.Stats()
	dbMaxOpenConnections.Set(float64(dbStats.MaxOpenConnections))
	dbOpenConnections.Set(float64(dbStats.OpenConnections))
	dbInUseConnections.Set(float64(dbStats.InUse))
	dbIdleConnections.Set(float64(dbStats.Idle))
	dbWaitCount.Set(float64(dbStats.WaitCount))
	dbWaitDuration.Set(float64(dbStats.WaitDuration))
	dbMaxIdleClosed.Set(float64(dbStats.MaxIdleClosed))
	dbMaxIdleTimeClosed.Set(float64(dbStats.MaxIdleTimeClosed))
	dbMaxLifetimeClosed.Set(float64(dbStats.MaxLifetimeClosed))

	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
