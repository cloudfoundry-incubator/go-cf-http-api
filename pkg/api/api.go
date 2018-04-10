package api

import (
	"log"
	"net"
	"net/http"
	"time"

	"encoding/json"
	"io/ioutil"

	"context"
	"strings"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/auth"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/uaaclient"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/viewer"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
)

type UAAClient interface {
	CheckToken(token string, ctx context.Context) (*uaaclient.User, error)
}

type Server struct {
	httpServer         *http.Server
	uaaClient          UAAClient
	Endpoints          []*Endpoint
	hostname           string
	templatesDirectory string
	logRequest         requestLogger
}

type requestLogger func(req Request, resp Response, endpoint *Endpoint, startTime time.Time, totalTime time.Duration)

type Config struct {
	UAAClient          UAAClient
	Hostname           string
	Port               string
	Endpoints          []*Endpoint
	LogRequest         requestLogger
	TemplatesDirectory string
	AssetsDirectory    string
}

func New(apiConfig Config) *Server {
	router := mux.NewRouter()
	router.Handle("/debug/metrics", exp.ExpHandler(metrics.DefaultRegistry)).Methods("GET")

	if apiConfig.LogRequest == nil {
		apiConfig.LogRequest = func(req Request, resp Response, endpoint *Endpoint, startTime time.Time, totalTime time.Duration) {}
	}

	server := &Server{
		httpServer: &http.Server{
			Addr:    net.JoinHostPort("", apiConfig.Port),
			Handler: router,
		},
		uaaClient:          apiConfig.UAAClient,
		logRequest:         apiConfig.LogRequest,
		hostname:           apiConfig.Hostname,
		templatesDirectory: apiConfig.TemplatesDirectory,
	}

	router.Handle("/assets/{rest}", http.StripPrefix("/assets/", http.FileServer(http.Dir(apiConfig.AssetsDirectory))))
	for _, e := range apiConfig.Endpoints {
		router.Handle(e.Path, server.handle(e)).Methods(e.Method)
	}

	return server
}

func (s *Server) Start() func() {
	go s.httpServer.ListenAndServe()

	log.Println("HTTP metrics now available at /debug/metrics")

	return func() {
		s.httpServer.Close()
	}
}

func (s *Server) handle(endpoint *Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "bearer ")
		token = strings.TrimPrefix(token, "Bearer ")

		currentUser, _ := s.uaaClient.CheckToken(token, context.Background())
		req := &realRequest{
			httpRequest: r,
			currentUser: currentUser,
		}

		start := time.Now()
		var resp Response

		if s.passesAuth(endpoint.Auth, currentUser) {
			resp = *endpoint.Handle(req)
		} else {
			resp = Response{
				StatusCode: http.StatusUnauthorized,
			}
		}

		s.writeResponse(w, resp)
		s.logRequest(req, resp, endpoint, start, time.Since(start))
	}
}

func (s *Server) writeResponse(w http.ResponseWriter, resp Response) {
	var bodyBytes []byte
	status := resp.StatusCode

	if resp.Body != nil {
		switch body := resp.Body.(type) {

		case HTMLTemplate:
			f, err := ioutil.ReadFile(s.templatesDirectory + "/" + body.Name + ".html")
			if err != nil {
				status = http.StatusInternalServerError
			} else {
				bodyBytes = []byte(viewer.Parse(string(f), map[string]string{"hostname": s.hostname}))
				w.Header().Set("Content-Type", "text/html")
			}

		default:
			b, err := json.Marshal(resp.Body)
			if err != nil {
				status = http.StatusInternalServerError
			} else {
				bodyBytes = b
				w.Header().Set("Content-Type", "application/json")
			}
		}
	}

	w.WriteHeader(status)
	w.Write(bodyBytes)
}

func (s *Server) passesAuth(authConfig *auth.Config, currentUser *uaaclient.User) bool {
	if authConfig.AuthType != auth.NONE {
		if currentUser == nil {
			return false
		}

		for _, s := range authConfig.Scopes {
			for _, authScope := range currentUser.Scopes {
				if s == authScope {
					return true
				}
			}
		}

		// none of the user's scopes matched the expected scopes
		if len(authConfig.Scopes) > 0 {
			return false
		}
	}

	return true
}

type HTMLTemplate struct {
	Name string
}
