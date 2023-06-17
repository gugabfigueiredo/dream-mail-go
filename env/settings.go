package env

import (
	"encoding/json"
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/service"
	"time"
)

type settings struct {
	// Server
	Server struct {
		Port              string        `default:"8080"`
		Context           string        `default:"sw-api"`
		UpdateRefsTimeout time.Duration `default:"4h"`
	}

	// Log
	Log *log.Config

	// Service
	Service service.Config
}

var Settings settings

func (s settings) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
