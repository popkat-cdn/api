package main

import (
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"Popkat/api"
	"Popkat/constants"
	"Popkat/routes/auth"
	"Popkat/routes/tests"
	"Popkat/state"
	"Popkat/types"

	"github.com/cloudflare/tableflip"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
	"go.uber.org/zap"

	"github.com/infinitybotlist/eureka/zapchi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "embed"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//go:embed docs/docs.html
var docsHTML string

var openapi []byte

// Simple middleware to handle CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// limit body to 10mb
		r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")

		if r.Method == "OPTIONS" {
			w.Write([]byte{})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func main() {
	state.Setup()

	docs.DocsSetupData = &docs.SetupData{
		URL:         state.Config.Sites.API.Parse(),
		ErrorStruct: types.ApiError{},
		Info: docs.Info{
			Title:          "Popkat CDN",
			TermsOfService: "https://popkat.select-list.xyz/legal#tos",
			Version:        "0.1",
			Description:    "Advanced CDN for Select List",
			Contact: docs.Contact{
				Name:  "Select List",
				URL:   "https://select-list.xyz",
				Email: "contact@select-list.xyz",
			},
			License: docs.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
		},
	}

	docs.Setup()

	docs.AddSecuritySchema("User", "User-Auth", "Requires a user token. Should be prefixed with `User ` in `Authorization` header.")
	docs.AddSecuritySchema("Service", "Service-Auth", "Requires a service token. Should be prefixed with `Service ` in `Authorization` header.")

	api.Setup()

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.CleanPath,
		corsMiddleware,
		zapchi.Logger(state.Logger, "api"),
		middleware.Timeout(30*time.Second),
	)

	routers := []uapi.APIRouter{
		// Use same order as routes folder
		auth.Router{},
		tests.Router{},
	}

	for _, router := range routers {
		name, desc := router.Tag()
		if name != "" {
			docs.AddTag(name, desc)
			uapi.State.SetCurrentTag(name)
		} else {
			panic("Router tag name cannot be empty")
		}

		router.Routes(r)
	}

	r.Get("/openapi", func(w http.ResponseWriter, r *http.Request) {
		w.Write(openapi)
	})

	docsTempl := template.Must(template.New("docs").Parse(docsHTML))

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/popkat", http.StatusFound)
	})

	r.Get("/docs/{srv}", func(w http.ResponseWriter, r *http.Request) {
		var docMap = map[string]string{
			"popkat": "/openapi",
		}

		srv := chi.URLParam(r, "srv")

		if srv == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid service name"))
			return
		}

		v, ok := docMap[srv]

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid service"))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		docsTempl.Execute(w, map[string]string{
			"url": v,
		})
	})

	// Load openapi here to avoid large marshalling in every request
	var err error
	openapi, err = json.Marshal(docs.GetSchema())

	if err != nil {
		panic(err)
	}

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(constants.EndpointNotFound))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(constants.MethodNotAllowed))
	})

	// If GOOS is windows, do normal http server
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		upg, _ := tableflip.New(tableflip.Options{})
		defer upg.Stop()

		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP)
			for range sig {
				state.Logger.Info("Received SIGHUP, upgrading server")
				upg.Upgrade()
			}
		}()

		// Listen must be called before Ready
		ln, err := upg.Listen("tcp", state.Config.Meta.Port.Parse())

		if err != nil {
			state.Logger.Fatal("Error binding to socket", zap.Error(err))
		}

		defer ln.Close()

		server := http.Server{
			ReadTimeout: 30 * time.Second,
			Handler:     r,
		}

		go func() {
			err := server.Serve(ln)
			if err != http.ErrServerClosed {
				state.Logger.Error("Server failed due to unexpected error", zap.Error(err))
			}
		}()

		if err := upg.Ready(); err != nil {
			state.Logger.Fatal("Error calling upg.Ready", zap.Error(err))
		}

		<-upg.Exit()
	} else {
		// Tableflip not supported
		state.Logger.Warn("Tableflip not supported on this platform, this is not a production-capable server.")
		err = http.ListenAndServe(state.Config.Meta.Port.Parse(), r)

		if err != nil {
			state.Logger.Fatal("Error binding to socket", zap.Error(err))
		}
	}
}
