package main

import (
	"os"

	"github.com/zachlefevre/project_knuth/sql"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
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
	router.HandleFunc("/api/sig", sig).Methods("GET")
	return router
}
func sig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Made by zachlefevre@gmail.com</h1>"))
}

type particleCollisionEvent struct {
	P1       string `json:"p1Name"`
	P2       string `json:"p2Name"`
	Epoch    int    `json:"epoch"`
	Timestep int    `json:"timestep"`
}

func particleCollision(w http.ResponseWriter, r *http.Request) {
	var event particleCollisionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid Particle Collision Event", 500)
		return
	}

	resp, err := sql.PersistParticleCollision(event.P1, event.P2, event.Epoch, event.Timestep)
	if err != nil {
		http.Error(w, "Failed to persist particle collision", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	log.Println(w.Header())
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}

type wallCollisionEvent struct {
	P        string `json:"p"`
	Wall     string `json:"wall"`
	Epoch    int    `json:"epoch"`
	Timestep int    `json:"timestep"`
}

func wallCollision(w http.ResponseWriter, r *http.Request) {
	var event wallCollisionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid Wall Collision Event", 500)
		return
	}

	resp, err := sql.PersistWallCollisionEvent(event.P, event.Wall, event.Epoch, event.Timestep)
	if err != nil {
		http.Error(w, "Failed to persist wall collision", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}

type particleLocationObj struct {
	ParticleName string `json:"particleName"`
	Epoch        int    `json:"epoch"`
	Timestep     int    `json:"timestep"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

func particleLocation(w http.ResponseWriter, r *http.Request) {
	var loc particleLocationObj
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil {
		http.Error(w, "Invalid Wall Collision Event", 500)
		return
	}
	resp, err := sql.PersistParticleLocation(loc.ParticleName, loc.Epoch, loc.Timestep, loc.X, loc.Y)
	if err != nil {
		http.Error(w, "Failed to persist wall collision", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
