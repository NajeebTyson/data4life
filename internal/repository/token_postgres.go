package repository

import (
	"data4life/pkg/token"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "postgres123"
	dbname   = "data4life"
)

type TokenStorePostgres struct {
	conn *sql.DB
}

func (s *TokenStorePostgres) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *TokenStorePostgres) AddToken(t token.Token) error {
	sqlStatement := `
		INSERT INTO tokens (token)
		VALUES ($1)
	`
	err := s.conn.QueryRow(sqlStatement, t).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenStorePostgres) AddTokenBatch(tokens []token.Token) error {
	var builder strings.Builder
	for i := 1; i <= len(tokens); i++ {
		builder.WriteString(fmt.Sprintf("($%d),", i))
	}
	q := builder.String()
	q = q[:len(q)-1] // to remove the comma from the last
	sqlStatement := fmt.Sprintf("INSERT INTO tokens (token) VALUES %s;", q)

	_, err := s.conn.Exec(sqlStatement, token.ConvertToInterfaceSlice(tokens)...)
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenStorePostgres) GetToken(t token.Token) (*token.Token, error) {
	var queryToken token.Token

	sqlStatement := `
		SELECT token FROM tokens WHERE token=$1;
	`

	row := s.conn.QueryRow(sqlStatement, t)
	switch err := row.Scan(&queryToken); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &queryToken, nil
	default:
		return nil, err
	}
}

func (s *TokenStorePostgres) DeleteToken(t token.Token) error {
	sqlStatement := `
		DELETE FROM tokens
		WHERE token = $1;
	`

	_, err := s.conn.Exec(sqlStatement, t)
	if err != nil {
		return err
	}
	return nil
}

func NewTokenStorePostgres() (*TokenStorePostgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &TokenStorePostgres{
		conn: db,
	}, nil

}
