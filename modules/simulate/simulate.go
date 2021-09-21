package simulate

import "github.com/zincium/zinc/modules/auth"

type authSimulate struct {
}

func (as *authSimulate) DiscoverKey(fingerprint string) (*auth.SSHKeyResult, error) {

	return &auth.SSHKeyResult{
		Content: "OK",
		KID:     1,
	}, nil
}
func (as *authSimulate) Authorize(req *auth.AuthRequest) (*auth.AuthResult, error) {

	return &auth.AuthResult{
		Status:   0,
		Message:  "OK",
		RID:      1,
		UID:      1,
		Location: req.RelativePath,
		Address:  "127.0.0.1",
	}, nil
}

func NewAuthorizationClient() auth.AuthorizationClient {
	return &authSimulate{}
}
