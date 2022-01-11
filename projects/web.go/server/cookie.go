package server

import (
	"encoding/base64"
	"net/http"

	"github.com/gorilla/securecookie"
)

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    base64.URLEncoding.EncodeToString([]byte(value)),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}

	http.SetCookie(c.Writer, cookie)
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {

	var err error

	if cookie, err := c.Req.Cookie(name); err == nil {
		val, err := base64.URLEncoding.DecodeString(cookie.Value)
		return string(val), err
	}
	return "", err
}

func NewSecureCookie(hashkey, blockkey string) *securecookie.SecureCookie {
	if len(hashkey) < 16 || len(blockkey) < 16 {
		return nil
	}
	// Hash keys should be at least 32 bytes long
	var hashKey  = []byte(hashkey)
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte(blockkey)
	var s2 = securecookie.New(hashKey, blockKey)

	return s2
}

// SetSecureCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetSecureCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, s2 *securecookie.SecureCookie) {
	if path == "" {
		path = "/"
	}
	emValue := map[string]interface{}{
		"oob": "rab",
		"xxx": value,
	}

	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.MaxAge = maxAge
	cookie.Path = path
	cookie.Domain = domain
	cookie.Secure = secure
	cookie.HttpOnly = httpOnly

	if encoded, err := s2.Encode(name, emValue); err == nil {
		cookie.Value = encoded
	} else {
		cookie.Value = base64.URLEncoding.EncodeToString([]byte(value))
	}
	http.SetCookie(c.Writer, cookie)
}

// SecureCookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) SecureCookie(name string, s2 *securecookie.SecureCookie) (string, error) {

	var err error
	if cookie, err := c.Req.Cookie(name); err == nil {
		value := make(map[string]interface{})
		if err = s2.Decode(name, cookie.Value, &value); err == nil {
			if value["oob"] == "rab" {
				return value["xxx"].(string), nil
			}
		} else {
			val, err := base64.URLEncoding.DecodeString(cookie.Value)
			return string(val), err
		}
	}
	return "", err
}

func (c *Context) DelCookie(name, path string) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		MaxAge:   -1,
		Path:     path,
	})
}
