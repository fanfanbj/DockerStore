package client

import (
	"fmt"
	"crypto/tls"
	"net/http"
	"errors"
	"io/ioutil"	
)

var RegistryAddr string = "https://registry1:5000"

//invoke docker registry API by httpclient.
func RegistryAPI(method,path,acceptHeader string)([]byte,string,error) {
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
	}else{
		return nil, "", errors.New(string(result))
	}
}
