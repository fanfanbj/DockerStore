package registry

import (
	"fmt"
        "crypto/tls"
        "net/http"
        "errors"
        "io/ioutil"
	"encoding/base64"
	"encoding/json"
)

type Token struct {
	Token string `json:"token"`
}
}

//retrieve JWT token by invoking flex_auth_service.
func (registry *Registry) GenToken(username,password,realm,service,scope string) (string,error) {
	url := fmt.Sprintf("realm=%s&service=%s&scope=%s",realm,service,scope)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
 	       return nil, "", err
        }

	request.Header.Add("Authorization", "Basic "+base64UrlEncode([]byte(fmt.Sprintf("%s:%s",username,password))

	tr := &http.Transport{
        	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
 
        client := &http.Client{Transport: tr}
        response, err := client.Do(request)
        if err != nil {
        	return nil, "", err
        }
	result, err := ioutil.ReadAll(response.Body)
          	if err != nil {
                	return nil, "", err
        }
        defer response.Body.Close()
        if response.StatusCode == http.StatusOK {
		json.Unmarshal([]byte(result), &Token)
 		return Token.Token, "", nil
	}else {
		return nil, "", errors.New(string(result))
	}
	
}

func base64UrlEncode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
} 
