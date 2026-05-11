// internal/service/user.go
package service

import (
	"backend/internal/model"
	repository "backend/internal/repo"
	"errors"
	"fmt"
)

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetMe(userID int) (*model.User, error) {
	return s.GetByID(userID)
}

func (s *userService) UpdateMe(userID int, req *model.UpdateUserRequest) (*model.User, error) {
	return s.Update(userID, req)
}

func (s *userService) GetByID(id int) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetByID: %w", err)
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

func (s *userService) GetAll() ([]model.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("UserService.GetAll: %w", err)
	}
	return users, nil
}

func (s *userService) SearchByName(name string) ([]model.User, error) {
	if len(name) < 2 {
		return nil, errors.New("имя для поиска должно быть минимум 2 символа")
	}
	users, err := s.repo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("UserService.SearchByName: %w", err)
	}
	return users, nil
}

func (s *userService) Update(id int, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("UserService.Update: %w", err)
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		existing, _ := s.repo.FindByEmail(*req.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email уже используется")
		}
		user.Email = *req.Email
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Position != nil {
		user.Position = req.Position
	}
	if req.Department != nil {
		user.Department = req.Department
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("UserService.Update: %w", err)
	}

	return user, nil
}

func (s *userService) Delete(id int) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("UserService.Delete: %w", err)
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}
	return s.repo.Delete(id)
}
