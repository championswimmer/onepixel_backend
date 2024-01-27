package controllers

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
	"onepixel_backend/src/db"
)

type EventsController struct {
	// event logging db (not the main app db)
	db *gorm.DB
}

func CreateEventsController() *EventsController {
	eventsDB := lo.Must(db.GetEventsDB())
	return &EventsController{
		db: eventsDB,
	}
}
