package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"hammerpost/model"
)

func getConn(pgDSN string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), pgDSN)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func UpdatePgParameter(pgDSN string, param []model.Param) error {
	conn, err := getConn(pgDSN)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	for _, p := range param {
		_, err = conn.Exec(context.Background(), fmt.Sprintf("ALTER SYSTEM SET %s TO '%s'", p.Name, p.Value))
		if err != nil {
			return err
		}
	}
	return nil
}
