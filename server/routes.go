package main

import "github.com/gorilla/mux"

// I might split the two login routes into its own service.
func connectRoutes(r *mux.Router, s *Server) {
	r.HandleFunc("/login", s.handlBnetLogin)
	r.HandleFunc("/bnet_oauth_cb", s.handleBnetCallback)

	r.HandleFunc("/player", s.jwtMiddleware(s.handlerPlayer)).Methods("GET") //player
	r.HandleFunc("/players", s.jwtMiddleware(s.handlerPlayers)).Methods("GET")

	r.HandleFunc("/rankings", s.jwtMiddleware(s.listRanking)).Methods("GET")
	r.HandleFunc("/rankings/{id}", s.jwtMiddleware(s.getRanking)).Methods("GET")
	r.HandleFunc("/rankings", s.jwtMiddleware(s.createRanking)).Methods("POST")
	r.HandleFunc("/rankings/promote", s.jwtMiddleware(s.promoteRanking)).Methods("POST")
	// r.HandleFunc("rankings/{id}", s.jwtMiddleware(s.updateRanking)).Methods("PUT")
	r.HandleFunc("/rankings/{id}", s.jwtMiddleware(s.deleteRanking)).Methods("DELETE")
}

// r.HandleFunc("/players", nil)
// r.HandleFunc("/players/:id", nil)

// r.HandleFunc("/rankings").Method("GET")
// r.HandleFunc("/rankings").Method("POST")
// r.HandleFunc("/rankings/:id").Method("GET")
// r.HandleFunc("/rankings/:id").Method("PUT")

// r.HandleFunc("/challenges")

// r.HandleFunc("/check")
