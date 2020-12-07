package controllers

import "github.com/diptyojha/goLngFirstProject/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Locations routes
	s.Router.HandleFunc("/locations", middlewares.SetMiddlewareJSON(s.CreateLocation)).Methods("POST")
	s.Router.HandleFunc("/locations", middlewares.SetMiddlewareJSON(s.GetLocations)).Methods("GET")
	s.Router.HandleFunc("/locations/{id}", middlewares.SetMiddlewareJSON(s.GetLocation)).Methods("GET")
	s.Router.HandleFunc("/locations/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateLocation))).Methods("PUT")
	s.Router.HandleFunc("/locations/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteLocation)).Methods("DELETE")

	//Locations routes
	s.Router.HandleFunc("/userlocations", middlewares.SetMiddlewareJSON(s.CreateUserLocation)).Methods("POST")
	s.Router.HandleFunc("/userlocations", middlewares.SetMiddlewareJSON(s.GetUserLocation)).Methods("GET")
	s.Router.HandleFunc("/userlocations/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUserLocation))).Methods("GET")
	s.Router.HandleFunc("/userlocations/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUserLocation))).Methods("PUT")
	s.Router.HandleFunc("/userlocations/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUserLocation)).Methods("DELETE")
}
