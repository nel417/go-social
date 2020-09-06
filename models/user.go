package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
)

type User struct {
	key string
}

func NewUser(username string, hash []byte) (*User, error) {
	ctx := context.TODO()

	id, err := client.Incr(ctx, "user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "username", username)
	pipe.HSet(ctx, key, "hash", hash)
	pipe.HSet(ctx, "user:by-username", username, id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &User{key}, nil
}

func (user *User) GetUsername() (string, error) {
	ctx := context.TODO()

	return client.HGet(ctx, user.key, "username").Result()
}

func (user *User) GetHash() ([]byte, error) {
	ctx := context.TODO()

	return client.HGet(ctx, user.key, "hash").Bytes()
}

func (user *User) Authenticate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err

}

func GetUserByUsername(username string) (*User, error) {
	ctx := context.TODO()

	id, err := client.HGet(ctx, "user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound

	} else if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	return &User{key}, nil

}

func AuthenticateUser(username, password string) error {
	user, err := GetUserByUsername(username)
	if err != nil {
		return err
	}
	return user.Authenticate(password)

}

func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return err
	}
	_, err = NewUser(username, hash)
	return err
}
