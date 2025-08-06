package contract

import (
	"context"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/utils"
)

type EventFilter struct {
	Category  string
	Status    string
	StartDate *string
	EndDate   *string
	Search    string
	Page      int
	Limit     int
}

type EventRepository interface {
	Create(ctx context.Context, event *entities.Event) *exception.AppError
	GetByID(ctx context.Context, id utils.BinaryUUID) (*entities.Event, *exception.AppError)
	GetAll(ctx context.Context, filter EventFilter) ([]*entities.Event, int64, *exception.AppError)
	Update(ctx context.Context, event *entities.Event) *exception.AppError
	Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	CheckNameExists(ctx context.Context, name string, excludeID *utils.BinaryUUID) (bool, *exception.AppError)
	GetAvailableCapacity(ctx context.Context, eventID utils.BinaryUUID) (int, *exception.AppError)
}

type EventService interface {
	CreateEvent(ctx context.Context, event *entities.Event, logCtx entities.ActivityLogContext) (*entities.Event, *exception.AppError)
	GetEventByID(ctx context.Context, id utils.BinaryUUID) (*entities.Event, *exception.AppError)
	GetAllEvents(ctx context.Context, filter EventFilter, logCtx entities.ActivityLogContext) ([]*entities.Event, int64, *exception.AppError)
	UpdateEvent(ctx context.Context, event *entities.Event, logCtx entities.ActivityLogContext) (*entities.Event, *exception.AppError)
	DeleteEvent(ctx context.Context, id utils.BinaryUUID, logCtx entities.ActivityLogContext) *exception.AppError
}