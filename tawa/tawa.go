package tawa

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	_url "net/url"
	"strconv"
	"strings"
)

type Tawa struct {
	redis *redis.Client
}

func New(url string) (*Tawa, error) {
	u, err := _url.Parse(url)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "redis" {
		return nil, errors.New("Please, use redis:// scheme")
	}
	var pass string
	user := u.User
	if user != nil {
		p, ok := u.User.Password()
		if ok {
			pass = p
		} else {
			pass = ""
		}
	}
	var db int64
	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		db = 0
	} else {
		db, err = strconv.ParseInt(path, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid db name %s", err)
		}
	}
	c := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: pass,
		DB:       int(db),
	})
	return &Tawa{
		redis: c,
	}, nil
}
