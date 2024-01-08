package handlers

import (
	"encoding/json"
	"io"
	syslog "log"
	"net/http"
	"os"

	"github.com/ddatdt12/kapo-play-ws-server/src/middlewares"
	"github.com/ddatdt12/kapo-play-ws-server/src/services"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type HttpHandlerImpl struct {
	gameService services.IGameService
}

func NewHttpHandler(gameService services.IGameService) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		gameService: gameService,
	}
}

func (h *HttpHandlerImpl) Router(router *mux.Router) {
	logger := syslog.New(os.Stdout, "", syslog.LstdFlags)
	logMiddleware := middlewares.NewLogMiddleware(logger)
	router.Use(logMiddleware.Func())
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/games/{code}", h.GetGame).Methods(http.MethodGet)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Info().Msgf("serveHome")
	io.WriteString(w, "Hello HOME KAPO")
}

func (h *HttpHandlerImpl) GetGame(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	game, err := h.gameService.GetGame(r.Context(), code)
	if err != nil {
		if errors.Is(err, services.ErrGameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RespondJSON(w, game, "success", http.StatusOK)
}

type JSONResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondJSON(w http.ResponseWriter, data interface{}, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponse{
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}
