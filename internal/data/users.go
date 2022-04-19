package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/sqlpipe/sqlpipe/internal/globals"
	"github.com/sqlpipe/sqlpipe/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)

var AnonymousUser = &User{}

type User struct {
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
	Username     string    `json:"username"`
	Password     password  `json:"-"`
	Version      int64     `json:"version"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 500, "username", "must not be more than 500 bytes long")

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	Etcd *clientv3.Client
}

func (m UserModel) Insert(user *User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), globals.EtcdTimeout)
	resp, err := m.Etcd.Get(ctx, fmt.Sprintf("sqlpipe/users/%v", user.Username))
	defer cancel()
	if err != nil {
		return err
	}
	if resp.Count != 0 {
		return ErrDuplicateUsername
	}

	creationTime := time.Now()
	user.CreatedAt = creationTime
	user.LastModified = creationTime

	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), globals.EtcdTimeout)
	_, err = m.Etcd.Put(
		ctx,
		fmt.Sprintf("sqlpipe/users/%v", user.Username),
		string(bytes),
	)
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) Get(username string) (*User, error) {

	context, cancel := context.WithTimeout(context.Background(), globals.EtcdTimeout)
	resp, err := m.Etcd.Get(context, fmt.Sprintf("sqlpipe/users/%v", username))
	defer cancel()
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, ErrRecordNotFound
	}

	var user User
	err = json.Unmarshal(resp.Kvs[0].Value, &user)
	if err != nil {
		return nil, err
	}
	user.Version = resp.Kvs[0].Version

	return &user, nil
}

func (m UserModel) GetByUsername(username string) (user *User, err error) {
	// query := `
	//     SELECT id, created_at, name, email, password_hash, activated, version
	//     FROM users
	//     WHERE email = $1`

	// var user User

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err := m.DB.QueryRowContext(ctx, query, email).Scan(
	// 	&user.Id,
	// 	&user.CreatedAt,
	// 	&user.Name,
	// 	&user.Email,
	// 	&user.Password.hash,
	// 	&user.Activated,
	// 	&user.Version,
	// )

	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return nil, ErrRecordNotFound
	// 	default:
	// 		return nil, err
	// 	}
	// }

	return user, nil
}

func (m UserModel) Update(user *User) (err error) {
	// query := `
	//     UPDATE users
	//     SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
	//     WHERE id = $5 AND version = $6
	//     RETURNING version`

	// args := []interface{}{
	// 	user.Name,
	// 	user.Email,
	// 	user.Password.hash,
	// 	user.Activated,
	// 	user.Id,
	// 	user.Version,
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	// if err != nil {
	// 	switch {
	// 	case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
	// 		return ErrDuplicateEmail
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return ErrEditConflict
	// 	default:
	// 		return err
	// 	}
	// }

	return err
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (user *User, err error) {
	// tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// query := `
	//     SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
	//     FROM users
	//     INNER JOIN tokens
	//     ON users.id = tokens.user_id
	//     WHERE tokens.hash = $1
	//     AND tokens.scope = $2
	//     AND tokens.expiry > $3`

	// args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	// var user User

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err := m.DB.QueryRowContext(ctx, query, args...).Scan(
	// 	&user.Id,
	// 	&user.CreatedAt,
	// 	&user.Name,
	// 	&user.Email,
	// 	&user.Password.hash,
	// 	&user.Activated,
	// 	&user.Version,
	// )
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return nil, ErrRecordNotFound
	// 	default:
	// 		return nil, err
	// 	}
	// }

	return user, nil
}