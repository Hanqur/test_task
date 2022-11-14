package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"path/filepath"
	"strconv"

	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Item struct {
	Product       string
	Price, Rating int
}

type Config struct {
	LogLevel    string   `envconfig:"log_level"`
	DbFileNames []string `envconfig:"db_file_names"`
}

func main() {
	var conf Config

	if err := envconfig.Process("app", &conf); err != nil {
		log.Panic().Err(err).Msg("failed to process config")
	}

	logLevel, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Panic().Err(err).Msg("failed to set log level")
	}

	log := zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()

	log.Info().Msg("start program")

	for _, fileName := range conf.DbFileNames {
		ext := filepath.Ext(fileName)

		switch ext {
		case ".json":
			processJson(fileName, log)
		case ".csv":
			processCsv(fileName, log)
		}
	}

	log.Info().Msg("end program")
}

func processJson(fileName string, log zerolog.Logger) {
	log = log.With().Str("file", "json").Logger()

	log.Info().Msg("start process file")

	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Error().Err(err).Msg("fail to open file")

		return
	}

	defer jsonFile.Close()

	dec := json.NewDecoder(jsonFile)

	_, err = dec.Token()
	if err != nil {
		log.Error().Err(err).Msg("fail to read open bracket")

		return
	}

	var maxPrice Item
	var maxRating Item

	for dec.More() {
		var item Item

		err := dec.Decode(&item)
		if err != nil {
			log.Error().Err(err).Msg("fail to decode line")

			return
		}

		if item.Price > maxPrice.Price {
			maxPrice = item
		}

		if item.Rating > maxRating.Rating {
			maxRating = item
		}
	}

	_, err = dec.Token()
	if err != nil {
		log.Error().Err(err).Msg("fail to read open bracket")

		return
	}

	log.Info().Interface("item", maxPrice).Msg("highest priced product")
	log.Info().Interface("item", maxRating).Msg("top rated product")

	log.Info().Msg("end process file")
}

func processCsv(fileName string, log zerolog.Logger) {
	log = log.With().Str("file", fileName).Logger()

	log.Info().Msg("start process file")

	file, err := os.Open(fileName)
	if err != nil {
		log.Error().Err(err).Msg("fail to open file")

		return
	}

	parser := csv.NewReader(file)

	var maxPrice Item
	var maxRating Item

	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("fail to read line")

			continue
		}

		if len(record) != 3 {
			continue
		}

		if record[1] == "Price" {
			continue
		}

		price, err := strconv.Atoi(record[1])
		if err != nil {
			log.Error().Err(err).Msg("fail to convert price of product")
		}

		rating, err := strconv.Atoi(record[2])
		if err != nil {
			log.Error().Err(err).Msg("fail to convert rating of product")
		}

		item := Item{
			Product: record[0],
			Price:   price,
			Rating:  rating,
		}

		if item.Price > maxPrice.Price {
			maxPrice = item
		}

		if item.Rating > maxRating.Rating {
			maxRating = item
		}
	}

	log.Info().Interface("item", maxPrice).Msg("highest priced product")
	log.Info().Interface("item", maxRating).Msg("top rated product")

	log.Info().Msg("end process file")
}
