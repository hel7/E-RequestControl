package services

import (
	"mime/multipart"
	Request_Manager "request_manager_api"
	"request_manager_api/pkg/repository"
)

type AdministratorService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdministratorService {
	return &AdministratorService{repo: repo}
}
func (s *AdministratorService) GetUserByID(userID int) (Request_Manager.User, error) {
	return s.repo.GetUserByID(userID)
}
func (s *AdministratorService) CreateUser(user Request_Manager.User) (int, error) {
	if err := user.ValidatePassword(); err != nil {
		return 0, err
	}
	if err := user.ValidateEmail(); err != nil {
		return 0, err
	}
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}
func (s *AdministratorService) UpdateUser(UserID int, input Request_Manager.UpdateUserInput, user Request_Manager.User) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if err := user.ValidateEmail(); err != nil {
		return err
	}
	if err := user.ValidatePassword(); err != nil {
		return err
	}
	if input.Password != nil {
		user.Password = *input.Password
		if err := user.ValidatePassword(); err != nil {
			return err
		}
		hashed := generatePasswordHash(*input.Password)
		input.Password = &hashed
	}
	return s.repo.UpdateUser(UserID, input, user)
}

func (s *AdministratorService) GetAllUsers() ([]Request_Manager.User, error) {
	return s.repo.GetAllUsers()
}
func (s *AdministratorService) Delete(UserID int) error {
	return s.repo.Delete(UserID)
}
func (s *AdministratorService) GetFilteredTickets(filter Request_Manager.TicketFilter) ([]Request_Manager.Ticket, error) {
	return s.repo.GetFilteredTickets(filter)
}
func (s *AdministratorService) BackupData(backupPath string) error {
	err := s.repo.BackupData(backupPath)
	if err != nil {
		return err
	}

	return nil
}
func (s *AdministratorService) RestoreData(backupPath multipart.File) error {
	err := s.repo.RestoreData(backupPath)
	if err != nil {
		return err
	}

	return nil
}
func (s *AdministratorService) ExportData(exportPath string) error {
	err := s.repo.ExportData(exportPath)
	if err != nil {
		return err
	}

	return nil
}
func (s *AdministratorService) ImportData(importPath string) error {
	err := s.repo.ImportData(importPath)
	if err != nil {
		return err
	}

	return nil
}
