package krb5

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/crypto"
	"github.com/jcmturner/gokrb5/v8/gssapi"
	"github.com/jcmturner/gokrb5/v8/iana/chksumtype"
	"github.com/jcmturner/gokrb5/v8/iana/flags"
	"github.com/jcmturner/gokrb5/v8/messages"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"github.com/jcmturner/gokrb5/v8/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

const (
	ContextFlagREADY = 128
)

func NewGSSAPIClientWithCCache(krb5confPath, ccachePath string) (ssh.GSSAPIClient, error) {
	cfg, err := config.Load(krb5confPath)
	if err != nil {
		return nil, err
	}
	cache, err := credentials.LoadCCache(ccachePath)
	if err != nil {
		return nil, err
	}
	client, err := client.NewFromCCache(cache, cfg)
	if err != nil {
		return nil, err
	}
	return newGSSAPIClient(client)
}

func NewGSSAPIClientWithPassword(krb5confPath, username, realm, password string) (ssh.GSSAPIClient, error) {
	cfg, err := config.Load(krb5confPath)
	if err != nil {
		return nil, err
	}
	client := client.NewWithPassword(username, realm, password, cfg)
	return newGSSAPIClient(client)
}

func newGSSAPIClient(client *client.Client) (ssh.GSSAPIClient, error) {
	err := client.Login()
	if err != nil {
		return nil, errors.Wrap(err, "failed to login kerberos")
	}
	err = client.AffirmLogin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to affirmLogin kerberos")
	}

	return &GSSAPIClient{
		client: client,
	}, nil
}

type GSSAPIClient struct {
	client *client.Client
	subkey types.EncryptionKey
}

// Create new authenticator checksum for kerberos MechToken
func (k *GSSAPIClient) newAuthenticatorChksum(flags []int) []byte {
	a := make([]byte, 24)
	binary.LittleEndian.PutUint32(a[:4], 16)
	for _, i := range flags {
		if i == gssapi.ContextFlagDeleg {
			x := make([]byte, 28-len(a))
			a = append(a, x...)
		}
		f := binary.LittleEndian.Uint32(a[20:24])
		f |= uint32(i)
		binary.LittleEndian.PutUint32(a[20:24], f)
	}
	return a
}

func (k *GSSAPIClient) InitSecContext(target string, token []byte, isGSSDelegCreds bool) ([]byte, bool, error) {
	GSSAPIFlags := []int{
		ContextFlagREADY,
		gssapi.ContextFlagInteg,
		gssapi.ContextFlagMutual,
	}
	if isGSSDelegCreds {
		GSSAPIFlags = append(GSSAPIFlags, gssapi.ContextFlagDeleg)
	}
	APOptions := []int{flags.APOptionMutualRequired}

	if len(token) == 0 {
		newTarget := strings.ReplaceAll(target, "@", "/")

		tkt, sKey, err := k.client.GetServiceTicket(newTarget)
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to get service ticket")
		}

		krb5Token, err := spnego.NewKRB5TokenAPREQ(k.client, tkt, sKey, GSSAPIFlags, APOptions) //nolint

		creds := k.client.Credentials
		auth, err := types.NewAuthenticator(creds.Domain(), creds.CName())
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to generate new authenticator")
		}
		auth.Cksum = types.Checksum{
			CksumType: chksumtype.GSSAPI,
			Checksum:  k.newAuthenticatorChksum(GSSAPIFlags),
		}
		etype, _ := crypto.GetEtype(sKey.KeyType)
		if err := auth.GenerateSeqNumberAndSubKey(sKey.KeyType, etype.GetKeyByteSize()); err != nil {
			return nil, false, errors.Wrap(err, "failed to generate seq number and sub key")
		}
		k.subkey = auth.SubKey

		apReq, err := messages.NewAPReq(tkt, sKey, auth)
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to create NewAPReq")
		}
		for _, o := range APOptions {
			types.SetFlag(&apReq.APOptions, o)
		}
		krb5Token.APReq = apReq

		outToken, err := krb5Token.Marshal()
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to marshal krb5 token")
		}
		return outToken, true, nil
	}

	var krb5Token spnego.KRB5Token
	if err := krb5Token.Unmarshal(token); err != nil {
		err := errors.Wrap(err, "unmarshal APRep token failed")
		return nil, false, err
	}
	if krb5Token.IsKRBError() {
		return nil, false, fmt.Errorf("received Kerberos error")
	}
	return nil, false, nil
}

func (k *GSSAPIClient) GetMIC(micFiled []byte) ([]byte, error) {
	micToken, err := gssapi.NewInitiatorMICToken(micFiled, k.subkey)
	if err != nil {
		return nil, err
	}
	token, err := micToken.Marshal()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (k *GSSAPIClient) DeleteSecContext() error {
	// k.client.Destroy()
	return nil
}
