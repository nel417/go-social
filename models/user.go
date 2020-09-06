package models

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
)

func AuthenticateUser(username, password string) error {
	ctx := context.TODO()
	hash, err := client.Get(ctx, "user:"+username).Bytes()
	if err == redis.Nil {
		return ErrUserNotFound

	} else if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return ErrInvalidLogin
	}
	return nil
}

func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return err
	}
	ctx := context.TODO()
	return client.Set(ctx, "user:"+username, hash, 0).Err()
}
