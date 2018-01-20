package tawa

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	_url "net/url"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Variables map[string]interface{} `json:"variables"`
	Tags      []string               `json:"tags"`
	Playbook  string                 `json:"playbook"`
	Hosts     []string               `json:"hosts"`
	Callback  string                 `json:"callback"`
	Id        string                 `json:"id"`
}

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

type Response struct {
	Id   uuid.UUID
	Chan chan interface{}
}

func (t *Tawa) Send(msg *Event) (*Response, error) {
	id := uuid.NewV4()
	msg.Id = id.String()
	m, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	r := &Response{
		Id:   id,
		Chan: make(chan interface{}),
	}
	err = t.redis.LPush("tawa", m).Err()
	if err != nil {
		return r, err
	}
	err = t.redis.Set(fmt.Sprintf("state:%s", id), "queued", 0).Err()
	if err != nil {
		return r, err
	}
	go func(client *redis.Client, r *Response) {
		cb := fmt.Sprintf("cb:%s", r.Id)
		result, err := client.BLPop(300*time.Second, cb).Result()
		if err != nil {
			//FIXME
			panic(err)
		}
		r.Chan <- result
	}(t.redis, r)
	return r, nil
}
