package sql

type PersistResponse struct {
	isSucces bool
}

func PersistParticleCollision(p1Name string, p2Name string, epoch int, timestep int) (PersistResponse, error) {
	return PersistResponse{}, nil
}

func PersistWallCollisionEvent(pName string, wallName string, epoch int, timestep int) (PersistResponse, error) {
	return PersistResponse{}, nil
}

func PersistParticleLocation(pName string, epoch int, timestep int, x int, y int) (PersistResponse, error) {
	return PersistResponse{}, nil
}
