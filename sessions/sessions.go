package sessions

import (
	"github.com/gorilla/sessions"
)

//byte array for keys
var Store = sessions.NewCookieStore([]byte("secret-lol"))
