package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http/ws"
	"github.com/ddatdt12/kapo-play-ws-server/src/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type HttpImpl struct {
	httpServer  *http.Server
	wsHub       *ws.Hub
	httpHandler *handlers.HttpHandlerImpl
}

func NewHttpProtocol(wsHub *ws.Hub, httpHandler *handlers.HttpHandlerImpl) *HttpImpl {
	return &HttpImpl{
		wsHub:       wsHub,
		httpHandler: httpHandler,
	}
}

func (p *HttpImpl) setupRouter(router *mux.Router) {

	router.Use(corsMiddleware)
	p.httpHandler.Router(router)

	router.HandleFunc("/ws/host", func(w http.ResponseWriter, r *http.Request) {
		p.wsHub.ServeHostWs(w, r)
	})
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		p.wsHub.ServeClientWs(w, r)
	})
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (p *HttpImpl) Listen() {
	addr := fmt.Sprintf(":%v", configs.EnvConfigs.SERVER_PORT)

	router := mux.NewRouter()

	go p.wsHub.Run()
	p.setupRouter(router)

	p.httpServer = &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}

	log.Info().Msgf("Server started on Port %s ", addr)
	err := p.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

func (p *HttpImpl) Shutdown(ctx context.Context) error {
	if err := p.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
