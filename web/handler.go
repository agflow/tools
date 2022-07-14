package web

import (
	"github.com/agflow/tools/agtime"
	"github.com/agflow/tools/cache"
	"github.com/agflow/tools/sql/db"
)

// Handler is a struct that packs all the necessary services for a web handler
type Handler struct {
	T        agtime.ClockTime
	DBSvc    *db.Service
	CacheSvc *cache.Service
}
