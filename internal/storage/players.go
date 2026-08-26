package storage

import (
	"context"
	"database/sql"
	"example.com/othello-records/internal/domain"
	"fmt"
)

type PlayerRepository struct{ database *Database }

func NewPlayerRepository(database *Database) *PlayerRepository {
	return &PlayerRepository{database: database}
}

func (r *PlayerRepository) Save(ctx context.Context, player domain.PlayerProfile) error {
	if err := player.Validate(); err != nil {
		return err
	}
	_, err := r.database.db.ExecContext(ctx, `INSERT INTO PlayerProfile (id, display_name, rank, active, games, wins, losses, draws)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET display_name=excluded.display_name, rank=excluded.rank, active=excluded.active,
games=excluded.games, wins=excluded.wins, losses=excluded.losses, draws=excluded.draws`,
		player.ID, player.DisplayName, player.Rank, boolInt(player.Active), player.Games, player.Wins, player.Losses, player.Draws)
	if err != nil {
		return fmt.Errorf("save player: %w", err)
	}
	return nil
}

func (r *PlayerRepository) Get(ctx context.Context, id string) (domain.PlayerProfile, error) {
	var player domain.PlayerProfile
	var active int
	err := r.database.db.QueryRowContext(ctx, `SELECT id, display_name, rank, active, games, wins, losses, draws FROM PlayerProfile WHERE id = ?`, id).
		Scan(&player.ID, &player.DisplayName, &player.Rank, &active, &player.Games, &player.Wins, &player.Losses, &player.Draws)
	if err != nil {
		return domain.PlayerProfile{}, fmt.Errorf("get player: %w", err)
	}
	player.Active = active != 0
	return player, nil
}

func (r *PlayerRepository) List(ctx context.Context, activeOnly bool) ([]domain.PlayerProfile, error) {
	query := `SELECT id, display_name, rank, active, games, wins, losses, draws FROM PlayerProfile ORDER BY display_name`
	args := []any{}
	if activeOnly {
		query = `SELECT id, display_name, rank, active, games, wins, losses, draws FROM PlayerProfile WHERE active = 1 ORDER BY display_name`
	}
	rows, err := r.database.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	defer rows.Close()
	players := make([]domain.PlayerProfile, 0)
	for rows.Next() {
		var player domain.PlayerProfile
		var active int
		if err := rows.Scan(&player.ID, &player.DisplayName, &player.Rank, &active, &player.Games, &player.Wins, &player.Losses, &player.Draws); err != nil {
			return nil, fmt.Errorf("scan player: %w", err)
		}
		player.Active = active != 0
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate players: %w", err)
	}
	return players, nil
}

func (r *PlayerRepository) Delete(ctx context.Context, id string) error {
	result, err := r.database.db.ExecContext(ctx, `DELETE FROM PlayerProfile WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete player: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
