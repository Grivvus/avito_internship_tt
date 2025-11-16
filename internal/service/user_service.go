package service

import (
	"context"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/models"
	"github.com/Grivvus/reviewers/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
	prRepo   *repository.PullRequestRepository
}

func NewUserservice(
	userRepo *repository.UserRepository,
	prRepo *repository.PullRequestRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (us *UserService) SetIsActive(
	ctx context.Context,
	model api.PostUsersSetIsActiveJSONBody,
) (api.User, error) {
	user, err := us.userRepo.Get(ctx, model.UserId)
	if err != nil {
		return user, ResourceNotFoundError
	}

	user.IsActive = model.IsActive

	err = us.userRepo.Update(ctx, user)

	return user, err
}

func (us *UserService) UserReviews(
	ctx context.Context,
	model api.GetUsersGetReviewParams,
) (models.UserReviewsResponse, error) {
	var response models.UserReviewsResponse
	response.UserId = model.UserId

	pullRequests, err := us.prRepo.FindByReviewer(ctx, model.UserId)
	if err != nil {
		return response, err
	}
	for _, pr := range pullRequests {
		prModel := api.PullRequestShort{
			AuthorId:        pr.AuthorId,
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			Status:          api.PullRequestShortStatus(pr.Status),
		}
		response.PullRequests = append(response.PullRequests, prModel)
	}

	return response, nil
}
