package sql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofrs/uuid"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	// connectionstring = "user=postgres dbname=postgres sslmode=disable"
	connectionstring = "postgres://postgres:PG_PASS@localhost:5439/postgres?sslmode=disable"
)

type PersistResponse struct {
	isSuccess bool
}

//(id UUID, p1 VARCHAR(20), p2 VARCHAR(20), epoch INT, timestep INT)
func PersistParticleCollision(p1Name string, p2Name string, epoch int, timestep int) (PersistResponse, error) {
	log.Println("epoch: ", epoch, "timestep: ", timestep, "p1: ", p1Name, "p2: ", p2Name)
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return PersistResponse{isSuccess: false}, err
	}
	id, _ := uuid.NewV4()

	collString := fmt.Sprintf("'%v', '%v', '%v', %v, %v",
		id.String(),
		p1Name,
		p2Name,
		epoch,
		timestep)
	sql := "INSERT INTO particleCollision VALUES (" + collString + ")"
	log.Println("executing: ", sql)

	if resp, err := db.Exec(
		sql); err != nil {
		log.Fatal("Failed to persist particle collision", err)
	} else {
		log.Println("Persisted particle collision ", resp)
	}

	return PersistResponse{
		isSuccess: true,
	}, nil
}
func GetAllParticleCollisions() []string {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Println("error connecting to the database: ", err)
	}

	sql := "SELECT * FROM particleCollision"
	log.Println("executing: ", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal("Failed to get all particle collisions", err)
	}
	defer rows.Close()
	var collisions []string
	for rows.Next() {
		var (
			id       string
			p1       string
			p2       string
			epoch    int
			timestep int
		)
		if err := rows.Scan(&id, &p1, &p2, &epoch, &timestep); err != nil {
			log.Fatal("Failed to read particle collision row", err)
		}
		collision := fmt.Sprintf("id %v : p1 %v collided with p2 %v at epoch %v timestep %v",
			id,
			p1,
			p2,
			epoch,
			timestep)
		collisions = append(collisions, collision)
	}
	return collisions

}

func PersistWallCollisionEvent(pName string, wallName string, epoch int, timestep int) (PersistResponse, error) {
	log.Println("epoch: ", epoch, "timestep: ", timestep, "wallName: ", wallName, "pName: ", pName)
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return PersistResponse{isSuccess: false}, err
	}
	id, _ := uuid.NewV4()

	collString := fmt.Sprintf("'%v', '%v', '%v', %v, %v",
		id.String(),
		pName,
		wallName,
		epoch,
		timestep)
	sql := "INSERT INTO wallCollision VALUES (" + collString + ")"
	log.Println("executing: ", sql)

	if resp, err := db.Exec(
		sql); err != nil {
		log.Fatal("Failed to persist wall collision", err)
	} else {
		log.Println("Persisted wall collision ", resp)
	}

	return PersistResponse{
		isSuccess: true,
	}, nil
}

//(id UUID, particle VARCHAR(20), obj VARCHAR(20), epoch INT, timestep INT)`
func GetAllWallCollisionEvents() []string {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Println("error connecting to the database: ", err)
	}

	sql := "SELECT * FROM wallCollision"
	log.Println("executing: ", sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal("Failed to get all wall collisions", err)
	}
	defer rows.Close()
	var collisions []string
	for rows.Next() {
		var (
			id       string
			particle string
			objName  string
			epoch    int
			timestep int
		)
		if err := rows.Scan(&id, &particle, &objName, &epoch, &timestep); err != nil {
			log.Fatal("Failed to read wallCollision row", err)
		}
		collision := fmt.Sprintf("id %v : Particle %v collided with %v at epoch %v timestep %v",
			id,
			particle,
			objName,
			epoch,
			timestep)
		collisions = append(collisions, collision)
	}
	return collisions

}

func PersistParticleLocation(pName string, epoch int, timestep int, x int, y int) (PersistResponse, error) {
	return PersistResponse{}, nil
}

func createParticleCollision() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS particleCollision
	(id UUID, p1 VARCHAR(20), p2 VARCHAR(20), epoch INT, timestep INT)`
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
	sqlString := `CREATE TABLE IF NOT EXISTS wallCollision
	(id UUID, particle VARCHAR(20), obj VARCHAR(20), epoch INT, timestep INT)`
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute "+sqlString, err)
	} else {
		log.Println("Executed sql string "+sqlString, resp)
	}
}
func createParticle() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS particle
	(name VARCHAR(20), mass FLOAT)`
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
	sqlString := `CREATE TABLE IF NOT EXISTS location
	(particle VARCHAR(20), epoch INT, timestep INT, x INT, y INT)`
	if resp, err := db.Exec(sqlString); err != nil {
		log.Fatal("Failed to Execute"+sqlString, err)
	} else {
		log.Println("Executed sql string"+sqlString, resp)
	}
}

func init() {
	log.Println("initializing DB")
	createParticleCollision()
	createWallCollision()
	createParticle()
	createLocation()

}
