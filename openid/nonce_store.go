package openid

import (
	"errors"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var max_nonce_age = flag.Duration("openid-max-nonce-age",
	60*time.Second,
	"Maximum accepted age for openid nonces. The bigger, the more"+
		"memory is needed to store used nonces.")

type NonceStore interface {
	Accept(endpoint, nonce string) error
}

type Nonce struct {
	T time.Time
	S string
}

type RedisNonceStore struct {
	Pool *redis.Pool
}

func (store *RedisNonceStore) Accept(endpoint, nonce string) error {
	if len(nonce) < 20 || len(nonce) > 256 {
		return errors.New("Invalid nonce")
	}

	ts, err := time.Parse(time.RFC3339, nonce[0:20])
	if err != nil {
		return err
	}

	now := time.Now()
	diff := now.Sub(ts)
	if diff > *max_nonce_age {
		return fmt.Errorf("Nonce too old: %ds", diff.Seconds())
	}

	s := nonce[20:]
	_, err = store.Pool.Get().Do("SET", "steam_nonces:"+ts.String(), s)
	if err != nil {
		return err
	}
	_, err = store.Pool.Get().Do("EXPIRE", "steam_nonces:"+ts.String(), 60)
	if err != nil {
		return err
	}
	return nil
}
