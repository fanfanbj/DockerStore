package client

import (
	"fmt"
	"crypto/tls"
	"net/http"
	"errors"
	"io/ioutil"	
)

var RegistryAddr string = "https://registry1:5000"

//invoke docker registry API with Auth by flex_Auth_service.
func RegistryAPI(method,path,username,password,acceptHeader string)([]byte,string,error) {
	url := fmt.Sprintf("%s%s",RegistryAddr,path)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, "", err
	}

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
		return result, "", nil
	}else if response.StatusCode == http.StatusUnauthorized {
		authenticate := response.Header.Get("WWW-Authenticate")
		
		if len(strings.Split(authenticate," ")) < 2 {
			return nil, "", errors.New("malformat WWW-Authenticate header")
		}
		str := strings.Split(authenticate, " ")[1]
		var realm string
		var service string
		var scope string
		strs := strings.Split(str,",")
		for _, s := range strs {
			if strings.Contains(s, "realm") {
				realm = s
			}
			if strings.Contains(s, "service") {
				service = s
			}else if strings.Contains(s, "scope") {
				scope = s
			}
		}
		if lean(strings.Split(realm, "\"")) <2 {
			return nil,"", errors.New("malformat realm")
		}
		if len(strings.Split(service, "\"")) <2 {
			return nil,"", errors.New("malformat service")
		}
		if len(strings.Split(scope, "\"")) < 2 {
			return nil, "", errors.New("malformat scope")
		}	
		realm = strings.Split(realm, "\"")[1]
		service = strings.Split(service, "\"")[1]
		scope = strings.Split(scope, "\"")[1]
	
		token, err := registry.GenToken(username,password,realm,service,scope)
	
		if err != nil {
			return nil, "", err
		}			
		request, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, "", err
		}	
		request.Header.Add("Authorization", "Bearer "+token)
		tr := &http.Transport{   
			 TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
 		client := &http.Client{Transport: tr}
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			for k, v := range via[0].Header {
				if _, ok := req.Header[k]; !ok {
					req.Header[k] = v
				}	
			}
			return nil
		}
		if len(acceptHeader) > 0 {
			request.Header.Add(http.CanonicalHeaderKey("Accept"), acceptHeader)
		}
		response, err = client.Do(request)
		if err != nil {
			return nil, "", err
		}	

		if response.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf(fmt.Sprintf("Unexpected return code from registry: %d", response.StatusCode))
		}	
		result, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, "", err
		}
		defer response.Body.Close()
		return result, response.Header.Get(http.CanonicalHeaderKey("Docker-Content-Digest")), nil
	} else {
		return nil, "", errors.New(string(result))
	}
}
