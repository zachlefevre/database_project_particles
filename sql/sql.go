package sql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	connectionstring = "postgres://root@localhost.com"
)

type PersistResponse struct {
	isSuccess bool
}

func PersistParticleCollision(p1Name string, p2Name string, epoch int, timestep int) (PersistResponse, error) {
	log.Println(p1Name + " hit " + p2Name)
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return PersistResponse{isSuccess: false}, err
	}
	sqlString, err := db.Prepare("INSERT INTO Knuth.particlecollision(p1, p2, epoch, timestep) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println("error inserting into particlecollision table: ", err)
		return PersistResponse{isSuccess: false}, err
	}

	sqlString.Exec(p1Name, p2Name, epoch, timestep)

	return PersistResponse{
		isSuccess: true,
	}, nil
}

func PersistWallCollisionEvent(pName string, wallName string, epoch int, timestep int) (PersistResponse, error) {
	log.Println("epoch: ", epoch, "timestep: ", timestep, "wallName: ", wallName, "pName: ", pName)
	return PersistResponse{}, nil
}

func PersistParticleLocation(pName string, epoch int, timestep int, x int, y int) (PersistResponse, error) {
	return PersistResponse{}, nil
}

func createDatabase() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := "CREATE DATABASE IF NOT EXISTS Knuth"
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute"+sqlString, err)
	} else {
		log.Println("Executed sql string"+sqlString, resp)
	}
}
func createParticleCollision() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS Knuth.particleCollision
	(id UUID, p1 STRING, p2 STRING, epoch INT, timestep INT)`
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute"+sqlString, err)
	} else {
		log.Println("Executed sql string"+sqlString, resp)
	}
}
func createWallCollision() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS Knuth.wallCollision
	(id UUID, particle STRING, obj STRING, epoch INT, timestep INT)`
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute"+sqlString, err)
	} else {
		log.Println("Executed sql string"+sqlString, resp)
	}
}
func createParticle() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS Knuth.particle
	(name STRING, mass FLOAT (3, 2))`
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute"+sqlString, err)
	} else {
		log.Println("Executed sql string"+sqlString, resp)
	}
}
func createLocation() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS Knuth.location
	(particle STRING, epoch INT, timestep INT, x INT, y INT)`
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute"+sqlString, err)
	} else {
		log.Println("Executed sql string"+sqlString, resp)
	}
}

func init() {
	// createDatabase()
	// createParticleCollision()
	// createWallCollision()
	// createParticle()
	// createLocation()

}
