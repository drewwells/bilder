package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/handlers"
)

type syncFile struct {
	sync.Mutex
	os.File
	f *os.File
}

func (sf *syncFile) Write(p []byte) (n int, err error) {
	sf.Lock()
	defer sf.Unlock()
	return sf.f.Write(p)
}

type server struct {
	http.Server
	sync.RWMutex
	albumUpdates <-chan []album
	dir          string
	accessLog    string
	logFile      *syncFile
	albums       map[string]http.Handler
}

func newServer(d, al string, au <-chan []album) *server {
	return &server{dir: d, accessLog: al, albumUpdates: au}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		http.Error(w, "404 page not found", 404)
		return
	}

	sep := strings.Index(r.URL.Path, "/")
	an := r.URL.Path
	if sep > 0 {
		an = r.URL.Path[:sep]
	}

	s.RLock()
	h, ok := s.albums[an]
	s.RUnlock()

	if !ok {
		http.Error(w, "404 page not found", 404)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, an)
	r.URL.Path = p
	h.ServeHTTP(w, r)
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path[len("/a/"):]
	a, ok := assets[p]
	if !ok {
		log.Printf("Got asset request %v but not available.\n", r.URL.Path)
		http.Error(w, "404 page not found", 404)
		return
	}

	w.Header().Set("Content-Type", a.ContentType)
	w.Write(a.Content)
}

func (s *server) listenForUpdates() {
	for as := range s.albumUpdates {
		hs := make(map[string]http.Handler)
		for _, a := range as {
			hs[a.name] = http.FileServer(http.Dir(filepath.Join(s.dir, a.name)))
			if s.logFile != nil {
				hs[a.name] = handlers.CombinedLoggingHandler(s.logFile, hs[a.name])
			}
		}

		s.Lock()
		s.albums = hs
		s.Unlock()
	}
}

func (s *server) serve() {
	go s.listenForUpdates()

	if s.accessLog != "" {
		lf, err := os.OpenFile(s.accessLog, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("Cannot open access log %#v with write access, err=%v", s.accessLog, err)
		} else {
			s.logFile = &syncFile{f: lf}
			defer lf.Close()
		}
	}

	mux := http.NewServeMux()
	if len(assets) == 0 {
		log.Printf("Serving assets for assets directory.")
		mux.Handle("/a/", http.StripPrefix("/a/", http.FileServer(http.Dir("assets"))))
	} else {
		log.Printf("Serving assets from memory.")
		mux.HandleFunc("/a/", assetsHandler)
	}

	mux.Handle("/b/", http.StripPrefix("/b/", s))
	log.Printf("serving %#v\n", s.dir)

	s.Addr = ":8173"
	s.Handler = mux

	log.Printf("Serving on http://0.0.0.0:8173")
	log.Fatal(s.ListenAndServe())
}
