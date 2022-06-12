package tinyCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_cache"

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	self     string
	basePath string
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path + "\n")
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// /<basePath>/<groupName>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath)+1:], "/", 2)
	// parts := strings.SplitN(r.URL.Path, "/", 2)

	fmt.Println()
	fmt.Println(parts)
	fmt.Println()
	// for debug

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(view.ByteSlice())
	if err != nil {
		return
	}
}
