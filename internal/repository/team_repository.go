package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var NoTeamFoundError error = errors.New("No team with this name was found")

type TeamRepository struct {
	conn *pgx.Conn
}

func NewTeamRepository(conn *pgx.Conn) TeamRepository {
	return TeamRepository{
		conn: conn,
	}
}

func (tr TeamRepository) Create(ctx context.Context, teamName string) error {
	const query = `insert into public."teams" (name) values ($1)`
	rows, err := tr.conn.Query(ctx, query, teamName)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (tr TeamRepository) Get(ctx context.Context, teamName string) (string, error) {
	var teamNameDBCheck string

	const query = `select name from public."teams" where teams.name = $1`

	err := tr.conn.QueryRow(ctx, query, teamName).Scan(&teamNameDBCheck)
	if err != nil {
		return "", NoTeamFoundError
	}

	return teamNameDBCheck, nil
}
