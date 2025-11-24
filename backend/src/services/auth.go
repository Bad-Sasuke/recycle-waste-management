package services

import (
	"fmt"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/infrastructure/httpclients"
	"recycle-waste-management-backend/src/infrastructure/providers"
	"recycle-waste-management-backend/src/middlewares"
	"recycle-waste-management-backend/src/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type authService struct {
	authGithub httpclients.IAuthGithub
	UserRepo   repositories.IUsersRepository
	Firebase   providers.IFirebaseProvider
}

type IAuthService interface {
	AuthGithub(code string) (entities.TokenResponseGithub, error)
	AuthGoogle(uid string) (entities.TokenResponseGoogle, error)
}

func NewAuthService(userRepo repositories.IUsersRepository) IAuthService {
	return &authService{
		UserRepo:   userRepo,
		authGithub: httpclients.NewAuthGithub(),
		Firebase:   providers.NewFirebaseProvider(),
	}
}

func (s *authService) AuthGithub(code string) (entities.TokenResponseGithub, error) {
	data, err := s.authGithub.GitHubAccessToken(code)
	if err != nil {
		return data, err
	}
	user, err := s.authGithub.GetUserGithub(data.AccessToken)
	if err != nil {
		return data, err
	}

	var userData entities.UserDataFormat

	token, err := middlewares.GenerateJWTToken(user.UserID)
	if err != nil {
		return data, err
	}
	data.JWT = *token.Token

	go func() {

		userOld, err := s.UserRepo.GetUser(user.UserID)
		if err == nil && mongo.ErrNoDocuments != err {
			fmt.Println("User already exists")
		}

		if userOld == nil || user.UserID == "" {
			if user.UserID != "" {
				userData = entities.UserDataFormat{
					UserID:    user.UserID,
					Username:  user.Username,
					Email:     user.Email,
					ImageURL:  user.ImageURL,
					JWT:       *token.Token,
					CreatedAt: time.Now().UTC().Add(7 * time.Hour),
					LastLogin: time.Now().UTC().Add(7 * time.Hour),
				}
			}

			if err := s.UserRepo.InsertNewUser(&userData); err != nil {
				fmt.Println(err)
			}
		} else {
			userOld.LastLogin = time.Now().UTC().Add(7 * time.Hour)
			userOld.JWT = *token.Token
			if err := s.UserRepo.UpdateUser(user.UserID, userOld); err != nil {
				fmt.Println(err)
			}
		}

	}()
	return data, nil
}

func (s *authService) AuthGoogle(uid string) (entities.TokenResponseGoogle, error) {
	var userNew entities.UserDataFormat
	userData, err := s.Firebase.GetUserFirebaseAuth(uid)
	if err != nil {
		return entities.TokenResponseGoogle{}, err
	}

	token, err := middlewares.GenerateJWTToken(userData.UID)
	if err != nil {
		return entities.TokenResponseGoogle{}, err
	}
	go func() {

		userOld, err := s.UserRepo.GetUser(userData.UID)
		if err == nil && mongo.ErrNoDocuments != err {
			fmt.Println("User already exists")
		}
		username := strings.Split(userData.Email, "@")[0]
		if userOld == nil || userData.UID == "" {
			if userData.UID != "" {
				userNew = entities.UserDataFormat{
					UserID:    userData.UID,
					Username:  username,
					Email:     userData.Email,
					ImageURL:  userData.PhotoURL,
					JWT:       *token.Token,
					CreatedAt: time.Now().UTC().Add(7 * time.Hour),
					LastLogin: time.Now().UTC().Add(7 * time.Hour),
				}
			}

			if err := s.UserRepo.InsertNewUser(&userNew); err != nil {
				fmt.Println(err)
			}
		} else {
			userOld.LastLogin = time.Now().UTC().Add(7 * time.Hour)
			userOld.JWT = *token.Token
			if err := s.UserRepo.UpdateUser(userData.UID, userOld); err != nil {
				fmt.Println(err)
			}
		}

	}()
	return entities.TokenResponseGoogle{JWT: *token.Token}, nil
}
