package services

import (
	"context"
	"errors"

	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
)

var (
	ErrUserNotFound  = errors.New("game not found")
	ErrUsernameExist = errors.New("username exist")
)

type IUserService interface {
	UsernameExist(ctx context.Context, code string, username string) (bool, error)
	JoinGame(ctx context.Context, code string, user dto.UserDto) (*models.User, error)
}

type UserService struct {
	userRepo repositories.IUserRepository
	gameRepo repositories.IGameRepository
}

func NewUserService(
	userRepo repositories.IUserRepository,
	gameRepo repositories.IGameRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		gameRepo: gameRepo,
	}
}

func (s *UserService) UsernameExist(ctx context.Context, code string, username string) (bool, error) {
	return s.userRepo.UsernameExist(ctx, code, username)
}

func (s *UserService) JoinGame(ctx context.Context, code string, userDto dto.UserDto) (*models.User, error) {
	_, err := s.gameRepo.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	exist, err := s.userRepo.UsernameExist(ctx, code, userDto.Username)

	if err != nil {
		return nil, err
	} else if exist {
		return nil, ErrUsernameExist
	}

	user := models.User{
		Username: userDto.Username,
	}

	_, err = s.userRepo.AddUser(ctx, code, user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
