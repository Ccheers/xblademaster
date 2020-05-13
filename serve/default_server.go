package serve

import (
	"github.com/Ccheers/xblademaster"
	"github.com/Ccheers/xblademaster/middleware"
)

// DefaultServer returns an Engine instance with the Recovery and Logger middleware already attached.
func DefaultServer(conf *xblademaster.ServerConfig) *xblademaster.Engine {
	engine := xblademaster.NewServer(conf)
	engine.Use(middleware.Recovery(), middleware.Trace(), middleware.Logger())
	return engine
}
