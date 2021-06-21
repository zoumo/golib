package krb5

import (
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/pkg/errors"
)

func GetDefaultPrinciplaNameFromCCache(ccache string) (string, error) {
	if ccache == "" {
		return "", errors.New("invalid ccache file path")
	}
	c, err := credentials.LoadCCache(ccache)
	if err != nil {
		return "", errors.Wrap(err, "failed to load kerberos ccache")
	}
	if len(c.DefaultPrincipal.PrincipalName.NameString) == 0 {
		return "", errors.New("empty principalName")
	}
	return c.DefaultPrincipal.PrincipalName.NameString[0], nil
}
