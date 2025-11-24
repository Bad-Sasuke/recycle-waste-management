package services

import (
	"crypto/rand"
	"encoding/base32"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/repositories"
)

type usersService struct {
	UsersRepository repositories.IUsersRepository
}

type IUsersService interface {
	GetAllUser() (*[]entities.UserDataFormat, error)
	InsertNewAccount(data *entities.NewUserBody) error
	InsertNewAccountWithRole(data *entities.NewUserBody, role entities.UserRole) error
	UpdateUser(userID string, data *entities.NewUserBody) error
	UpdateUserWithRole(userID string, data *entities.NewUserBody, role entities.UserRole) error
	UpdateUserImage(userID string, imageURL string) error
	DeleteUser(userID string) error
	GetUser(userID string) (*entities.UserDataFormat, error)
	UpdateUserRole(userID string, role entities.UserRole) error
}

func NewUsersService(repo0 repositories.IUsersRepository) IUsersService {
	return &usersService{
		UsersRepository: repo0,
	}
}

func (sv *usersService) InsertNewAccount(data *entities.NewUserBody) error {
	secret := make([]byte, 4)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return err
	}
	var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)
	userID := b32NoPadding.EncodeToString(secret)
	newData := &entities.UserDataFormat{
		UserID:   userID,
		Username: data.Username,
		Email:    data.Email,
		// Role will be set to default by repository layer
	}

	err = sv.UsersRepository.InsertNewUser(newData)
	if err != nil {
		return err
	}
	return nil
}

func (sv *usersService) InsertNewAccountWithRole(data *entities.NewUserBody, role entities.UserRole) error {
	secret := make([]byte, 4)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return err
	}
	var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)
	userID := b32NoPadding.EncodeToString(secret)
	newData := &entities.UserDataFormat{
		UserID:   userID,
		Username: data.Username,
		Email:    data.Email,
		Role:     string(role),
	}

	err = sv.UsersRepository.InsertNewUser(newData)
	if err != nil {
		return err
	}
	return nil
}

func (sv *usersService) GetAllUser() (*[]entities.UserDataFormat, error) {
	userData, err := sv.UsersRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return userData, nil

}

func (sv *usersService) UpdateUser(userID string, data *entities.NewUserBody) error {
	updateData := &entities.UserDataFormat{
		Username: data.Username,
		Email:    data.Email,
	}

	err := sv.UsersRepository.UpdateUser(userID, updateData)
	if err != nil {
		return err
	}
	return nil
}

func (sv *usersService) UpdateUserWithRole(userID string, data *entities.NewUserBody, role entities.UserRole) error {
	updateData := &entities.UserDataFormat{
		Username: data.Username,
		Email:    data.Email,
		Role:     string(role),
	}

	err := sv.UsersRepository.UpdateUser(userID, updateData)
	if err != nil {
		return err
	}
	return nil
}

func (sv *usersService) DeleteUser(userID string) error {
	if err := sv.UsersRepository.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}

func (sv *usersService) UpdateUserImage(userID string, imageURL string) error {
	updateData := &entities.UserDataFormat{
		ImageURL: imageURL,
	}

	err := sv.UsersRepository.UpdateUser(userID, updateData)
	if err != nil {
		return err
	}
	return nil
}

func (sv *usersService) GetUser(userID string) (*entities.UserDataFormat, error) {
	userData, err := sv.UsersRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (sv *usersService) UpdateUserRole(userID string, role entities.UserRole) error {
	userData, err := sv.UsersRepository.GetUser(userID)
	if err != nil {
		return err
	}

	updateData := &entities.UserDataFormat{
		UserID:    userData.UserID,
		Username:  userData.Username,
		Email:     userData.Email,
		ImageURL:  userData.ImageURL,
		Role:      string(role),
		CreatedAt: userData.CreatedAt,
		LastLogin: userData.LastLogin,
	}

	err = sv.UsersRepository.UpdateUser(userID, updateData)
	if err != nil {
		return err
	}
	return nil
}
