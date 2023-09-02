package credential

import (
	"context"

	"github.com/intelops/go-common/credentials"
	"github.com/pkg/errors"
)

const (
	clusterCertEntity            = "client-cert"
	serviceClientOAuthEntityName = "service-client-oauth"
	oauthClientIdKey             = "CLIENT_ID"
	oauthClientSecretKey         = "CLIENT_SECRET"
	captenConfigEntityName       = "capten-config"
	globalValuesCredIdentifier   = "global-values"
)

func GetServiceUserCredential(ctx context.Context, svcEntity, userName string) (cred credentials.ServiceCredential, err error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return
	}
	cred, err = credReader.GetServiceUserCredential(context.Background(), svcEntity, userName)
	if err != nil {
		err = errors.WithMessagef(err, "error in reading credential for %s/%s", svcEntity, userName)
	}
	return
}

func GetClusterCerts(ctx context.Context, clusterID string) (cred credentials.CertificateData, err error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return
	}

	cred, err = credReader.GetCertificateData(context.Background(), clusterCertEntity, clusterID)
	if err != nil {
		err = errors.WithMessagef(err, "error in reading cert for %s/%s", clusterCertEntity, clusterID)
	}
	return
}

func GetGenericCredential(ctx context.Context, entityName, credIndentifer string) (cred map[string]string, err error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return
	}

	cred, err = credReader.GetCredential(context.Background(), credentials.GenericCredentialType, entityName, credIndentifer)
	if err != nil {
		err = errors.WithMessagef(err, "error in reading cred for %s/%s", clusterCertEntity, credIndentifer)
	}
	return
}

func PutClusterCerts(ctx context.Context, clusterID, clientCAChainData, clientKeyData, clientCertData string) error {
	credReader, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing credential admin")
	}

	err = credReader.PutCertificateData(context.Background(), clusterCertEntity, clusterID,
		credentials.CertificateData{
			CACert: clientCAChainData,
			Key:    clientKeyData,
			Cert:   clientCertData,
		})
	if err != nil {
		return errors.WithMessagef(err, "error in put cert for %s/%s", clusterCertEntity, clusterID)
	}
	return nil
}

func DeleteClusterCerts(ctx context.Context, clusterID string) (err error) {
	credReader, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential admin")
		return
	}

	err = credReader.DeleteCertificateData(context.Background(), clusterCertEntity, clusterID)
	if err != nil {
		err = errors.WithMessagef(err, "error in delete cert for %s/%s", clusterCertEntity, clusterID)
	}
	return
}

func StoreAppOauthCredential(ctx context.Context, serviceName, clientId, clientSecret string) error {
	credWriter, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing credential admin")
	}

	cred := map[string]string{
		oauthClientIdKey:     clientId,
		oauthClientSecretKey: clientSecret,
	}

	err = credWriter.PutCredential(ctx, credentials.GenericCredentialType,
		serviceClientOAuthEntityName, serviceName, cred)
	if err != nil {
		return errors.WithMessagef(err, "error while storing service oauth credential %s/%s into the vault",
			serviceClientOAuthEntityName, serviceName)
	}
	return nil
}

func GetAppOauthCredential(ctx context.Context, serviceName string) (clientId, clientSecret string, err error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		serviceClientOAuthEntityName, serviceName)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading service oauth credential %s/%s from the vault",
			serviceClientOAuthEntityName, serviceName)
		return
	}

	clientId = cred[oauthClientIdKey]
	clientSecret = cred[oauthClientSecretKey]
	if len(clientId) == 0 || len(clientSecret) == 0 {
		err = errors.WithMessagef(err, "invalid service oauth credential %s/%s in the vault",
			serviceClientOAuthEntityName, serviceName)
		return
	}
	return
}

func GetClusterGlobalValues(ctx context.Context) (globalValues string, err error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		captenConfigEntityName, globalValuesCredIdentifier)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading cluster global values %s/%s from the vault",
			captenConfigEntityName, globalValuesCredIdentifier)
		return
	}

	globalValues = cred[globalValuesCredIdentifier]
	if len(globalValues) == 0 {
		err = errors.WithMessagef(err, "invalid cluster global values %s/%s in the vault",
			captenConfigEntityName, globalValuesCredIdentifier)
		return
	}
	return
}
