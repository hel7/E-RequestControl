package services

import (
	Request_Manager "request_manager_api"
	"request_manager_api/pkg/repository"
)

type TicketService struct {
	repo repository.Ticket
}

func NewTicketService(repo repository.Ticket) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) CreateTicket(userID int, ticket Request_Manager.Ticket) (int, error) {
	return s.repo.CreateTicket(userID, ticket)
}

func (s *TicketService) GetUserTickets(ticketID int) ([]Request_Manager.Ticket, error) {
	return s.repo.GetUserTickets(ticketID)
}

func (s *TicketService) DeleteUserTicket(ticketID, userID int) error {
	return s.repo.DeleteUserTicket(ticketID, userID)
}
func (s *TicketService) UpdateTicket(userID int, ticketID int, input Request_Manager.UpdateTicketInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateTicket(userID, ticketID, input)
}

func (s *TicketService) DeleteTicketAdmin(ticketID int) error {
	return s.repo.DeleteTicketAdmin(ticketID)
}
func (s *TicketService) GetAllTickets() ([]Request_Manager.Ticket, error) {
	return s.repo.GetAllTickets()
}
func (s *TicketService) GetTicketByID(ticketID int) (Request_Manager.Ticket, error) {
	return s.repo.GetTicketByID(ticketID)
}
