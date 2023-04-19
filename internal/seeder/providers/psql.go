package providers

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PsqlProvider struct {
	conn *sqlx.DB
	tx   *sqlx.Tx
}

func NewPsqlProvider(dsn string) (*PsqlProvider, error) {
	conn, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	tx, err := conn.Beginx()
	if err != nil {
		return nil, err
	}

	return &PsqlProvider{
		conn: conn,
		tx:   tx,
	}, nil
}

func (p *PsqlProvider) GetAll(tableName string) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	rows, err := p.tx.Queryx(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		val := make(map[string]any, 0)
		scanErr := rows.MapScan(val)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, val)
	}

	return result, nil
}

func (p *PsqlProvider) Truncate(tableName string) error {
	return p.tx.QueryRowx("TRUNCATE TABLE " + tableName).Err()
}

func (p *PsqlProvider) Insert(tableName string, columns []string, values []any) error {
	rowsCountForInsert := len(values) / len(columns)
	sql := p.tx.Rebind(generateInsertSQL(tableName, columns, rowsCountForInsert))
	return p.tx.QueryRowx(sql, values...).Err()
}

func (p *PsqlProvider) Commit() error {
	defer p.conn.Close()
	return p.tx.Commit()
}

func (p *PsqlProvider) Rollback() error {
	defer p.conn.Close()
	return p.tx.Rollback()
}
