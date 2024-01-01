package tests

import (
	"github.com/samber/lo"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
)

var App = server.CreateApp(lo.Must(db.GetDB()))
