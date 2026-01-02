package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netbill/pgx"
)

const profilesTable = "profiles"

type Profile struct {
	AccountID   uuid.UUID `db:"account_id"`
	Username    string    `db:"username"`
	Official    bool      `db:"official"`
	Pseudonym   *string   `db:"pseudonym"`
	Description *string   `db:"description"`
	Avatar      *string   `db:"avatar"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (p Profile) scan(row sq.RowScanner) error {
	err := row.Scan(
		&p.AccountID,
		&p.Username,
		&p.Official,
		&p.Pseudonym,
		&p.Description,
		&p.Avatar,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scanning profile: %w", err)
	}
	return nil
}

type ProfilesQ struct {
	db       pgx.DBTX
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewProfilesQ(db pgx.DBTX) ProfilesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ProfilesQ{
		db:       db,
		selector: builder.Select("*").From(profilesTable),
		inserter: builder.Insert(profilesTable),
		updater:  builder.Update(profilesTable),
		deleter:  builder.Delete(profilesTable),
		counter:  builder.Select("COUNT(*) AS count").From(profilesTable),
	}
}

func (q ProfilesQ) Insert(ctx context.Context, input Profile) (Profile, error) {
	values := map[string]interface{}{
		"account_id":  input.AccountID,
		"username":    input.Username,
		"official":    input.Official,
		"pseudonym":   input.Pseudonym,
		"description": input.Description,
		"avatar":      input.Avatar,
		"created_at":  input.CreatedAt,
		"updated_at":  input.UpdatedAt,
	}

	query, args, err := q.inserter.
		SetMap(values).
		Suffix("RETURNING account_id, username, official, pseudonym, description, avatar, created_at, updated_at").
		ToSql()
	if err != nil {
		return Profile{}, fmt.Errorf("building insert query for %s: %w", profilesTable, err)
	}

	row := q.db.QueryRowContext(ctx, query, args...)

	var p Profile
	err = p.scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, nil
		}
		return Profile{}, err
	}

	return p, nil
}

func (q ProfilesQ) Update(ctx context.Context) ([]Profile, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.
		Suffix("RETURNING account_id, username, official, pseudonym, description, avatar, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building update query for %s: %w", profilesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Profile
	for rows.Next() {
		var p Profile
		err = p.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning profile: %w", err)
		}
		out = append(out, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q ProfilesQ) UpdateOne(ctx context.Context) (Profile, error) {
	rows, err := q.Update(ctx)
	if err != nil {
		return Profile{}, err
	}
	if len(rows) != 1 {
		return Profile{}, fmt.Errorf("expected 1 profile to be updated, got %d", len(rows))
	}
	return rows[0], nil
}

func (q ProfilesQ) UpdateUsername(username string) ProfilesQ {
	q.updater = q.updater.Set("username", username)
	return q
}

func (q ProfilesQ) UpdateOfficial(official bool) ProfilesQ {
	q.updater = q.updater.Set("official", official)
	return q
}

func (q ProfilesQ) UpdatePseudonym(pseudonym *string) ProfilesQ {
	q.updater = q.updater.Set("pseudonym", pseudonym)
	return q
}

func (q ProfilesQ) UpdateDescription(description *string) ProfilesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q ProfilesQ) UpdateAvatar(avatar *string) ProfilesQ {
	q.updater = q.updater.Set("avatar", avatar)
	return q
}

func (q ProfilesQ) Get(ctx context.Context) (Profile, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Profile{}, fmt.Errorf("building get query for %s: %w", profilesTable, err)
	}

	row := q.db.QueryRowContext(ctx, query, args...)

	var p Profile
	err = row.Scan(
		&p.AccountID,
		&p.Username,
		&p.Official,
		&p.Pseudonym,
		&p.Description,
		&p.Avatar,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, nil
		}
		return Profile{}, err
	}

	return p, nil
}

func (q ProfilesQ) Select(ctx context.Context) ([]Profile, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", profilesTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Profile
	for rows.Next() {
		var p Profile
		err = rows.Scan(
			&p.AccountID,
			&p.Username,
			&p.Official,
			&p.Pseudonym,
			&p.Description,
			&p.Avatar,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning profile: %w", err)
		}
		out = append(out, p)
	}

	return out, nil
}

func (q ProfilesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", profilesTable, err)
	}

	_, err = q.db.ExecContext(ctx, query, args...)

	return err
}

func (q ProfilesQ) FilterAccountID(accountID ...uuid.UUID) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"account_id": accountID})
	return q
}

func (q ProfilesQ) FilterUsername(username string) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"username": username})
	q.counter = q.counter.Where(sq.Eq{"username": username})
	q.deleter = q.deleter.Where(sq.Eq{"username": username})
	q.updater = q.updater.Where(sq.Eq{"username": username})
	return q
}

func (q ProfilesQ) FilterOfficial(official bool) ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"official": official})
	q.counter = q.counter.Where(sq.Eq{"official": official})
	q.deleter = q.deleter.Where(sq.Eq{"official": official})
	q.updater = q.updater.Where(sq.Eq{"official": official})
	return q
}

func (q ProfilesQ) FilterLikePseudonym(pseudonym string) ProfilesQ {
	q.selector = q.selector.Where(sq.ILike{"pseudonym": "%" + pseudonym + "%"})
	q.counter = q.counter.Where(sq.ILike{"pseudonym": "%" + pseudonym + "%"})
	q.updater = q.updater.Where(sq.ILike{"pseudonym": "%" + pseudonym + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"pseudonym": "%" + pseudonym + "%"})

	return q
}

func (q ProfilesQ) FilterLikeUsername(username string) ProfilesQ {
	q.selector = q.selector.Where(sq.ILike{"username": "%" + username + "%"})
	q.counter = q.counter.Where(sq.ILike{"username": "%" + username + "%"})
	q.updater = q.updater.Where(sq.ILike{"username": "%" + username + "%"})
	q.deleter = q.deleter.Where(sq.ILike{"username": "%" + username + "%"})

	return q
}

func (q ProfilesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", profilesTable, err)
	}

	var count uint
	err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q ProfilesQ) Page(limit, offset uint) ProfilesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q ProfilesQ) OrderCreatedAt(ascending bool) ProfilesQ {
	if ascending {
		q.selector = q.selector.OrderBy("created_at ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC")
	}
	return q
}
