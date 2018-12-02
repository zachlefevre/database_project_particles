package main

import (
	"bytes"
	"fmt"
	"strings"

	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	router.HandleFunc("/sig", sig).Methods("GET")
	return router
}
func sig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Made by zachlefevre@gmail.com</h1>"))
}

type collisionEvent struct {
	p1       particle.Particle
	p2       particle.Particle
	epoch    int
	timestep int
}

func particleCollision(w http.ResponseWriter, r *http.Request) {
	var event collisionEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid Collision Event", 500)
		return
	}

	sql.persistParticleCollision(event.p1, event.p2, event.epoch, event.timestep)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func createAlgorithmRPC(cmd *pb.CreateAlgorithmCommand) (*pb.Algorithm, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlgorithmAggregateClient(conn)
	return client.CreateAlgorithm(context.Background(), cmd)
}
func getAlgorithm(w http.ResponseWriter, r *http.Request) {
	var algo pb.Algorithm
	err := json.NewDecoder(r.Body).Decode(&algo)
	if err != nil {
		http.Error(w, "Invalid algoritm", 500)
		return
	}

	queryID, _ := uuid.NewV4()
	getQuery := pb.GetAlgorithmQuery{
		Algorithm: &algo,
		CreatedOn: time.Now().Unix(),
		Id:        queryID.String(),
	}
	resp, err := getAlgorithmRPC(&getQuery)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get algorithm", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if resp == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusFound)
	}
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func getAlgorithmRPC(query *pb.GetAlgorithmQuery) (*pb.Algorithm, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlgorithmAggregateClient(conn)
	return client.GetAlgorithm(context.Background(), query)
}

func addFile(w http.ResponseWriter, r *http.Request) {
	algoName := r.URL.Query().Get("algorithm")
	version := r.URL.Query().Get("version")
	if algoName == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	if version == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileContent := make([]byte, 100)
	_, err = file.Read(fileContent)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to read file", 500)
		return
	}
	fileContent = bytes.Trim(fileContent, "\x00")
	fileStr := string(fileContent)

	cd := strings.Split(handler.Header.Get("Content-Disposition"), "; ")
	nameSection := cd[2]
	log.Println(cd, nameSection)
	name := strings.Split(nameSection, "=")[1]
	name = name[1 : len(name)-1]
	log.Println(name)
	queryID, _ := uuid.NewV4()
	getQuery := pb.GetAlgorithmQuery{
		Algorithm: &pb.Algorithm{
			Name:    algoName,
			Version: version,
		},
		CreatedOn: time.Now().Unix(),
		Id:        queryID.String(),
	}
	algo, err := getAlgorithmRPC(&getQuery)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get algorithm", 500)
		return
	}
	log.Println(handler.Header)
	algoFile := pb.AlgorithmFile{
		Content:  fileStr,
		Name:     name,
		Filetype: handler.Header.Get("Content-Type"),
	}

	addFileCmd := pb.AssociateFileCommand{
		Algorithm:     algo,
		AlgorithmFile: &algoFile,
	}
	resp, err := addFileRPC(&addFileCmd)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func addFileRPC(cmd *pb.AssociateFileCommand) (*pb.Algorithm, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlgorithmAggregateClient(conn)
	return client.AssociateFile(context.Background(), cmd)
}
