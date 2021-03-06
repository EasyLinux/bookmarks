package main

import (
	"flag"

	"github.com/nrocco/bookmarks/api"
	"github.com/nrocco/bookmarks/queue"
	"github.com/nrocco/bookmarks/scheduler"
	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Version stores the current version of Bend
	Version string

	// Workers stores the amount of workers that can do async tasks
	Debug = flag.Bool("debug", false, "Enable debug mode")

	// Workers stores the amount of workers that can do async tasks
	Workers = flag.Int("workers", 4, "The number of workers to start")

	// Interval controls how often feeds should be checked for new items
	Interval = flag.Int("interval", 30, "Fetch new feeds with this interval in minutes")

	// HTTPAddr stores the value for the --http option and defaults to 0.0.0.0:8000
	HTTPAddr = flag.String("http", "0.0.0.0:3000", "Address to listen for HTTP requests on")

	// Database holds the connection string for the database connection
	Database = flag.String("database", "data.sqlite", "The location to the sqlite database")
)

func main() {
	// Parse flags
	flag.Parse()

	log.Info().
		Bool("debug", *Debug).
		Int("workers", *Workers).
		Int("interval", *Interval).
		Str("http", *HTTPAddr).
		Str("database", *Database).
		Msg("Starting bookmarks")

	// Setup the global logger
	if *Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Setup the database
	store, err := storage.New(*Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open the database")
	}

	// Setup the async job queue
	queue := queue.New(store, *Workers)

	// Setup the http server
	api := api.New(store, queue)

	if *Interval != 0 {
		// Setup the periodic scheduler
		scheduler.New(store, queue, *Interval)
	} else {
		log.Info().Msg("Scheduler is not enabled")
	}

	// Run the http server
	if err := api.ListenAndServe(*HTTPAddr); err != nil {
		log.Warn().Err(err).Msg("Stopped the api server")
	}

	log.Info().Msg("Stopping bookmarks")
}
