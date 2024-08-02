package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	Server ConfServer
}

type ConfServer struct {
	Port         int           `env:"PORT,required"`
	ReadTimeout  time.Duration `env:"TIMEOUT_READ,default=5s"`
	WriteTimeout time.Duration `env:"TIMEOUT_WRITE,default=10s"`
	IdleTimeout  time.Duration `env:"TIMEOUT_IDLE,default=15s"`
	Debug        bool          `env:"DEBUG,default=false"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
