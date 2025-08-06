package contract

import (
	"context"
	"root-app/internal/api/dto"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/utils"
)

type TicketFilter struct {
	EventID utils.BinaryUUID
	UserID  utils.BinaryUUID
	Status  string
	Page    int
	Limit   int
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *entities.Ticket) *exception.AppError
	GetByID(ctx context.Context, id utils.BinaryUUID) (*entities.Ticket, *exception.AppError)
	GetAll(ctx context.Context, filter TicketFilter) ([]*entities.Ticket, int64, *exception.AppError)
	Update(ctx context.Context, ticket *entities.Ticket) *exception.AppError
	GetByUser(ctx context.Context, userID utils.BinaryUUID, filter TicketFilter) ([]*entities.Ticket, int64, *exception.AppError)
	Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	CountTicketsByEvent(ctx context.Context, eventID utils.BinaryUUID) (int, *exception.AppError)
	GetTicketsByEvent(ctx context.Context, eventID utils.BinaryUUID) ([]*entities.Ticket, *exception.AppError)
	GenerateUniqueTicketCode(ctx context.Context) (string, *exception.AppError)
}

type TicketService interface {
	PurchaseTicket(ctx context.Context, ticket *entities.Ticket, logCtx entities.ActivityLogContext) (*entities.Ticket, *exception.AppError)
	GetTicketByID(ctx context.Context, id utils.BinaryUUID) (*entities.Ticket, *exception.AppError)
	GetUserTickets(ctx context.Context, userID utils.BinaryUUID, filter TicketFilter, logCtx entities.ActivityLogContext) ([]*entities.Ticket, int64, *exception.AppError)
	GetAllTickets(ctx context.Context, filter TicketFilter, logCtx entities.ActivityLogContext) ([]*entities.Ticket, int64, *exception.AppError)
	UpdateTicketStatus(ctx context.Context, id utils.BinaryUUID, req *dto.UpdateTicketStatusRequest, reason string, logCtx entities.ActivityLogContext) (*entities.Ticket, *exception.AppError)
	DeleteTicket(ctx context.Context, id utils.BinaryUUID, logCtx entities.ActivityLogContext) *exception.AppError
}
