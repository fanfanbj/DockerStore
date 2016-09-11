package registry

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	_"github.com/docker/distribution/manifest/schema1"
	_"github.com/docker/distribution/manifest/schema2"
)

func (registry *Registry) RegistryAPIGet(path, username, password string) ([]byte, string, error) {
	return registry.RegistryAPI("GET", path, username, password, "")
}

