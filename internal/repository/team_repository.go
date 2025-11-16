package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var NoTeamFoundError error = errors.New("No team with this name was found")

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{
		pool: pool,
	}
}

func (tr *TeamRepository) Create(ctx context.Context, teamName string) error {
	const query = `insert into public."teams" (name) values ($1)`
	_, err := tr.pool.Exec(ctx, query, teamName)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TeamRepository) Get(ctx context.Context, teamName string) (string, error) {
	var teamNameDBCheck string

	const query = `select name from public."teams" where teams.name = $1`

	err := tr.pool.QueryRow(ctx, query, teamName).Scan(&teamNameDBCheck)
	if err != nil {
		return "", NoTeamFoundError
	}

	return teamNameDBCheck, nil
}
