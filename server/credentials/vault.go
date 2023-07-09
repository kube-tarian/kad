package credentials

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	vaultauth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type config struct {
	Address                string        `envconfig:"VAULT_ADDR" required:"true"`
	CACert                 string        `envconfig:"VAULT_CACERT" required:"false"`
	ReadTimeout            time.Duration `envconfig:"VAULT_READ_TIMEOUT" default:"60s"`
	MaxRetries             int           `envconfig:"VAULT_MAX_RETRIES" default:"5"`
	VaultKVMountPath       string        `envconfig:"VAULT_KV_MOUNT_PATH" default:"secret"`
	VaultToken             string        `envconfig:"VAULT_TOKEN"`
	VaultRole              string        `envconfig:"VAULT_ROLE" required:"true"`
	ServiceAccoutTokenPath string        `envconfig:"SERVICE_ACCOUNT_TOKEN_PATH" default:"/var/run/secrets/kubernetes.io/serviceaccount/token"`
}

type client struct {
	c    *api.Client
	conf config
}

func newClientWithAuth(ctx context.Context) (*client, error) {
	conf := config{}
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, err
	}

	vc, err := newClient(conf)
	if err != nil {
		return nil, err
	}

	if len(conf.VaultToken) != 0 {
		vc.c.SetToken(conf.VaultToken)
		return vc, nil
	}

	err = vc.configureAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

func newClient(conf config) (*client, error) {
	cfg, err := prepareVaultConfig(conf)
	if err != nil {
		return nil, err
	}

	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &client{
		c:    c,
		conf: conf,
	}, nil
}

func prepareVaultConfig(conf config) (cfg *api.Config, err error) {
	cfg = api.DefaultConfig()
	cfg.Address = conf.Address
	cfg.Timeout = conf.ReadTimeout
	cfg.Backoff = retryablehttp.DefaultBackoff
	cfg.MaxRetries = conf.MaxRetries
	if conf.CACert != "" {
		tlsConfig := api.TLSConfig{CACert: conf.CACert}
		err = cfg.ConfigureTLS(&tlsConfig)
	}
	return
}

func prepareCredentialSecretPath(credentialType, credEntityName, credIdentifier string) string {
	return fmt.Sprintf("%s/%s/%s", credentialType, credEntityName, credIdentifier)
}

func readFileContent(path string) (s string, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	s = string(b)
	return
}

func (vc *client) configureAuthToken(ctx context.Context) (err error) {
	serviceToken, err := readFileContent(vc.conf.ServiceAccoutTokenPath)
	if err != nil {
		return err
	}

	k8sAuth, err := vaultauth.NewKubernetesAuth(
		vc.conf.VaultRole,
		vaultauth.WithServiceAccountToken(serviceToken),
	)
	if err != nil {
		return errors.WithMessagef(err, "error in initializing Kubernetes auth method")
	}

	authInfo, err := vc.c.Auth().Login(ctx, k8sAuth)
	if err != nil {
		return errors.WithMessagef(err, "error in login with Kubernetes auth")
	}
	if authInfo == nil {
		return errors.New("no auth info was returned after login")
	}
	return nil
}

func (vc *client) getCredential(ctx context.Context, secretPath string) (cred map[string]string, err error) {
	secretValByPath, err := vc.c.KVv2(vc.conf.VaultKVMountPath).Get(context.Background(), secretPath)
	if err != nil {
		err = errors.WithMessagef(err, "error in reading certificate data from %s", secretPath)
		return
	}

	if secretValByPath == nil {
		err = errors.WithMessagef(err, "crdentaial not found at %s", secretPath)
		return
	}
	if secretValByPath.Data == nil {
		err = errors.WithMessagef(err, "crdentaial data is corrupted for %s", secretPath)
		return
	}
	cred = map[string]string{}
	for key, val := range secretValByPath.Data {
		cred[key] = val.(string)
	}
	return
}

func (vc *client) putCredential(ctx context.Context, secretPath string, cred map[string]string) (err error) {
	credData := map[string]interface{}{}
	for key, val := range cred {
		credData[key] = val
	}
	_, err = vc.c.KVv2(vc.conf.VaultKVMountPath).Put(ctx, secretPath, credData)
	if err != nil {
		err = errors.WithMessagef(err, "error in putting credentail at %s", secretPath)
	}
	return
}

func (vc *client) deleteCredential(ctx context.Context, secretPath string) (err error) {
	err = vc.c.KVv2(vc.conf.VaultKVMountPath).Delete(ctx, secretPath)
	if err != nil {
		err = errors.WithMessagef(err, "error in deleting credentail at %s", secretPath)
	}
	return
}
