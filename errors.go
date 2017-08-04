package lanes

import "errors"

var (
	ErrMissingAccessKey  = errors.New("missing AWS access key ID")
	ErrMissingSecretKey  = errors.New("missing AWS secret key")
	ErrMissingSSHProfile = errors.New("missing SSH profile")
)
