package web

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marciomarquesdesouza/go-rate-limiter/internal/infra/database"
	ratelimiter "github.com/marciomarquesdesouza/go-rate-limiter/internal/rate-limiter"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

type Webserver struct {
	RequestLimit         int
	BlockingTimeSeconds  int
	RequestInfoInterface *database.LimiterInfoRepositoryInterface
}

// NewServer creates a new server instance
func NewServer(requestLimit int, blockingTimeSeconds int, requestInfoInterface database.LimiterInfoRepositoryInterface) *Webserver {
	return &Webserver{
		RequestLimit:         requestLimit,
		BlockingTimeSeconds:  blockingTimeSeconds,
		RequestInfoInterface: &requestInfoInterface,
	}
}

// createServer creates a new server instance with go chi router
func (web *Webserver) CreateServer() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(web.RateLimiter)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Get("/", web.HandleRequest)
	return router
}

type TemplateData struct {
	Title              string
	BackgroundColor    string
	ResponseTime       time.Duration
	ExternalCallMethod string
	ExternalCallURL    string
	Content            string
	RequestNameOTEL    string
	OTELTracer         trace.Tracer
}

func (web *Webserver) HandleRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are allowed make some requests."))
}

func (web *Webserver) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Error getting host ip", http.StatusInternalServerError)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		apiKey := r.Header.Get("API_KEY")
		if apiKey != "" {
			web.RequestLimit = viper.GetInt(apiKey)
		}

		if blocked, err := ratelimiter.CheckLimitReached(host, web.RequestLimit, web.BlockingTimeSeconds, *web.RequestInfoInterface); blocked {
			http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
