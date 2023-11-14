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
}

func NewHttpProtocol() *HttpImpl {
	return &HttpImpl{}
}

func (p *HttpImpl) setupRouter(hub *ws.Hub) {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r)
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

	hub := ws.NewHub()
	go hub.Run()
	p.setupRouter(hub)

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
