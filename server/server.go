// Local server capabilities
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mactep/agryo/hedera"
)

type Server struct {
	hedera hedera.HederaAPI
}

// Returns a new instance of the local server or an error if it fails to connect to the database
func NewServer(accountID, privateKey, user, password, host, port, DBName string) (Server, error) {
	hederaClient, err := hedera.NewHederaAPI(accountID, privateKey)
	if err != nil {
		return Server{}, err
	}

	return Server{
		hedera: hederaClient,
	}, nil
}

// Attach the handlers and run the server on port 8080
func (server Server) Run() {
	http.HandleFunc("/polygon", server.handlePolygon)

	http.ListenAndServe(":8080", nil)
}

// Handler that parses the quote request, forward it to the API, parses the
// response and returns it
func (server Server) handlePolygon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()

	var polygon hedera.Polygon
	err := json.NewDecoder(r.Body).Decode(&polygon)
	if err != nil {
		fmt.Fprint(w, err)
	}

	ch := make(chan []byte)
	defer close(ch)
	go server.hedera.HashPolygon(polygon, ch)
	hash := <- ch

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hash)
}
