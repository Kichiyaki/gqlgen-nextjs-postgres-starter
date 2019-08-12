package sessions

import "github.com/gorilla/sessions"

type Serializer interface {
	Deserialize(d []byte, ss *sessions.Session) error
	Serialize(ss *sessions.Session) ([]byte, error)
}
