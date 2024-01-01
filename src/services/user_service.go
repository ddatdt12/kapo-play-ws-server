package services

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameExist = errors.New("username exist")
)

type IUserService interface {
	UsernameExist(ctx context.Context, code string, username string) (bool, error)
	JoinGame(ctx context.Context, code string, user *models.User) error
	QuitGame(ctx context.Context, code string, username string) error
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

func (s *UserService) JoinGame(ctx context.Context, code string, userDto *models.User) error {
	_, err := s.gameRepo.GetByCode(ctx, code)

	if err != nil {
		return errors.Wrap(err, "join game")
	}

	exist, err := s.userRepo.UsernameExist(ctx, code, userDto.Username)

	if err != nil {
		return errors.Wrap(err, "join game")
	} else if exist {
		return ErrUsernameExist
	}

	user := models.User{
		Username: userDto.Username,
	}

	err = s.userRepo.Add(ctx, code, &user)

	if err != nil {
		return errors.Wrap(err, "join game")
	}

	return nil
}

func (s *UserService) QuitGame(ctx context.Context, code string, username string) error {
	err := s.userRepo.Remove(ctx, code, username)

	if err != nil {
		return errors.Wrap(err, "quit game")
	}

	return nil
}
