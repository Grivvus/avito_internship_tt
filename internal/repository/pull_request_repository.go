package repository

import (
	"context"
	"fmt"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/jackc/pgx/v5"
)

type PullRequestRepository struct {
	conn *pgx.Conn
}

func NewPullRequestRepository(conn *pgx.Conn) PullRequestRepository {
	return PullRequestRepository{
		conn: conn,
	}
}

func (pr PullRequestRepository) FindByReviewer(
	ctx context.Context,
	reviewerId string,
) ([]api.PullRequest, error) {
	const query = `select 
		id
		from public."pull_requests" inner join public."pull_request_reviewers"
			on "pull_request_reviewers".pr_id = "pull_requests".id
		where reviewer_id = $1
	`

	rows, err := pr.conn.Query(ctx, query, reviewerId)
	if err != nil {
		return nil, fmt.Errorf("On query execution: %w", err)
	}
	defer rows.Close()

	pullRequests := make([]api.PullRequest, 0, 10)

	for rows.Next() {
		var prId string
		err := rows.Scan(&prId)
		if err != nil {
			return nil, fmt.Errorf("On scaning: %w", err)
		}
		pullRequest, err := pr.Get(ctx, prId)
		if err != nil {
			return nil, fmt.Errorf("On getting pr: %w", err)
		}
		pullRequests = append(pullRequests, pullRequest)
	}
	return pullRequests, nil
}

func (pr PullRequestRepository) Create(
	ctx context.Context,
	model api.PullRequestShort,
) error {
	const query = `insert into public."pull_requests"
		(id, title, author_id)
		values ($1, $2, $3)
	`
	rows, err := pr.conn.Query(ctx, query, model.PullRequestId, model.PullRequestName, model.AuthorId)
	if err != nil {
		return fmt.Errorf("On query execution: %w", err)
	}
	defer rows.Close()
	return nil
}

func (pr PullRequestRepository) Get(
	ctx context.Context,
	prId string,
) (api.PullRequest, error) {

	var responseModel api.PullRequest

	const selectPRInfoQuery = `
		select id, title, author_id, status, merged_at, created_at
			from public."pull_requests"
			where id = $1
	`

	const selectPRReviewersQuery = `
		select reviewer_id from public."pull_request_reviewers"
			where pr_id = $1
	`

	err := pr.conn.QueryRow(ctx, selectPRInfoQuery, prId).Scan(
		&responseModel.PullRequestId, &responseModel.PullRequestName,
		&responseModel.AuthorId, &responseModel.Status,
		&responseModel.MergedAt, &responseModel.CreatedAt,
	)
	if err != nil {
		return responseModel, fmt.Errorf("Can't scan result: %w", err)
	}

	prReviewersRows, err := pr.conn.Query(ctx, selectPRReviewersQuery, prId)
	if err != nil {
		return responseModel, fmt.Errorf("On query execution: %w", err)
	}
	defer prReviewersRows.Close()

	for prReviewersRows.Next() {
		var reviewerId string
		err := prReviewersRows.Scan(&reviewerId)
		if err != nil {
			return responseModel, fmt.Errorf("Can't scan result: %w", err)
		}
		responseModel.AssignedReviewers = append(responseModel.AssignedReviewers, reviewerId)
	}
	return responseModel, nil
}

func (pr PullRequestRepository) Merge(
	ctx context.Context,
	pullRequest api.PullRequest,
) (api.PullRequest, error) {
	const query = `update public."pull_requests" 
		(status, merged_at)
		values ('MERGED', now())
		where id = $1
	`

	rows, err := pr.conn.Query(ctx, query, pullRequest.PullRequestId)
	if err != nil {
		return api.PullRequest{}, fmt.Errorf("On query execution: %w", err)
	}
	defer rows.Close()

	return pr.Get(ctx, pullRequest.PullRequestId)
}

func (pr PullRequestRepository) AssignReviewers(
	ctx context.Context,
	pullRequestId string,
	reviewersId []string,
) error {
	const query = `insert into public."pull_request_reviewers"
		(pr_id, reviewer_id)
		values ($1, $2)
	`

	tx, err := pr.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("On begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, reviewerId := range reviewersId {
		_, err := tx.Exec(ctx, query, pullRequestId, reviewerId)
		if err != nil {
			return fmt.Errorf("On insert: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("On commit: %w", err)
	}
	return nil
}

func (pr PullRequestRepository) ReassignReviewer(
	ctx context.Context,
	pullRequestId, oldReviewerId, newReviewerId string,
) error {
	const deleteQuery = `delete from public."pull_request_reviewers"
		where reviewer_id = $1`
	const insertQuery = `insert into public."pull_request_reviewers"
		(pr_id, reviewer_id) values ($1, $2)`

	tx, err := pr.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("On transaction begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteQuery, oldReviewerId)
	if err != nil {
		return fmt.Errorf("On delete query: %w", err)
	}

	_, err = tx.Exec(ctx, insertQuery, pullRequestId, newReviewerId)
	if err != nil {
		return fmt.Errorf("On insert query: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("On commit: %w", err)
	}
	return nil
}
