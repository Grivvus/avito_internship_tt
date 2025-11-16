package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type PullRequestService struct {
	prRepo   *repository.PullRequestRepository
	userRepo *repository.UserRepository
}

func NewPullRequestService(
	prRepo *repository.PullRequestRepository,
	userRepo *repository.UserRepository,
) *PullRequestService {
	return &PullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
	}
}

func (pr *PullRequestService) Create(
	ctx context.Context, prCreate api.PostPullRequestCreateJSONBody,
) (api.PullRequest, error) {
	var now time.Time
	response := api.PullRequest{
		AuthorId:        prCreate.AuthorId,
		PullRequestId:   prCreate.PullRequestId,
		PullRequestName: prCreate.PullRequestName,
		Status:          api.PullRequestStatusOPEN,
		CreatedAt:       &now,
	}
	err := pr.prRepo.Create(ctx, api.PullRequestShort{
		AuthorId:        prCreate.AuthorId,
		PullRequestId:   prCreate.PullRequestId,
		PullRequestName: prCreate.PullRequestName,
		Status:          api.PullRequestShortStatusOPEN,
	})

	if err != nil {
		var errToReturn error = fmt.Errorf("Unkown server error: %w", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				errToReturn = ResourceNotFoundError
			case "23505":
				errToReturn = PRAlreadyExistError
			}
		}
		return response, errToReturn
	}

	potentialReviewers, err := pr.userRepo.FindOtherMembers(ctx, prCreate.AuthorId)
	if err != nil {
		return response, fmt.Errorf("On FindOtherMembers: %w", err)
	}

	filteredReviewers := pr.filterReviewers(ctx, potentialReviewers, []string{prCreate.AuthorId})
	if len(filteredReviewers) > 2 {
		filteredReviewers = filteredReviewers[:2]
	}

	err = pr.prRepo.AssignReviewers(ctx, prCreate.PullRequestId, filteredReviewers)
	if err != nil {
		return response, fmt.Errorf("On AssignReviewers: %w", err)
	}

	for _, reviewer := range filteredReviewers {
		response.AssignedReviewers = append(response.AssignedReviewers, reviewer)
	}

	now = time.Now()

	return response, nil

}

func (pr *PullRequestService) Merge(
	ctx context.Context,
	prToMerge api.PostPullRequestMergeJSONBody,
) (api.PullRequest, error) {
	prToChange, err := pr.prRepo.Get(ctx, prToMerge.PullRequestId)
	if err != nil {
		return prToChange, ResourceNotFoundError
	}

	// already merged, do nothing
	if prToChange.Status == api.PullRequestStatusMERGED {
		return prToChange, nil
	}

	prToChange.Status = api.PullRequestStatusMERGED
	return pr.prRepo.Merge(ctx, prToChange)
}

func (pr *PullRequestService) Reassign(
	ctx context.Context,
	prToReassign api.PostPullRequestReassignJSONBody,
) (api.PullRequest, error) {
	log.Println(prToReassign)
	prToChange, err := pr.prRepo.Get(ctx, prToReassign.PullRequestId)
	if err != nil {
		return prToChange, ResourceNotFoundError
	}

	if prToChange.Status == api.PullRequestStatusMERGED {
		return prToChange, CantReassignOnMergedPRError
	}

	log.Println(prToChange.AssignedReviewers)
	log.Println(prToReassign.OldUserId)
	if !slices.Contains(prToChange.AssignedReviewers, prToReassign.OldUserId) {
		return prToChange, ReviewerNotAssignedError
	}

	potentialReviewers, err := pr.userRepo.FindOtherMembers(ctx, prToChange.AuthorId)
	if err != nil {
		return prToChange, fmt.Errorf("On FindOtherMembers: %w", err)
	}

	filter := append([]string{}, prToChange.AssignedReviewers...)
	filter = append(filter, prToChange.AuthorId)
	filteredReviewers := pr.filterReviewers(ctx, potentialReviewers, filter)
	if len(filteredReviewers) == 0 {
		return prToChange, NoCandidatesError
	}
	chosen := filteredReviewers[0]

	err = pr.prRepo.ReassignReviewer(ctx, prToChange.PullRequestId, prToReassign.OldUserId, chosen)
	if err != nil {
	}

	prToChange.AssignedReviewers = slices.DeleteFunc(
		prToChange.AssignedReviewers,
		func(reviewerId string) bool {
			return reviewerId == prToReassign.OldUserId
		},
	)
	prToChange.AssignedReviewers = append(prToChange.AssignedReviewers, chosen)

	return prToChange, nil
}

func (pr *PullRequestService) filterReviewers(
	ctx context.Context,
	potentialReviewers []api.User,
	inappropriateIds []string,
) []string {
	validReviewers := make([]string, 0, 2)
	for _, reviewer := range potentialReviewers {
		if reviewer.IsActive && !slices.Contains(inappropriateIds, reviewer.UserId) {
			validReviewers = append(validReviewers, reviewer.UserId)
		}
	}
	return validReviewers
}
