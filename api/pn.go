package main

import (
	"github.com/zachlefevre/project_knuth/sql"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zachlefevre/project_knuth/particle"
)

func main() {
	server := &http.Server{
		Addr:    ":3030",
		Handler: initRoutes(),
	}
	log.Println("Http Server Listening...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/collision/particle", particleCollision).Methods("POST")
	router.HandleFunc("/api/collision/wall", wallCollision).Methods("POST")
	router.HandleFunc("/api/location/particle", particleLocation).Methods("POST")
	router.HandleFunc("/sig", sig).Methods("GET")
	return router
}
func sig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Made by zachlefevre@gmail.com</h1>"))
}

type particleCollisionEvent struct {
	p1       particle.Particle
	p2       particle.Particle
	epoch    int
	timestep int
}

func particleCollision(w http.ResponseWriter, r *http.Request) {
	var event particleCollisionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid Particle Collision Event", 500)
		return
	}

	resp, err := sql.PersistParticleCollision(event.p1.Name, event.p2.Name, event.epoch, event.timestep)
	if err != nil {
		http.Error(w, "Failed to persist particle collision", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}

type wallCollisionEvent struct {
	p        particle.Particle
	wall     string
	epoch    int
	timestep int
}

func wallCollision(w http.ResponseWriter, r *http.Request) {
	var event wallCollisionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid Wall Collision Event", 500)
		return
	}

	resp, err := sql.PersistWallCollisionEvent(event.p.Name, event.wall, event.epoch, event.timestep)
	if err != nil {
		http.Error(w, "Failed to persist wall collision", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}

type particleLocationObj struct {
	particleName string
	epoch        int
	timestep     int
	x            int
	y            int
}

func particleLocation(w http.ResponseWriter, r *http.Request) {
	var loc particleLocationObj
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil {
		http.Error(w, "Invalid Wall Collision Event", 500)
		return
	}
	resp, err := sql.PersistParticleLocation(loc.particleName, loc.epoch, loc.timestep, loc.x, loc.y)
	if err != nil {
		http.Error(w, "Failed to persist wall collision", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
