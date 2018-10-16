package api

import (
	"net/http"
	"strings"
	"time"

	"nyota/backend/api/requestinterceptor"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/store"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricReqCount *prometheus.CounterVec // metric - number of requests/tenant
	metricReqTimes *prometheus.SummaryVec // metric - time per req
)

type Service struct {
	Router *mux.Router
	Store  *store.Store
}

//InitAPI - initialize in api package
func initAPI() {
	metricReqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_req_total",
			Help: "Total number of API reqs processed, partitioned by tenantID and url",
		},
		[]string{"tenantID", "url"},
	)

	metricReqTimes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "api_req_time",
			Help:       "Total time taken to process API requests, partitioned by tenantID and url",
			Objectives: map[float64]float64{0.25: 0.05, 0.5: 0.05, 0.75: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"tenantID", "url"},
	)

	prometheus.MustRegister(metricReqCount)
	prometheus.MustRegister(metricReqTimes)
}

// Chain applies Prizm Handler to a http.HandlerFunc
func chain(realFunc requestinterceptor.PrizmHandler,
	interceptors ...requestinterceptor.Interceptor) http.HandlerFunc {

	// Chaining of all requests....
	for _, m := range interceptors {
		realFunc = m(realFunc)
	}

	// Creating a Session context which will be passed to chaining...
	u := &model.SessionContext{User: &model.UserContext{TenantId: "", UserName: ""}, Err: nil}

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// Invoke the chaining...
		realFunc(u, w, r)

		if u.Err != nil {
			logutil.Errorf(u, "Method:%s, URL:%s, Type:%s, Message: %s", r.Method, r.URL, u.Err.Type, u.Err.Message)
			http.Error(w, u.Err.Message, u.Err.Code)
			//u.Err = nil
		}

		//add url tenant details to prom service
		if !strings.Contains(r.URL.Path, "/metrics") {
			timeTaken := time.Since(startTime)
			metricReqTimes.WithLabelValues(u.User.TenantId, r.URL.Path).Observe(float64(timeTaken / time.Millisecond))
			metricReqCount.WithLabelValues(u.User.TenantId, r.URL.Path).Inc()
		}
	}
}

/*NewRoute Adds all routes exposed by ABS*/
func NewRoute() *mux.Router {

	store, err := store.New()
	if err != nil {
		return nil
	}

	srv := &Service{
		Router: mux.NewRouter(),
		Store:  store,
	}
	initAPI()
	// Add user records to db
	// srv.addRecords()

	r := srv.Router

	//added for withoutPrefix route
	r.Handle("/metrics", promhttp.Handler())

	// All API's to use this...
	apiRoute := r.PathPrefix("/api/v1/").Subrouter()

	nologinRoutes, guardedRoutes := getAllRoutes(srv)
	for _, route := range nologinRoutes {
		apiRoute.Handle(route.Path, chain(route.RealHandler,
			requestinterceptor.TrackReqResp(),
			requestinterceptor.AddNoCacheHeader())).Methods(route.Method)
	}

	for _, route := range guardedRoutes {
		apiRoute.Handle(route.Path, chain(route.RealHandler,
			requestinterceptor.RBACCheck(route.Group, route.Permission),
			requestinterceptor.TrackReqResp(),
			requestinterceptor.AddNoCacheHeader(),
			requestinterceptor.ValidateSession())).Methods(route.Method)
	}

	// This will serve static html files
	fileServer := http.FileServer(http.Dir("static/"))

	oldWebHandler := http.StripPrefix("/ui/", fileServer)
	r.PathPrefix("/ui/").Handler(oldWebHandler)

	webHandler := http.StripPrefix("/", fileServer)
	r.PathPrefix("/").Handler(webHandler)

	return r
}
