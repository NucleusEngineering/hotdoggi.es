package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	auth "firebase.google.com/go/v4/auth"
	gin "github.com/gin-gonic/gin"
)

// Authenticate implements the security middleware
func Authenticate(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("trace.context", ctx)

	// Skip verification in non-prod
	if Global["environment"].(string) != "prod" {
		c.Next()
		return
	}

	// Trace verification
	traceToken := c.GetHeader("X-Cloud-Trace-Context")
	if traceToken != "" {
		// TODO do something with the trace context
		c.Set("trace.id", traceToken)
	}

	// Protocol verification
	protoHeader := c.GetHeader("X-Forwarded-Proto")
	if protoHeader == "" {
		Respond(c, http.StatusUnauthorized, "missing protocol header")
		return
	}
	if protoHeader != "https" {
		Respond(c, http.StatusUnauthorized, "refusing to serve unencrypted traffic")
		return
	}

	// Gateway/Proxy verification
	gatewayHeader := c.GetHeader("Authorization")
	if gatewayHeader == "" {
		Respond(c, http.StatusUnauthorized, "missing gateway authorization header")
		return
	}
	gatewayToken := strings.Split(gatewayHeader, " ")[1]
	clientHTTP := Global["client.http"].(*http.Client)
	response, err := clientHTTP.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", gatewayToken))
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to reach gateway verification endpoint")
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to read verification response")
		return
	}
	var responseMap map[string]interface{}
	json.Unmarshal(body, &responseMap)
	if _, ok := responseMap["error"]; ok {
		Respond(c, http.StatusUnauthorized, "failed to verify gateway token")
		return
	}
	if responseMap["iss"] != "https://accounts.google.com" {
		Respond(c, http.StatusUnauthorized, "failed to verify gateway token issuer")
		return
	}
	gatewayIdentity := Global["gateway.service_account"].(string)
	if responseMap["email"] != gatewayIdentity {
		Respond(c, http.StatusUnauthorized, "failed to verify gateway token identity")
		return
	}

	// Original user authorization verification
	firebaseHeader := c.GetHeader("X-Forwarded-Authorization")
	if firebaseHeader == "" {
		Respond(c, http.StatusUnauthorized, "missing authorization user authorization header")
		return
	}
	firebaseToken := strings.Split(firebaseHeader, " ")[1]
	clientFirebase := Global["client.firebase"].(*auth.Client)
	token, err := clientFirebase.VerifyIDToken(ctx, firebaseToken)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "unable to verify user token")
		return
	}

	// Gateway/Proxy user pre-flight authorization verification
	encoded := c.Request.Header.Get("X-Endpoint-API-UserInfo")
	if encoded == "" {
		Respond(c, http.StatusUnauthorized, "missing gateway user info header")
		return
	}
	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to decode user info header")
		return
	}
	var principal Principal
	err = json.Unmarshal(bytes, &principal)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to deserialize user info header")
		return
	}
	if token.UID != principal.ID {
		Respond(c, http.StatusUnauthorized, "mismatching inbound principal identities")
		return
	}

	// Verification OK

	c.Set("principal.email", principal.Email)
	c.Set("principal.id", principal.ID)
	c.Set("principal.name", principal.Name)
	c.Set("principal.picture", principal.PictureURL)
	c.Next()
}
