// Package db provides access to a database and common operations
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

// Returns a new instance of the database or a error if something goes wrong
// while trying to connect to the specified database
func NewDB(user, password, host, port, DBName string) (DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		host, user, password, DBName, port,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{}, err
	}

	err = initDB(db)
	return DB{
		conn: db,
	}, err
}

func (db DB) CreatePolygon(geometry string, hash string) error {
	_, err := db.conn.Exec("INSERT INTO polygons (geom, hash) VALUES (ST_GeomFromGeoJSON($1), $2)", geometry, hash)
	return err
}

func (db DB) FindPolygonByID(id string) (string, error) {
	row := db.conn.QueryRow(`
		SELECT ST_AsGeoJSON(geom)::json
		FROM polygons WHERE id=$1;
	`, id)

	var geoJSON string
	err := row.Scan(&geoJSON)
	return geoJSON, err
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS polygons (
			id SERIAL PRIMARY KEY,
			geom geometry(POLYGON, 4326) NOT NULL,
			hash VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS properties (
			id SERIAL PRIMARY KEY,
			polygonId INT NOT NULL,
			gid INT UNIQUE,
			idfarmer VARCHAR(255),
			companyid VARCHAR(255),
			regionid VARCHAR(255),
			countryid VARCHAR(255),
			stateid VARCHAR(255),
			municipalityid VARCHAR(255),
			technicalid VARCHAR(255),
			status VARCHAR(255),
			activity VARCHAR(255),
			bsow INT,
			product VARCHAR(255),
			eharvest VARCHAR(255),
			latcenter VARCHAR(255),
			loncenter VARCHAR(255),
			CONSTRAINT fk_polygon FOREIGN KEY(polygonId) REFERENCES polygons(id)
		);
	`)

	return err
}
