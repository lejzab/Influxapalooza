package main

import (
	"flag"
	"lejzab/influxapalooza/configuration"
	"lejzab/influxapalooza/db"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// configures logger
func set_logger(debug *bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func main() {

	debug := flag.Bool("debug", false, "sets log level to debug")
	// data_to_write_filename := flag.String("filename", "data.csv", "name of file to read data")
	config_filename := flag.String("config", "config.toml", "name of config file")
	dbname := flag.String("dbname", "testdb", "name of influx database")
	createdb := flag.Bool("createdb", false, "creates influx database")
	flag.Parse()

	set_logger(debug)

	log.Info().Msg("No cześć")
	defer log.Info().Msg("papa, biedaczyska")

	config, err := configuration.NewConfiguration(*config_filename)
	if err != nil {
		log.Error().Str("filename", *config_filename).Msgf("error reading config: %s", err)
		return
	}
	log.Debug().Interface("combined config", config).Msg("read config")
	influx_client := db.NewInflux(config.Influx, *dbname)
	if *createdb {
		influx_client.CreateDB(*dbname)
		if err != nil {
			log.Error().Err(err).Msg("error writing to influx")
		}
	}
	// write 10 times in loop
	for i := 0; i < 1; i++ {
		influx_client.WriteTestData("test_measurement")
	}
}
