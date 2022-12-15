package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/user"
)

var logger = log.Get()

type Guid struct {
	Token string
}

var (
	policiesToApply = []string{
		"639a488d6aff8e0001fc6aa6",
	}
)

// AddFooBarHeader adds custom "Foo: Bar" header to the request
func KeyCreation(rw http.ResponseWriter, res *http.Response, req *http.Request) {
	var g Guid
	b, err := ioutil.ReadAll(req.Body)
	logger.Info("JSON VALUE: ", b)
	json.Unmarshal([]byte(b), &g)
	// err := json.NewDecoder(r.Body).Decode(&g)
	// google how to unmarshall key value pairs in golang
	// need to construct struct & parse json into the struct

	// if token exists in temporary redis then create key and g.token not nil
	// else fail request
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	session := &user.SessionState{
		Alias: g.Token,
		MetaData: map[string]interface{}{
			"token": g.Token,
		},
		KeyID:         g.Token,
		ApplyPolicies: policiesToApply,
	}
	ctx.SetSession(req, session, true)

}

// Custom Auth, applies a rate limit of
// 2 per 10 given a token of "abc"
func AuthCheck(rw http.ResponseWriter, r *http.Request) {
	logger.Info("JSON VALUE: ")
	token := r.Header.Get("Authorization")
	if token != "d3fd1a57-94ce-4a36-9dfe-679a8f493b49" && token != "3be61aa4-2490-4637-93b9-105001aa88a5" {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	session := &user.SessionState{
		Alias: token,
		MetaData: map[string]interface{}{
			token: token,
		},
		KeyID: token,
	}
	ctx.SetSession(r, session, true)
}

func main() {}

func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")
}
