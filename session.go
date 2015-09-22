package happyngine

import (
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"net/http"
)

var store *CookieStore = nil

func init() {

	if env.Get("SESSION_AUTHENTICATION_KEY") == "" {
		log.Debugln("No SESSION_AUTHENTICATION_KEY, ignoring session")
		return
	}

	store = NewCookieStore([]byte(env.Get("SESSION_AUTHENTICATION_KEY")),
		[]byte(env.Get("SESSION_ENCRYPTION_KEY")),
		[]byte(env.Get("SESSION_AUTHENTICATION_KEY_OLD")),
		[]byte(env.Get("SESSION_ENCRYPTION_KEY_OLD")))
}

// Options stores configuration for a session or session store.
//
// Fields are a subset of http.Cookie fields.
type SessionOptions struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

type Session struct {
	ID      string
	Values  map[string]interface{}
	Options *SessionOptions
	IsNew   bool
	name    string
	changed bool
}

// Save is a convenience method to save this session. It is the same as calling
// store.Save(request, response, session). You should call Save before writing to
// the response or returning from the handler.
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {

	return store.Save(r, w, s)
}

// Name returns the name used to register the session.
func (s *Session) Name() string {
	return s.name
}

func (s *Session) Get(name string) interface{} {

	i, ok := s.Values[name]
	if !ok {
		return nil
	}

	return i
}

func (s *Session) Set(name string, value interface{}) {

	s.changed = true
	s.Values[name] = value
}

func (s *Session) SetMap(m map[string]interface{}) {

	s.changed = true

	for key, value := range m {
		s.Values[key] = value
	}
}

func (s *Session) Del(name string) {

	s.changed = true
	delete(s.Values, name)
}
func (s *Session) Destroy() {

	s.changed = true
	s.Options.MaxAge = -1
}

func (s *Session) Changed() bool {

	return s.changed
}

func GetSession(r *http.Request, name string) *Session {

	s, _ := store.Get(r, name)

	if s == nil || s.IsNew {
		return nil
	}

	s.changed = false

	return s
}

func NewSession(name string, options *SessionOptions) *Session {

	return &Session{
		Values:  make(map[string]interface{}),
		Options: options,
		name:    name,
		changed: true,
	}
}
