package lanes

import "errors"

var (
	ErrAbort             = errors.New("aborted")
	ErrMissingAWSProfile = errors.New("missing AWS profile")
	ErrMissingAccessKey  = errors.New("missing AWS access key ID")
	ErrMissingSecretKey  = errors.New("missing AWS secret key")
	ErrMissingSSHProfile = errors.New("missing SSH profile")
)
