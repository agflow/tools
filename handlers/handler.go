package handler

import (
	"github.com/go-redis/redis/v8"

	"github.com/agflow/tools/agtime"
	"github.com/agflow/tools/sql"
)

type Handlers struct {
	T        agtime.ClockTime
	DBSvc    *sql.DBService
	CacheSvc *redis.Client
}
