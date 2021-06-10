package krb5

import (
	"fmt"

	"github.com/jcmturner/gokrb5/v8/credentials"
)

func GetDefaultPrinciplaNameFromCCache(ccache string) (string, error) {
	if ccache == "" {
		return "", fmt.Errorf("invalid ccache file path")
	}
	c, err := credentials.LoadCCache(ccache)
	if err != nil {
		return "", fmt.Errorf("failed to load kerberos ccache: %v", err)
	}
	if len(c.DefaultPrincipal.PrincipalName.NameString) == 0 {
		return "", fmt.Errorf("empty principalName")
	}
	return c.DefaultPrincipal.PrincipalName.NameString[0], nil
}
