package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Komplementariteten/lutra"

	"github.com/Komplementariteten/lutra/auth"
	"github.com/Komplementariteten/lutra/db"
	"github.com/Komplementariteten/lutra/util"
)

// HTTPServer to handle Stuff
type HTTPServer struct {
	s         *http.Server
	a         *authHandle
	p         *portalHandle
	i         *apiHandle
	IsStarted bool
}

// Start starts a HTTPServer
func (h *HTTPServer) Start() {
	log.Fatal(h.s.ListenAndServe())
}

func (h *HTTPServer) Shutdown() {
	h.s.Close()
}

// CreateNew creates a default HttpServer
func CreateNew(port int, staticFilePath string, config *lutra.LutraConfig) (*HTTPServer, error) {
	db, err := db.NewConnection(context.Background(), config.MongoDbConnectionStr, config.MongoDbName)
	if err != nil {
		panic(err)
	}
	s, err := util.NewMailClient(config.SmtpServer, config.SmtpPort, config.SmtpLogin, config.SmtpPassword)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	a := &authHandle{db: db, mail: s}
	a.sessions = make(map[string]*Session)
	p := &portalHandle{db: db, mail: s}
	i := &apiHandle{db: db}
	mux.Handle("/auth/", http.StripPrefix("/auth/", a))
	mux.Handle("/portal/", http.StripPrefix("/portal/", p))
	mux.Handle("/api/", http.StripPrefix("/api/", i))
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(staticFilePath))))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	handle := &HTTPServer{s: server, a: a, p: p}
	auth.Environment.Config = config
	return handle, nil
}
