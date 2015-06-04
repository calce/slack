// POST webhooks listener that does your chores
package slack

import (
	"bytes"
	"errors"
	"strconv"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

const basicAuthPrefix string = "Basic "

type Slack struct {
	host string
	port string
	username []byte
	password []byte
	cert string
	key string
	engines map[string]*Engine
	routes map[string]
}

func New(host string, port int, username string, password string, cert string, key string) {
	return Slack{
		host: host,
		port: strconv.Itoa(port),
		username: []byte(username),
		password: []byte(password),
		cert: cert,
		key: key,
		engines: make(map[string]*Engine),
	}
}

func (this *Slack) Register(engines ...*Engine) *Slack {
	for engine := range engines {
		this.engines[engine.GetName()] = engine
	}
	return this
}

func (this *Slack) isAuthenticated(res http.ResponseWriter, req *http.Request) bool {

	// Get the Basic Authentication credentials
	auth := req.Header.Get("Authorization")

	if strings.HasPrefix(auth, basicAuthPrefix) {

		// Check credentials
		payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)

			if len(pair) == 2 &&
				bytes.Equal(pair[0], this.username) &&
				bytes.Equal(pair[1], this.password) {
					return true
			}

		}
	}
	res.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return false	
}

func (this *Slack) handle(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if !this.isAuthenticated(res, req) { return }
	
	endpoint := ps.ByName("endpoint")
	if engine := this.engines[endpoint]; engine != nil {
		action := ps.ByName("action")
		engine.Do(action, req.PostForm)
	}
}

func (this *Slack) auth() httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		}
}

func (this *Slack) Do(name string, action string, params ...interfaces{}) (interface{}, error) {
	engine := this.engines[name]
	if engine == nil { return nil, errors.New("Engine not found: " + name)}
	return engine.Do(action, params)
}

func (this *Slack) Serve(tls bool) {

	router := httprouter.New()
	router.POST("/:endpoint/:action", this.auth())
	
	if tls {
		http.ListenAndServeTLS(this.host + ":" + this.port, this.cert, this.key, router)
	} else {
		http.ListenAndServe(this.host + ":" + this.port, router)
	}
	
}