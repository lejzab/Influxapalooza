package db

import (
	"fmt"
	"lejzab/influxapalooza/configuration"
	"math/rand"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/rs/zerolog/log"
)

// influxdb conf and funcs
type Influx struct {
	// Host     string
	// Port     int
	User      string
	Pass      string
	Addr      string
	Timeout   time.Duration
	Dbname    string
	Precision string
}

// NewInflux create a new influxdb client and tries to ping database.
func NewInflux(c configuration.Influx, dbname string) *Influx {
	addr := fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
	return &Influx{
		// Host:    c.Host,
		// Port:    c.Port,
		User:      c.Username,
		Pass:      c.Password,
		Addr:      addr,
		Timeout:   time.Second * time.Duration(c.Timeout),
		Precision: c.Precision,
		Dbname:    dbname,
	}
}

// Create influx http client
func newClient(i *Influx) (client.HTTPClient, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     i.Addr,
		Username: i.User,
		Password: i.Pass,
	})
	if err != nil {
		return nil, err
	}
	t, r, err := c.Ping(i.Timeout)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("influxdb ping success, %v, %s", t, r)
	return c, nil
}

// Write point to influxdb
func (i *Influx) Write(measurement string, tags map[string]string, fields map[string]interface{}) error {
	// create a new client
	c, err := newClient(i)
	if err != nil {
		return err
	}
	defer c.Close()
	// create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.Dbname,
		Precision: i.Precision,
	})
	if err != nil {
		return err
	}
	// create a point and add to batch
	pt, err := client.NewPoint(measurement, tags, fields)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	// write the batch
	if err := c.Write(bp); err != nil {
		return err
	}
	return nil
}

// create db in  influx, takes a dbname as a parameter
func (i *Influx) CreateDB(dbname string) error {
	// create a new client
	c, err := newClient(i)
	defer c.Close()
	if err != nil {
		return err
	}
	// create query to create a database
	q := client.Query{
		Command: fmt.Sprintf("CREATE DATABASE %s", dbname),
	}
	_, err = c.Query(q)
	if err != nil {
		return err
	}
	log.Debug().Msgf("influxdb db %s created", dbname)
	replication_name := "three_days"
	replicas := 1
	duration := time.Hour * 24 * 3
	q = client.Query{
		Command: fmt.Sprintf("CREATE RETENTION POLICY %s ON %s DURATION %s REPLICATION %d DEFAULT", replication_name, dbname, duration, replicas),
	}
	_, err = c.Query(q)
	if err != nil {
		return err
	}

	return nil
}

// writes test data to influx, with Write function.
func (i *Influx) WriteTestData(measurement string) error {
	tags := map[string]string{"host": "server01", "region": "us-west"}
	fields := map[string]interface{}{
		"percent": rand.Intn(100),
		"cpu":     rand.Intn(200),
	}
	// write the data
	err := i.Write(measurement, tags, fields)
	if err != nil {
		if err != nil {
			log.Error().Err(err).Msg("error writing test data to influx")
		}
		return err
	}
	log.Debug().Str("measurement", measurement).
		Str("tags", fmt.Sprintf("%v", tags)).
		Str("fields", fmt.Sprintf("%v", fields)).
		Msgf("influxdb test data written")
	return nil
}
