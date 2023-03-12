package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"alox.sh"
)

// type contextKey struct {
// 	name string
// }

// var stateContextKey = &contextKey{"state"}

// if request.Context().Value(stateContextKey) == nil {
// 	request = request.WithContext(context.WithValue(request.Context(), stateContextKey, server))
// }

const serverContextKey = "config"

type Config struct {
	DBPath string
	DB     *sql.DB
}

func SetupConfig() (config *Config, err error) {
	config = &Config{
		DBPath: "data.db",
	}

	// _ = os.Remove(config.DBPath)

	if config.DB, err = sql.Open("sqlite3", config.DBPath); err != nil {
		log.Fatal(err)
	}

	if _, err = config.DB.Exec(`
CREATE TABLE todos (
	id      INT NOT NULL PRIMARY KEY,
	title   CHAR(128),
	content CHAR(512)
);
`); err != nil {
		return
	}

	return
}

func FromContext(context context.Context) (*Config, error) {
	switch value := context.Value(serverContextKey).(type) {
	case *Config:
		return value, nil
	}

	return nil, fmt.Errorf("*Config not found within context")
}

func FromServer(server alox.Server) (*Config, error) {
	return FromContext(server.Context())
}

func (config *Config) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, serverContextKey, config)
}

func (config *Config) Dispose() {
	_ = config.DB.Close()
}
