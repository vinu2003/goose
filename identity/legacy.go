package identity

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Legacy struct {
	client *http.Client
}

func (l *Legacy) Auth(creds *Credentials) (*AuthDetails, error) {
	if l.client == nil {
		l.client = &http.Client{CheckRedirect: nil}
	}
	request, err := http.NewRequest("GET", creds.URL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Auth-User", creds.User)
	request.Header.Set("X-Auth-Key", creds.Secrets)
	response, err := l.client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusNoContent {
		content, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("Failed to Authenticate (code %d %s): %s",
			response.StatusCode, response.Status, content)
	}
	details := &AuthDetails{}
	details.TokenId = response.Header.Get("X-Auth-Token")
	if details.TokenId == "" {
		return nil, fmt.Errorf("Did not get valid Token from auth request")
	}
	nova_url := response.Header.Get("X-Server-Management-Url")
	if nova_url == "" {
		return nil, fmt.Errorf("Did not get valid management URL from auth request")
	}
	details.ServiceURLs = map[string]string{"compute": nova_url}

	return details, nil
}
