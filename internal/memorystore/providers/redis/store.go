package redis

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	mfaSessionPrefix = "mfa_sess_"
	otpPrefix        = "otp_"
)

// SetUserSession sets the user session for given user identifier in form recipe:user_id
func (c *memoryStoreProvider) SetUserSession(userId, key, token string) error {
	err := c.store.Set(c.ctx, fmt.Sprintf("%s:%s", userId, key), token, -1).Err()
	if err != nil {
		log.Debug().Err(err).Msg("Error saving user session to redis")
		return err
	}
	return nil
}

// GetUserSession returns the user session from redis store.
func (c *memoryStoreProvider) GetUserSession(userId, key string) (string, error) {
	data, err := c.store.Get(c.ctx, fmt.Sprintf("%s:%s", userId, key)).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

// DeleteUserSession deletes the user session from redis store.
func (c *memoryStoreProvider) DeleteUserSession(userId, key string) error {
	if err := c.store.Del(c.ctx, fmt.Sprintf("%s:%s", userId, key)).Err(); err != nil {
		log.Debug().Err(err).Msg("Error deleting user session from redis")
		// continue
	}
	return nil
}

// DeleteAllUserSessions deletes all the user session from redis
func (c *memoryStoreProvider) DeleteAllUserSessions(userID string) error {
	res := c.store.Keys(c.ctx, fmt.Sprintf("*%s*", userID))
	if res.Err() != nil {
		log.Debug().Err(res.Err()).Msg("Error getting all user sessions from redis")
		return res.Err()
	}
	keys := res.Val()
	for _, key := range keys {
		fmt.Println("Deleting key", key)
		err := c.store.Del(c.ctx, key).Err()
		if err != nil {
			log.Debug().Err(err).Msg("Error deleting all user sessions from redis")
			continue
		}
	}
	return nil
}

// DeleteSessionForNamespace to delete session for a given namespace example google,github
func (c *memoryStoreProvider) DeleteSessionForNamespace(namespace string) error {
	res := c.store.Keys(c.ctx, fmt.Sprintf("%s:*", namespace))
	if res.Err() != nil {
		log.Debug().Err(res.Err()).Msg("Error getting all user sessions from redis")
		return res.Err()
	}
	keys := res.Val()
	for _, key := range keys {
		err := c.store.Del(c.ctx, key).Err()
		if err != nil {
			log.Debug().Err(err).Msg("Error deleting all user sessions from redis")
			continue
		}
	}
	return nil
}

// SetMfaSession sets the mfa session with key and value of userId
func (c *memoryStoreProvider) SetMfaSession(userId, key, otp string, expiration int64) error {
	currentTime := time.Now()
	expireTime := time.Unix(expiration, 0)
	duration := expireTime.Sub(currentTime)
	err := c.store.Set(c.ctx, fmt.Sprintf("%s%s:%s", mfaSessionPrefix, userId, key), otp, duration).Err()
	if err != nil {
		log.Debug().Err(err).Msg("Error saving user session to redis")
		return err
	}
	return nil
}

// GetMfaSession returns value of given mfa session
func (c *memoryStoreProvider) GetMfaSession(userId, key string) (string, error) {
	data, err := c.store.Get(c.ctx, fmt.Sprintf("%s%s:%s", mfaSessionPrefix, userId, key)).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

// DeleteMfaSession deletes given mfa session from in-memory store.
func (c *memoryStoreProvider) DeleteMfaSession(userId, key string) error {
	if err := c.store.Del(c.ctx, fmt.Sprintf("%s%s:%s", mfaSessionPrefix, userId, key)).Err(); err != nil {
		log.Debug().Err(err).Msg("Error deleting user session from redis")
		// continue
	}
	return nil
}
