package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

var db *sql.DB

func Connect(cred Credentials) error {
	dest := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cred.User, cred.Password, cred.Host, cred.Port, cred.Database)

	con, err := sql.Open("postgres", dest)
	if err != nil {
		return err
	}

	db = con

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return nil
}

func Disconnect() error {
	if db == nil {
		return errors.New("Nothing to close")
	}

	db.Close()

	return nil
}

func New(id int64) error {
	q := `INSERT INTO UserStats (id) VALUES ($1);`
	_, err := db.Exec(q, id)

	return err
}

func IncreaseReq(id int64) error {
	q := `UPDATE UserStats SET amount = amount + 1 WHERE id = $1;`
	_, err := db.Exec(q, id)

	return err
}

func FirstReq(id int64) error {
	q := `UPDATE UserStats SET firstReq = NOW() WHERE id = $1;`
	_, err := db.Exec(q, id)

	return err
}

func LastPair(id int64, pair string) error {
	q := `UPDATE UserStats SET lastPair = $1 WHERE id = $2;`
	_, err := db.Exec(q, pair, id)

	return err
}

func UserStats(id int64) (UserModel, error) {
	q := `SELECT * FROM UserStats WHERE id = $1`
	rows, err := db.Query(q, id)
	if err != nil {
		return UserModel{}, err
	}

	defer rows.Close()

	res := UserModel{}
	if rows.Next() {
		err = rows.Scan(&res.ID, &res.FirstReq, &res.ReqAmount, &res.LastPair)
		if err != nil {
			return UserModel{}, err
		}

		return res, nil
	}

	return UserModel{}, errors.New("no rows")
}

func StatsExist(id int64) bool {
	_, err := UserStats(id)

	return err == nil
}
