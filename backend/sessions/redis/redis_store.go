package redis

import (
	"encoding/base32"
	"errors"
	"net/http"
	"strings"

	ginSessions "github.com/gin-contrib/sessions"
	goRedis "github.com/go-redis/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_sessions "github.com/kichiyaki/graphql-starter/backend/sessions"
	"github.com/kichiyaki/graphql-starter/backend/sessions/serializer"
)

const (
	defaultMaxAge = 15 * 60
)

type RedisStore struct {
	conn       *goRedis.Client
	codecs     []securecookie.Codec
	options    *sessions.Options
	maxLength  int
	keyPrefix  string
	serializer _sessions.Serializer
}

func NewRedisStore(conn *goRedis.Client, keyPairs ...[]byte) *RedisStore {
	return &RedisStore{
		conn:   conn,
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		options: &sessions.Options{
			Path:   "/",
			MaxAge: defaultMaxAge,
		},
		maxLength:  4096,
		keyPrefix:  "session_",
		serializer: serializer.NewJSONSerializer(),
	}
}

func (s *RedisStore) SetKeyPrefix(p string) *RedisStore {
	s.keyPrefix = p
	return s
}

func (s *RedisStore) SetOptions(opts *sessions.Options) *RedisStore {
	s.options = opts
	return s
}

func (s *RedisStore) SetMaxLength(length int) *RedisStore {
	s.maxLength = length
	return s
}

func (s *RedisStore) SetSerializer(serializer _sessions.Serializer) *RedisStore {
	s.serializer = serializer
	return s
}

func (s *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *RedisStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.options
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.codecs...)
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

func (s *RedisStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.conn.Do("DEL", s.keyPrefix+session.ID).Err(); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

func (s *RedisStore) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if err := s.conn.Del(s.keyPrefix + session.ID).Err(); err != nil {
		return err
	}
	// Set cookie to expire.
	options := *session.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))
	// Clear session values.
	for k := range session.Values {
		delete(session.Values, k)
	}
	return nil
}

func (s *RedisStore) DeleteByID(ids ...string) error {
	formattedIds := []string{}
	for _, id := range ids {
		if !strings.Contains(id, s.keyPrefix) {
			formattedIds = append(formattedIds, s.keyPrefix+id)
		} else {
			formattedIds = append(formattedIds, id)
		}
	}
	return s.conn.Del(formattedIds...).Err()
}

func (s *RedisStore) GetAll() ([]*sessions.Session, error) {
	keys, _, err := s.conn.Scan(0, "", 0).Result()
	if err != nil {
		return nil, err
	}
	results := []*sessions.Session{}
	for _, key := range keys {
		val, err := s.conn.Get(key).Result()
		if err != nil {
			return nil, err
		}
		sess := &sessions.Session{
			ID:     strings.Replace(key, s.keyPrefix, "", -1),
			Values: make(map[interface{}]interface{}),
		}
		err = s.serializer.Deserialize([]byte(val), sess)
		results = append(results, sess)
	}
	return results, nil
}

func (c *RedisStore) Options(options ginSessions.Options) {
	c.options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (s *RedisStore) load(session *sessions.Session) (bool, error) {
	data, err := s.conn.Get(s.keyPrefix + session.ID).Result()
	if err != nil {
		return false, err
	}
	if data == "" {
		return false, nil // no data was associated with this key
	}
	return true, s.serializer.Deserialize([]byte(data), session)
}

func (s *RedisStore) save(session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	age := session.Options.MaxAge
	if age == 0 {
		age = defaultMaxAge
	}
	return s.conn.Do("SETEX", s.keyPrefix+session.ID, age, b).Err()
}
