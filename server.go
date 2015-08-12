package main // github.com/abduld/g

import (
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/tylerb/graceful.v1"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/gorilla/context"
	e "github.com/pjebs/jsonerror"
	//github.com/jinzhu/gorm
	//https://github.com/pjebs/optimus-go
	//  "github.com/xyproto/permissions2"
	//"github.com/mholt/binding"
	//github.com/throttled/throttled
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hd_errors_total",
		Help: "Number of hard-disk errors.",
	})
)

func C(r *http.Request, authenticatedKey string) {
	context.Set(r, 0, authenticatedKey) // Read http://www.gorillatoolkit.org/pkg/context about setting arbitary context key
}

func addStaticRoutes(routes *mux.Router) {
	routes.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/",
		http.FileServer(http.Dir("public"))))
}

func NewRoute() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)
	r := render.New(render.Options{
		Directory:     "templates",
		IndentJSON:    true,
		IsDevelopment: true,
		Layout:        "layout",
	})

	routes.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.Text(w, http.StatusOK, "Index page")
	})

	routes.HandleFunc("/template", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", nil)
	})

	routes.HandleFunc("/error", func(w http.ResponseWriter, req *http.Request) {

		err := e.New(12, "Unauthorized Access", "Please log in first to access this site")
		r.JSON(w, http.StatusUnauthorized, err.Render())
	})

	routes.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		r.Data(w, http.StatusOK, []byte("Some binary data here."))
	})

	routes.Handle("/metrics", prometheus.Handler())

	addStaticRoutes(routes)

	return routes
}

/* Add with
app.UseHandler(NewRecoveryMiddleware())
func NewRecoveryMiddleware() http.Handler {
	recoveryMiddleware := recovery.New(recovery.Options{
		Out:       os.Stderr,
		StackSize: 8 * 1024,
		Prefix:    "Recovery",
	})
	recoveryMiddleware.Logger = nil
	return recoveryMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("you should not have a handler that just panics ;)")
	}))
}
*/

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
	})

	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(NewLoggerMiddleware())
	//app.Use(negroni.NewStatic(http.Dir("public")))
	app.UseHandler(context.ClearHandler(NewRoute()))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	graceful.Run(":"+port, 10*time.Second, app)
}
