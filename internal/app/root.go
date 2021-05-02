package app

import (
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
)

type Root struct {
	Server     ServerConfig    `mapstructure:"server"`
	Log        log.Config      `mapstructure:"log"`
	MiddleWare mid.LogConfig   `mapstructure:"middleware"`
	Firestore  FirestoreConfig `mapstructure:"firestore"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type FirestoreConfig struct {
	File      string `mapstructure:"file"`
	ProjectId string `mapstructure:"project_id"`
}
