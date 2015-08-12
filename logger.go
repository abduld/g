package main // github.com/abduld/g

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

// LoggerMiddleware is a LoggerMiddleware handler that logs the request as it goes in and the response as it goes out.
type LoggerMiddleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger LoggerMiddleware
	Logger *logrus.Logger
	// Name is the name of the application as recorded in latency metrics
	Name string
}

// NewLoggerMiddleware returns a new *LoggerMiddleware, yay!
func NewLoggerMiddleware() *LoggerMiddleware {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	return &LoggerMiddleware{
		Logger: log,
		Name:   "web",
	}
}

func (l *LoggerMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	entry := l.Logger.WithFields(logrus.Fields{
		"request": r.RequestURI,
		"method":  r.Method,
		"remote":  r.RemoteAddr,
	})

	if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
		entry = entry.WithField("request_id", reqID)
	}
	entry.Info("started handling request")

	next(rw, r)

	latency := time.Since(start)
	res := rw.(negroni.ResponseWriter)
	entry.WithFields(logrus.Fields{
		"status":      res.Status(),
		"text_status": http.StatusText(res.Status()),
		"took":        latency,
		fmt.Sprintf("measure#%s.latency", l.Name): latency.Nanoseconds(),
	}).Info("completed handling request")
}
