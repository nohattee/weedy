package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
)

var validate = validator.New()

type httpServerCtxKey int

const (
	requestKey httpServerCtxKey = iota
)

type Response struct {
	Err        error       `json:"-"`
	StatusCode int         `json:"-"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
}

func RequestFromContext(ctx context.Context) interface{} {
	v := ctx.Value(requestKey)
	return v
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

type Route struct {
	Name    string
	Path    string
	Method  string
	Request interface{}
	Handler func(context.Context) Response
	Group   []Route
}

type Controller interface {
	Routes() []Route
}

type HttpServer struct {
	host        string
	port        string
	router      chi.Router
	controllers []Controller
}

func NewHttpServer(host string, port string, path string) *HttpServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return &HttpServer{
		host:        host,
		port:        port,
		router:      r,
		controllers: make([]Controller, 0),
	}
}

func (s *HttpServer) AddController(ctrl Controller) {
	s.controllers = append(s.controllers, ctrl)
}

func extractParams(r *http.Request, mapData map[string]interface{}) {
	queryParams := r.URL.Query()
	for k, v := range queryParams {
		// Check param's value if it is JSON.
		// Sometimes we will pass JSON in GET request like /users?filter={"name": "Admin"}.
		var value map[string]interface{}
		if err := json.Unmarshal([]byte(v[0]), &value); err == nil {
			mapData[k] = value
			continue
		}
		mapData[k] = v[0]
	}

	// Extract params in url like: /posts/{postId}/comments/{commentId}
	urlParam := chi.RouteContext(r.Context()).URLParams

	// Key, Value pair start with index=1
	for i := 1; i < len(urlParam.Keys); i++ {
		mapData[urlParam.Keys[i]] = urlParam.Values[i]
	}
}

func extractBody(r *http.Request, mapData map[string]interface{}) {
	r.ParseMultipartForm(32 << 20)

	if r.Form != nil {
		for k, v := range r.Form {
			mapData[k] = v
		}
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		for k, v := range r.MultipartForm.File {
			mapData[k] = v
		}
	}

	json.NewDecoder(r.Body).Decode(&mapData)
}

func parseRequest(route Route, r *http.Request) error {
	mapRequest := map[string]interface{}{}

	if r.Method != http.MethodGet {
		extractBody(r, mapRequest)
	}
	extractParams(r, mapRequest)
	byteRequest, err := json.Marshal(mapRequest)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(byteRequest, &route.Request); err != nil {
		return err
	}
	return validate.Struct(route.Request)
}

func initRoutes(r chi.Router, routes []Route) {
	for i := range routes {
		route := routes[i]
		if len(route.Group) > 0 {
			r.Route(route.Path, func(r chi.Router) {
				initRoutes(r, route.Group)
			})
		}
		if route.Handler != nil {
			r.Method(route.Method, route.Path, http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					err := parseRequest(route, r)
					if err != nil {
						render.Render(w, r, &Response{
							Err:        err,
							StatusCode: http.StatusUnprocessableEntity,
							Message:    "UNPROCESSABLE_ENTITY",
						})
						return
					}

					ctx := context.WithValue(r.Context(), requestKey, route.Request)

					resp := route.Handler(ctx)
					if resp.Err != nil {
						resp.StatusCode = http.StatusInternalServerError
						render.Render(w, r, &resp)
						return
					}
					render.Render(w, r, &resp)
				}))
		}
	}
}

func (s *HttpServer) Run() error {
	for _, ctrl := range s.controllers {
		routes := ctrl.Routes()
		initRoutes(s.router, routes)
	}

	fmt.Println("server listening at 8000")
	handler := cors.AllowAll().Handler(s.router)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), handler)
}
