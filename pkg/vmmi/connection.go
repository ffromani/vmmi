package vmmi

import (
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	libvirt "github.com/libvirt/libvirt-go"
	"io/ioutil"
)

func connect(opts *messages.Options) (*libvirt.Connect, error) {
	if opts.ConnectionCredentials.Username == "" || opts.ConnectionCredentials.PasswordFile == "" {
		return libvirt.NewConnect(opts.ConnectionURI)
	}

	auth := libvirt.ConnectAuth{
		CredType: []libvirt.ConnectCredentialType{
			libvirt.CRED_AUTHNAME,
			libvirt.CRED_PASSPHRASE,
		},
		Callback: func(creds []*libvirt.ConnectCredential) {
			for _, cred := range creds {
				if cred.Type == libvirt.CRED_AUTHNAME {
					cred.Result = opts.ConnectionCredentials.Username
					cred.ResultLen = len(cred.Result)
				} else if cred.Type == libvirt.CRED_PASSPHRASE {
					content, err := ioutil.ReadFile(opts.ConnectionCredentials.PasswordFile)
					if err == nil {
						cred.Result = string(content)
						cred.ResultLen = len(cred.Result)
					}
				}
			}
		},
	}
	return libvirt.NewConnectWithAuth(opts.ConnectionURI, &auth, 0)
}
