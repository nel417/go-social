package models

import "context"

func GetComments() ([]string, error) {
	// define context
	ctx := context.TODO()
	return client.LRange(ctx, "comments", 0, 10).Result()

}

func PostComment(comment string) error {
	//define context
	ctx := context.TODO()
	return client.LPush(ctx, "comments", comment).Err()
}
