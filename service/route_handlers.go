package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	"github.com/rancher/auth-proxy/server"
)

//CreateToken is a handler for route /token and returns the jwt token after authenticating the user
func CreateToken(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("GetToken failed with error: %v", err)
	}
	var jsonInput map[string]string

	err = json.Unmarshal(bytes, &jsonInput)
	if err != nil {
		log.Errorf("unmarshal failed with error: %v", err)
	}

	securityCode := jsonInput["code"]

	if securityCode != "" {
		log.Debugf("CreateToken called with securityCode")
		//getToken
		token, status, err := server.CreateToken(jsonInput)
		if err != nil {
			log.Errorf("CreateToken failed with error: %v", err)
			if status == 0 {
				status = http.StatusInternalServerError
			}
			ReturnHTTPError(w, r, status, fmt.Sprintf("%v", err))
			return
		}

		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(token)

	} else {
		ReturnHTTPError(w, r, http.StatusBadRequest, "Bad Request, Please check the request content")
		return
	}
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Info("Failed to get token cookie: %v", err)
		ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return
	}

	log.Infof("token cookie: %v %v", cookie.Name, cookie.Value)

	//getToken
	token, status, err := server.GetToken(cookie.Value)
	if err != nil {
		log.Errorf("GetToken failed with error: %v", err)
		if status == 0 {
			status = http.StatusInternalServerError
		}
		ReturnHTTPError(w, r, status, fmt.Sprintf("%v", err))
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(token)

}

func DeleteToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Info("Failed to get token cookie: %v", err)
		ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return
	}

	log.Infof("token cookie: %v %v", cookie.Name, cookie.Value)

	//getToken
	status, err := server.DeleteToken(cookie.Value)
	if err != nil {
		log.Errorf("DeleteToken failed with error: %v", err)
		if status == 0 {
			status = http.StatusInternalServerError
		}
		ReturnHTTPError(w, r, status, fmt.Sprintf("%v", err))
		return
	}
}

func GetIdentities(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Info("Failed to get token cookie: %v", err)
		ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return
	}

	log.Infof("token cookie: %v %v", cookie.Name, cookie.Value)

	//getToken
	identities, status, err := server.GetIdentities(cookie.Value)
	if err != nil {
		log.Errorf("DeleteToken failed with error: %v", err)
		if status == 0 {
			status = http.StatusInternalServerError
		}
		ReturnHTTPError(w, r, status, fmt.Sprintf("%v", err))
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(identities)
}
