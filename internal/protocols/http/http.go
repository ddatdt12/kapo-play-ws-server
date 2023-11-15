package http

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http/ws"
	"github.com/rs/zerolog/log"
)

type HttpImpl struct {
	httpServer *http.Server
	wsHub      *ws.Hub
}

func NewHttpProtocol(wsHub *ws.Hub) *HttpImpl {
	return &HttpImpl{
		wsHub: wsHub,
	}
}

func (p *HttpImpl) setupRouter() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		p.wsHub.ServeWs(w, r)
	})
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("serveHome")
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func (p *HttpImpl) Listen() {
	var addr = flag.String("addr", ":8080", "http service address")
	log.Info().Msgf("Server started on Port %s ", *addr)

	go p.wsHub.Run()
	p.setupRouter()

	p.httpServer = &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 5 * time.Second,
	}
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
