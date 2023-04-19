package providers

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MysqlProvider struct {
	conn *sqlx.DB
	tx   *sqlx.Tx
}

func NewMysqlProvider(dsn string) (*MysqlProvider, error) {
	conn, err := sqlx.Open("mysql", dsn)
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

	return &MysqlProvider{
		conn: conn,
		tx:   tx,
	}, nil
}

func (m *MysqlProvider) GetAll(tableName string) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	rows, err := m.tx.Queryx(fmt.Sprintf("SELECT * FROM %s", tableName))
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

func (m *MysqlProvider) Truncate(tableName string) error {
	return m.tx.QueryRowx(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Err()
}

func (m *MysqlProvider) Insert(tableName string, columns []string, values []any) error {
	rowsCountForInsert := len(values) / len(columns)
	sql := m.tx.Rebind(generateInsertSQL(tableName, columns, rowsCountForInsert))
	return m.tx.QueryRowx(sql, values...).Err()
}

func (m *MysqlProvider) Commit() error {
	defer m.conn.Close()
	return m.tx.Commit()
}

func (m *MysqlProvider) Rollback() error {
	defer m.conn.Close()
	return m.tx.Rollback()
}
