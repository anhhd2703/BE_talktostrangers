package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SignalServer struct {
	config signalConf
	Router *http.ServeMux
}

type signalConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func NewServer(config signalConf) *SignalServer {
	return &SignalServer{
		config: config,
	}
}

func (s *SignalServer) Start() {

	// use gorilla mux to middle ware
	//Use the default DefaultServeMux.
	mux := http.NewServeMux()
	myHttp := &http.Server{
		Addr:           s.config.Host + ":" + strconv.Itoa(s.config.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	myHttp.ListenAndServe()
	s.Router = mux
}

func (s *SignalServer) initializeRoutes() {
	// verify sdk
	s.Router.Handle("/", CreateNodeHandler())
}
func CreateNodeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// close resp.Body()
		sendJSON(w, "data", nil)
	})
}
func sendJSON(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	result := map[string]interface{}{"success": true}
	if err != nil {
		log.Println("[sendJSON] err: %v", err)
		result["success"] = false
		result["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else if data != nil {
		result["data"] = data
	}
	e := json.NewEncoder(w).Encode(result)
	if e != nil {
		log.Println("[sendJSON][json.NewEncoder] err: %v", e)
	}
}
