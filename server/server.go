// Local server capabilities
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mactep/agryo/db"
	"github.com/mactep/agryo/hedera"
)

type Server struct {
	hedera hedera.HederaAPI
	db     db.DB
}

// Returns a new instance of the local server or an error if it fails to connect to the database
func NewServer(accountID, privateKey, user, password, host, port, DBName string) (Server, error) {
	hederaClient, err := hedera.NewHederaAPI(accountID, privateKey)
	if err != nil {
		return Server{}, err
	}

	db, err := db.NewDB(user, password, host, port, DBName)
	if err != nil {
		return Server{}, err
	}

	return Server{
		hedera: hederaClient,
		db:     db,
	}, nil
}

// Attach the handlers and run the server on port 8080
func (server Server) Run() {
	http.HandleFunc("/hash-polygon", server.handleHashPolygon)
	http.HandleFunc("/polygon", server.handlePolygon)

	port := ":8080"
	fmt.Printf("Running server on port %s\n", port)
	http.ListenAndServe(port, nil)
}

// Handler that parses the quote request, forward it to the API, parses the
// response and returns it
func (server Server) handleHashPolygon(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()

	var collection hedera.FeatureCollection
	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// TODO: make it async
	for _, polygon := range collection.Features {
		ch := make(chan []byte)
		defer close(ch)
		go server.hedera.HashPolygon(polygon, ch)
		hash := <-ch

		geometry, err := json.Marshal(polygon.Geometry)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		err = server.db.CreatePolygon(string(geometry), hash, polygon.Properties)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Searches in the database for for a polygon with the given ID and returns
// it alongside the hash and it's properties
func (server Server) handlePolygon(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	defer r.Body.Close()

	id := r.URL.Query().Get("id")
	if id == "" {
		fmt.Fprint(w, "Invalid ID")
		return
	}

	polygon, err := server.db.FindPolygonByID(id)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, polygon)
}
