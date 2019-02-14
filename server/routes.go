package main

import "github.com/gorilla/mux"

// I might split the two login routes into its own service.
func connectRoutes(r *mux.Router, s *Server) {
	r.HandleFunc("/login", s.handlBnetLogin)
	r.HandleFunc("/bnet_oauth_cb", s.handleBnetCallback)

	r.HandleFunc("/player", s.jwtMiddleware(s.handlerPlayer)) //player
	r.HandleFunc("/players", s.jwtMiddleware(s.handlerPlayers))
	//players
	//players/{id}
}

// r.HandleFunc("/players", nil)
// r.HandleFunc("/players/:id", nil)

// r.HandleFunc("/rankings").Method("GET")
// r.HandleFunc("/rankings").Method("POST")
// r.HandleFunc("/rankings/:id").Method("GET")
// r.HandleFunc("/rankings/:id").Method("PUT")

// r.HandleFunc("/challenges")

// r.HandleFunc("/check")
