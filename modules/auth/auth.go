package auth

import "net/http"

// authorization

type Action int

const (
	Download Action = iota
	Upload
	Archive
)

type Protocol int

const (
	SSH Protocol = iota
	HTTP
	Extension
)

type AuthRequest struct {
	RelativePath string
	UserName     string
	Password     string
	KID          int64
	Action       Action
	Protocol     Protocol
}

type AuthResult struct {
	Status   int
	Message  string
	RID      int64
	UID      int64
	Address  string
	Location string
}

func (ar *AuthResult) Success() bool {
	return ar.Status == 0 || ar.Status == 200
}

func (ar *AuthResult) RenderHTML(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(ar.Status)
	_, err := w.Write([]byte(ar.Message))
	return err
}

type SSHKeyResult struct {
	Content string
	KID     int64
}

type AuthorizationClient interface {
	// SSH SHA256 fingerprint
	DiscoverKey(fingerprint string) (*SSHKeyResult, error)
	// Auth
	Authorize(req *AuthRequest) (*AuthResult, error)
}
