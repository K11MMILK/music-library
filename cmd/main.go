package main

import (
	timetracker "time-tracker"
	"time-tracker/pkg/handler"
	"time-tracker/pkg/repository"
	"time-tracker/pkg/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Music-Library
// @version 1.0
// @description API Server for Music-library Application

// @host localhost:8000
// @BasePath
const longText = `[Intro]
Is this the real life? Is this just fantasy?
Caught in a landslide, no escape from reality
Open your eyes, look up to the skies and see
I'm just a poor boy, I need no sympathy
Because I'm easy come, easy go, little high, little low
Any way the wind blows doesn't really matter to me, to me

[Verse 1]
Mama, just killed a man
Put a gun against his head, pulled my trigger, now he's dead
Mama, life had just begun
But now I've gone and thrown it all away
Mama, ooh, didn't mean to make you cry
If I'm not back again this time tomorrow
Carry on, carry on as if nothing really matters

[Verse 2]
Too late, my time has come
Sends shivers down my spine, body's aching all the time
Goodbye, everybody, I've got to go
Gotta leave you all behind and face the truth
Mama, ooh (Any way the wind blows)
I don't wanna die
I sometimes wish I'd never been born at all

[Verse 3]
I see a little silhouetto of a man
Scaramouche, Scaramouche, will you do the Fandango?
Thunderbolt and lightning, very, very frightening me
(Galileo) Galileo, (Galileo) Galileo, Galileo Figaro magnifico
But I'm just a poor boy, nobody loves me
He's just a poor boy from a poor family
Spare him his life from this monstrosity
Easy come, easy go, will you let me go?
Bismillah, no, we will not let you go
(Let him go) Bismillah, we will not let you go
(Let him go) Bismillah, we will not let you go
(Let me go) Will not let you go
(Let me go) Will not let you go
(Never, never, never, never let me go) Ah
No, no, no, no, no, no, no
(Oh, mamma mia, mamma mia) Mamma mia, let me go
Beelzebub has a devil put aside for me, for me, for me

[Bridge]
So you think you can stone me and spit in my eye?
So you think you can love me and leave me to die?
Oh, baby, can't do this to me, baby
Just gotta get out, just gotta get right outta here

[Outro]
(Ooh)
(Ooh, yeah, ooh, yeah)
Nothing really matters, anyone can see
Nothing really matters
Nothing really matters to me
Any way the wind blows`

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(logrus.StandardLogger().Out)

	logger.Info("Starting the application")

	if err := godotenv.Load("configs/config.env"); err != nil {
		logger.WithError(err).Fatal("Error loading .env file")
	}
	logger.Info(".env file loaded")

	if err := initConfig(); err != nil {
		logger.WithError(err).Fatal("Error occurred while initializing config")
	}

	m, err := migrate.New(
		"file://migrations",
		"postgres://"+
			viper.GetString("DB_USERNAME")+":"+
			viper.GetString("DB_PASSWORD")+"@"+
			viper.GetString("DB_HOST")+":"+
			viper.GetString("DB_PORT")+"/"+
			viper.GetString("DB_DBNAME")+"?sslmode="+
			viper.GetString("DB_SSLMODE"),
	)

	if err != nil {
		logger.WithError(err).Fatal("Error occurred while creating migration instance")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.WithError(err).Fatal("Error occurred while migrating database")
	}
	logger.Info("Database migration successful")

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_DBNAME"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
	})
	if err != nil {
		logger.WithError(err).Fatal("Error occurred while connecting to database")
	}
	logger.Info("Database connection established")

	logger.Info("Populating the database with test data")

	isEmpty, err := isDatabaseEmpty(db)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to check if the database is empty")
	}

	if isEmpty {
		logrus.Info("The database is empty, populating with test data")

		_, err = db.Exec(`
			INSERT INTO groupss (groupName) VALUES
			('Metallica'),
			('Nirvana'),
			('Queen')
			ON CONFLICT DO NOTHING;
		`)
		if err != nil {
			logger.WithError(err).Fatal("Failed to populate the database with test data in 'groupss'")
		}

		_, err = db.Exec(`
			INSERT INTO songs (songName, groupId) VALUES
			('Enter Sandman', 1),
			('Smells Like Teen Spirit', 2),
			('Bohemian Rhapsody', 3)
			ON CONFLICT DO NOTHING;
		`)
		if err != nil {
			logger.WithError(err).Fatal("Failed to populate the database with test data in 'songs'")
		}

		_, err = db.Exec(`
			INSERT INTO songDetails (releaseDate, text, link, songId) VALUES
			('1991-07-29', 'Say your prayers, little one...', 'https://example.com/enter-sandman', 1),
			('1991-09-10', 'Load up on guns, bring your friends...', 'https://example.com/smells-like-teen-spirit', 2),
			('1975-10-31', $1, 'https://example.com/bohemian-rhapsody', 3)
			ON CONFLICT DO NOTHING;
		`, longText)
		if err != nil {
			logger.WithError(err).Fatal("Failed to populate the database with test data in 'songDetails'")
		}
	} else {
		logrus.Info("The database is not empty")
	}

	if err != nil {
		logger.WithError(err).Fatal("Failed to populate the database with test data")
	}
	logger.Info("Test data inserted successfully")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	logger.Info("Repositories and services initialized")

	srv := new(timetracker.Server)
	port := viper.GetString("port")
	logger.Infof("Starting server on port %s", port)
	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logger.WithError(err).Fatal("Error occurred while running HTTP server")
	}
	logger.Info("Server started successfully")

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logger.WithError(err).Fatal("Error occurred while migrating database")
	}
	logger.Info("Database migration down")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	return viper.ReadInConfig()
}

func isDatabaseEmpty(db *sqlx.DB) (bool, error) {
	var exists bool
	query := `SELECT NOT EXISTS (SELECT 1 FROM songs)`
	err := db.Get(&exists, query)
	if err != nil {
		return false, err
	}
	return exists, nil
}
