package repository

import (
	"context"
	"errors"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/jackc/pgx/v5"
)

var UserNotFoundErr error = errors.New("No user was found with this userId")

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) UserRepository {
	return UserRepository{
		conn: conn,
	}
}

func (ur UserRepository) Create(ctx context.Context, user api.User) error {
	const query = `insert into public."users" (id, username, is_active, team)
	               values ($1, $2, $3, $4)`

	rows, err := ur.conn.Query(
		ctx, query, user.UserId,
		user.Username, user.IsActive, user.TeamName,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (ur UserRepository) Get(ctx context.Context, userId string) (api.User, error) {
	var userFromDB api.User

	const query = `select id, username, is_active, team from public."users" where id = $1`

	err := ur.conn.QueryRow(ctx, query, userId).Scan(
		&userFromDB.UserId, &userFromDB.Username,
		&userFromDB.IsActive, &userFromDB.TeamName,
	)

	if err != nil {
		return userFromDB, UserNotFoundErr
	}

	return userFromDB, nil
}

func (ur UserRepository) Update(ctx context.Context, updatedUser api.User) error {
	const query = `
		update public."users" set
			username = $1,
			is_active = $2,
			team = $3
		where id = $4
	`

	rows, err := ur.conn.Query(
		ctx, query,
		updatedUser.Username, updatedUser.IsActive,
		updatedUser.TeamName, updatedUser.UserId,
	)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (ur UserRepository) FindByTeam(ctx context.Context, teamName string) ([]api.User, error) {
	const query = `select id, username, is_active, team from public."users" where team = $1`

	rows, err := ur.conn.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	selectedUsers := make([]api.User, 0, 8)

	for rows.Next() {
		var user api.User
		if err := rows.Scan(&user.UserId, &user.Username, &user.IsActive, &user.TeamName); err != nil {
			return nil, err
		}
		selectedUsers = append(selectedUsers, user)
	}

	return selectedUsers, nil
}

func (ur UserRepository) FindOtherMembers(
	ctx context.Context,
	memberId string,
) ([]api.User, error) {
	const query = `
	select id, username, is_active, team from public."users"
	where team = (select team from public."users" where id = $1)
	and id <> $1
	`

	rows, err := ur.conn.Query(ctx, query, memberId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	otherMembers := make([]api.User, 0, 8)
	for rows.Next() {
		var user api.User
		if err := rows.Scan(&user.UserId, &user.Username, &user.IsActive, &user.TeamName); err != nil {
			return nil, err
		}
		otherMembers = append(otherMembers, user)
	}
	return otherMembers, nil
}

func (ur UserRepository) CreateOrUpdate(ctx context.Context, user api.User) error {
	_, err := ur.Get(ctx, user.UserId)
	if err != nil {
		// better to match err here
		err := ur.Create(ctx, user)
		if err != nil {
			return err
		}
		return nil
	}

	err = ur.Update(ctx, user)
	return err
}
