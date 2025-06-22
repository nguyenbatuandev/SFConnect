package service

import (
	"errors"
	"user_service/internal/entity"
	"user_service/internal/interface"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo _interface.UserRepository
	authRepo _interface.AuthSerice
}

func NewUserService(userRepo _interface.UserRepository, authRepo _interface.AuthSerice) *UserService {
	return &UserService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}
func (s *UserService) RegisterUser(User *entity.User) error {
	hashedPassword, err := s.authRepo.HashPassword(User.Password)
	if err != nil {
		return err
	}
	User.IsActive = true
	User.ID = uuid.New()
	User.Password = hashedPassword
	return s.userRepo.CreateUser(User)
}

func (s *UserService) LoginUser(user *entity.User) (*entity.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	isValid := s.authRepo.CheckPasswordHash(user.Password, existingUser.Password)
	if !isValid {
		return nil, errors.New("invalid email or password")
	}

	if !existingUser.IsActive {
		return nil, errors.New("user account is not active")
	}

	return existingUser, nil
}

func (s *UserService) GetProfile(id uuid.UUID) (*entity.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.DeleteUser(id)
}

func (s *UserService) UpdateProfile(id uuid.UUID, name string) (*entity.User, error) {
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		existingUser.Name = name
	}

	err = s.userRepo.UpdateUser(existingUser)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (s *UserService) ToggleUserLockByAdmin(id uuid.UUID) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	newStatus := !user.IsActive
	return s.userRepo.ToggleUserLockByAdmin(id, newStatus)
}

func (s *UserService) GetAllUsersByAdmin() ([]*entity.User, error) {
	return s.userRepo.GetAllUsersByAdmin()
}
