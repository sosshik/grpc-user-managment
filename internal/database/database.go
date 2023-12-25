package database

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"git.foxminded.ua/foxstudent106264/task-4.1/internal/domain"
	"git.foxminded.ua/foxstudent106264/task-4.1/pkg/config"
	proto "git.foxminded.ua/foxstudent106264/task-4.1/protos/gen/go/user_service"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type DBConfig struct {
	DbUrl       string `env:"DATABASE_URL"`
	ReconnTime  int    `env:"RECONN_TIME" envDefault:"5"`
	ConnCheck   bool   `env:"CONN_CHECK" envDefault:"true"`
	ReconnTries int    `env:"RECONN_TRIES" envDefault:"5"`
}

type Database struct {
	config *DBConfig
	DB     *sql.DB
}

var once sync.Once

var dbinstance *Database

func NewDatabase(cfg *config.Config) (*Database, error) {

	if dbinstance == nil {
		once.Do(func() {
			db, err := sql.Open("postgres", cfg.DbUrl)
			if err != nil {
				log.Warnf("unable to create db instance: %s", err)
			}

			dbinstance = &Database{&DBConfig{cfg.DbUrl, cfg.ReconnTime, cfg.ConnCheck, cfg.ReconnTries}, db}

			if cfg.ConnCheck {
				go dbinstance.connectionCheck(cfg.DbUrl)
			}
		})

	}

	return dbinstance, nil
}

func (d *Database) connectionCheck(conn string) {
	log.Info("Connection check started")
	var i int
	for {
		time.Sleep(time.Duration(d.config.ReconnTime) * time.Second)
		if err := d.DB.Ping(); err != nil {
			log.Warnf("Lost connection to Database. Attempting to reconnect.")
			if err := d.DB.Close(); err != nil {
				log.Warnf("Error while disconecting: %s", err)
				continue
			}
			if i <= d.config.ReconnTries {
				d.DB, err = sql.Open("postgres", conn)
				if err != nil {
					log.Warnf("Failed to reconnect: %s", err)
					i++
				} else {
					log.Infof("Reconnected to PostgreSQL!")
					i = 0
				}
			} else {
				break
			}

		}
	}
}

func (d *Database) CreateUser(user domain.UserProfileDTO) error {

	_, err := d.DB.Exec(`
	INSERT INTO users (oid, nickname, email, first_name, last_name, password, created_at, updated_at, state)
	VALUES ($1, $2, $3, $4, $5, $6, $7,$8, $9);
	`, user.OID, user.Nickname, user.Email, user.FirstName, user.LastName, user.Password, user.CreatedAt, user.UpdatedAt, user.State)
	if err != nil {
		return fmt.Errorf("unable to execute query to DB: %w", err)
	}
	return nil
}

func (d *Database) GetUserByEmail(email string) (domain.UserProfileDTO, error) {
	var user domain.UserProfileDTO
	err := d.DB.QueryRow(`
	SELECT oid, nickname, email, first_name, last_name FROM users
	WHERE email = $1;
	`, email).Scan(&user.OID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName)
	if err != nil && err != sql.ErrNoRows {
		return domain.UserProfileDTO{}, fmt.Errorf("unable to execute query to DB: %w", err)
	}
	return user, nil
}

func (d *Database) GetUserByID(oid uuid.UUID) (domain.UserProfileDTO, error) {
	var user domain.UserProfileDTO
	err := d.DB.QueryRow(`
	SELECT oid, nickname, email, first_name, last_name FROM users
	WHERE oid = $1;
	`, oid).Scan(&user.OID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName)
	if err != nil && err != sql.ErrNoRows {
		return domain.UserProfileDTO{}, fmt.Errorf("unable to execute query to DB: %w", err)
	}
	return user, nil
}

func (d *Database) GetUsers() ([]*proto.UserInfo, error) {
	rows, err := d.DB.Query(`
	SELECT oid, nickname, email, first_name, last_name
    FROM users
	`)
	if err != nil {
		return []*proto.UserInfo{}, fmt.Errorf("unable to execute query to DB: %w", err)
	}
	defer rows.Close()

	var users []*proto.UserInfo

	for rows.Next() {
		var user domain.UserProfileDTO
		err := rows.Scan(&user.OID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			return []*proto.UserInfo{}, fmt.Errorf("unable to scan row from DB: %w", err)
		}

		users = append(users, &proto.UserInfo{
			Oid:       &proto.UUID{Value: user.OID.String()},
			Email:     user.Email,
			Nickname:  user.Nickname,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}
	return users, nil
}

func (d *Database) UpdateUser(user domain.UserProfileDTO) error {

	_, err := d.DB.Exec(`
	UPDATE users
	SET nickname=$1, email = $2, first_name=$3, last_name=$4, updated_at=$5
	WHERE oid=$6;
	`, user.Nickname, user.Email, user.FirstName, user.LastName, time.Now().UTC(), user.OID)
	if err != nil {
		return fmt.Errorf("unable to execute query to DB: %w", err)
	}
	return nil
}

func (d *Database) DeleteUser(oid uuid.UUID) error {

	_, err := d.DB.Exec(`
	DELETE FROM users
    WHERE oid = $1;
	`, oid)
	if err != nil {
		return fmt.Errorf("unable to execute query to DB: %w", err)
	}

	return nil
}
