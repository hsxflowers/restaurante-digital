package db

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hsxflowers/restaurante-digital/workers"
	"github.com/labstack/gommon/log"
)

type SQLStore struct {
	db *sql.DB
}

type CatDatabase interface {
	CreatePedido(ctx context.Context, userID, itemID string, valor float64)
	GetItem(ctx context.Context, itemId string) (*workers.Item, error)
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

func (s *SQLStore) GetItem(ctx context.Context, itemId string) (*workers.Item, error) {
	var item workers.Item
	var valor float64

	query := `
		SELECT item_id, nome, tempo_corte, tempo_grelha, tempo_montagem, tempo_bebida, valor
		FROM pedido
		WHERE item_id = $1
		ORDER BY RANDOM() LIMIT 1`

	row := s.db.QueryRowContext(ctx, query, itemId)

	err := row.Scan(&item.ItemId, &item.Nome, &item.TempoCorte, &item.TempoGrelha, &item.TempoMontagem, &item.TempoBebida, &valor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.New(exceptions.ErrPedidoNotFound, err)
		}
		log.Error("Error fetching pedido from database: ", err)
		return nil, exceptions.New(exceptions.ErrInternalServer, err)
	}

	item.Valor = valor
	return &item, nil
}

func (s *SQLStore) CreatePedido(ctx context.Context, userID, itemID string, valor float64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO pedido (user_id, item_id, valor) 
		VALUES ($1, $2, $3)`,
		userID, itemID, valor,
	)
	if err != nil {
		log.Error("Error creating pedido in database: ", err)
		return exceptions.New(exceptions.ErrInternalServer, err)
	}
	return nil
}
