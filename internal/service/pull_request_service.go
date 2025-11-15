package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/repository"
)

var PRAlreadyExistError = errors.New("already exist")
var CantReassignOnMergedPRError = errors.New("cannot reassign on merged PR")
var ReviewerNotAssignedError = errors.New("reviewer is not assigned to this PR")
var NoCandidatesError = errors.New("no active replacement candidate in team")

type PullReqeustService struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
}

func NewPullRequestService(
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
) PullReqeustService {
	return PullReqeustService{
		prRepo: prRepo,
	}
}

func (pr PullReqeustService) Create(
	ctx context.Context, prCreate api.PostPullRequestCreateJSONBody,
) (api.PullRequest, error) {
	var response api.PullRequest
	err := pr.prRepo.Create(ctx, api.PullRequestShort{
		AuthorId:        prCreate.AuthorId,
		PullRequestId:   prCreate.AuthorId,
		PullRequestName: prCreate.PullRequestName,
		Status:          api.PullRequestShortStatusOPEN,
	})

	if err != nil {
		return response, ResourceNotFoundError
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
	return response, nil

}

func (pr PullReqeustService) Merge(
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

func (pr PullReqeustService) Reassign(
	ctx context.Context,
	prToReassign api.PostPullRequestReassignJSONBody,
) (api.PullRequest, error) {
	prToChange, err := pr.prRepo.Get(ctx, prToReassign.PullRequestId)
	if err != nil {
		return prToChange, ResourceNotFoundError
	}

	if prToChange.Status == api.PullRequestStatusMERGED {
		return prToChange, CantReassignOnMergedPRError
	}

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

func (pr PullReqeustService) filterReviewers(
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
