package main

import (
	"os"

	"github.com/rs/cors"

	"github.com/zachlefevre/project_knuth/sql"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3080"
	}
	handler := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "GET"},
			AllowCredentials: true,
			// Debug:            true,
		}).Handler(initRoutes())
	log.Println("Http Server Listening...")
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/collision/particle", particleCollision).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/collision/wall", wallCollision).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/location/particle", particleLocation).Methods("POST")
	router.HandleFunc("/api/entity/particle", particle).Methods("POST")

	router.HandleFunc("/api/collision/wall", getWallCollisions).Methods("GET")
	router.HandleFunc("/api/collision/particle", getParticleCollisions).Methods("GET")
	router.HandleFunc("/api/entity/particle", getParticles).Methods("GET")
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
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}

type particleType struct {
	Name string  `json:"name"`
	Mass float64 `json:"mass"`
}

func particle(w http.ResponseWriter, r *http.Request) {
	var particle particleType
	err := json.NewDecoder(r.Body).Decode(&particle)
	if err != nil {
		http.Error(w, "Invalid Wall Collision Event", 500)
		return
	}

	resp, err := sql.PersistParticle(particle.Name, particle.Mass)
	if err != nil {
		http.Error(w, "Failed to persist particle", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func getParticles(w http.ResponseWriter, r *http.Request) {
	collisions := sql.GetAllParticles()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(collisions)
	w.Write(j)
}
func getWallCollisions(w http.ResponseWriter, r *http.Request) {
	collisions := sql.GetAllWallCollisionEvents()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(collisions)
	w.Write(j)
}
func getParticleCollisions(w http.ResponseWriter, r *http.Request) {
	collisions := sql.GetAllParticleCollisions()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(collisions)
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
