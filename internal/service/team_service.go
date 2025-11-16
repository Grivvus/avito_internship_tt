package service

import (
	"context"
	"errors"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/repository"
)

var TeamAlreadyExistError error = errors.New("already exists")

type TeamService struct {
	teamRepo *repository.TeamRepository
	userRepo *repository.UserRepository
}

func NewTeamService(
	teamRepo *repository.TeamRepository,
	userRepo *repository.UserRepository,
) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (ts *TeamService) AddTeam(ctx context.Context, team api.Team) error {
	err := ts.teamRepo.Create(ctx, team.TeamName)
	if err != nil {
		// better to add additional error-matching
		return TeamAlreadyExistError
	}
	for _, user := range team.Members {
		model := api.User{
			IsActive: user.IsActive,
			UserId:   user.UserId,
			Username: user.Username,
			TeamName: team.TeamName,
		}
		err := ts.userRepo.CreateOrUpdate(ctx, model)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TeamService) GetTeam(
	ctx context.Context,
	params api.GetTeamGetParams,
) (api.Team, error) {
	var team api.Team

	_, err := ts.teamRepo.Get(ctx, params.TeamName)

	if err != nil {
		return team, ResourceNotFoundError
	}

	team.TeamName = params.TeamName

	users, err := ts.userRepo.FindByTeam(ctx, params.TeamName)
	if err != nil {
		return team, err
	}
	for _, user := range users {
		team.Members = append(team.Members, api.TeamMember{
			IsActive: user.IsActive,
			UserId:   user.UserId,
			Username: user.Username,
		})
	}

	return team, nil
}
