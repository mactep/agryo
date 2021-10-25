// Package db provides access to a database and common operations
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mactep/agryo/hedera"
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

func (db DB) CreatePolygon(geometry string, hash string, properties hedera.Properties) error {
	result, err := db.conn.Exec("INSERT INTO polygons (geom) VALUES (ST_GeomFromGeoJSON($1))", geometry)
	if err != nil {
		return err
	}

	polygonId, err := result.LastInsertId()
	polygonId += 1
	// TODO: refactor this ASAP
	query := fmt.Sprintf(`
		INSERT INTO properties (
			gid, idfarmer, companyid, regionid, countryid, stateid,
			municipalityid, technicalid, status, activity, bsow,
			product, eharvest, latcenter, loncenter, hash, polygonId
		)
		VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', %d)
		`,
		properties.Gid, properties.Idfarmer, properties.Companyid, properties.Regionid, properties.Countryid, properties.Stateid,
		properties.Municipalityid, properties.Technicalid, properties.Status, properties.Activity, properties.Bsow,
		properties.Product, properties.Eharvest, properties.Latcenter, properties.Loncenter, hash, polygonId,
	)

	_, err = db.conn.Exec(query)

	return err
}

func (db DB) FindPolygonByID(id string) (string, error) {
	row := db.conn.QueryRow(`
		SELECT json_build_object(
			'type',       'Feature',
			'geometry',   ST_AsGeoJSON(polygons.geom)::json,
			'properties', properties
		)
		FROM polygons
		join properties on polygons.id=properties.polygonid
		where polygons.id = $1
	`, id)

	var feature string
	err := row.Scan(&feature)
	return feature, err
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS polygons (
			id SERIAL PRIMARY KEY,
			geom geometry(POLYGON, 4326) NOT NULL
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
			hash VARCHAR(255) UNIQUE,
			CONSTRAINT fk_polygon FOREIGN KEY(polygonId) REFERENCES polygons(id)
		);
	`)

	return err
}
