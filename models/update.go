package models

import (
	"context"
	"fmt"
)

type Update struct {
	key string
}

func NewUpdate(userId int64, body string) (*Update, error) {
	ctx := context.TODO()

	id, err := client.Incr(ctx, "update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "user_id", userId)
	pipe.HSet(ctx, key, "body", body)
	pipe.LPush(ctx, "updates", id)
	pipe.LPush(ctx, fmt.Sprintf("user:%d:updates", userId), id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &Update{key}, nil
}

func (update *Update) GetBody() (string, error) {
	ctx := context.TODO()

	return client.HGet(ctx, update.key, "body").Result()
}

func (update *Update) GetUser() (*User, error) {
	ctx := context.TODO()

	userId, err := client.HGet(ctx, update.key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserById(userId)
}

func GetAllUpdates() ([]*Update, error) {
	ctx := context.TODO()

	updateIds, err := client.LRange(ctx, "updates", 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIds))
	for i, id := range updateIds {
		key := "update:" + id
		updates[i] = &Update{key}
	}
	return updates, nil
}

func GetUpdates(userId int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userId)
	ctx := context.TODO()

	updateIds, err := client.LRange(ctx, key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIds))
	for i, id := range updateIds {
		key := "update:" + id
		updates[i] = &Update{key}
	}
	return updates, nil
}

func PostUpdate(userId int64, body string) error {
	_, err := NewUpdate(userId, body)
	return err
}
