package credential

import (
	"context"
	"fmt"

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
	PluginCredentialType         = "plugin"
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

func PutServiceUserCredential(ctx context.Context, svcEntity, userIdentifer, userName, password string) error {
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing credential admin")
	}

	err = credAdmin.PutServiceUserCredential(context.Background(), svcEntity, userIdentifer,
		credentials.ServiceCredential{
			UserName: userName,
			Password: password,
		})
	if err != nil {
		return errors.WithMessagef(err, "error in put service cred for %s/%s", svcEntity, userIdentifer)
	}
	return nil
}

func PutClusterCerts(ctx context.Context, orgID, clusterName, clientCAChainData, clientKeyData, clientCertData string) error {
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing credential admin")
	}

	certIndetifier := getClusterCertIndentifier(orgID, clusterName)
	err = credAdmin.PutCertificateData(context.Background(), clusterCertEntity, certIndetifier,
		credentials.CertificateData{
			CACert: clientCAChainData,
			Key:    clientKeyData,
			Cert:   clientCertData,
		})
	if err != nil {
		return errors.WithMessagef(err, "error in put cert for %s/%s", clusterCertEntity, certIndetifier)
	}
	return nil
}

func PutGenericCredential(ctx context.Context, svcEntity, credId string, cred map[string]string) error {
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing credential admin")
	}

	err = credAdmin.PutCredential(context.Background(), credentials.GenericCredentialType,
		svcEntity, credId, cred)
	if err != nil {
		return errors.WithMessagef(err, "error in put generic cred for %s/%s", svcEntity, credId)
	}
	return nil
}

func getClusterCertIndentifier(orgID, clusterName string) string {
	return fmt.Sprintf("%s:%s", orgID, clusterName)
}

func PutPluginCredential(ctx context.Context, pluginName, svcEntity string, cred map[string]string) error {
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing credential admin")
	}

	err = credAdmin.PutCredential(context.Background(), PluginCredentialType,
		pluginName, svcEntity, cred)
	if err != nil {
		return errors.WithMessagef(err, "error in put generic cred for %s/%s", pluginName, svcEntity)
	}
	return nil
}

func GetPluginCredential(ctx context.Context, pluginName, svcEntity string) (cred credentials.ServiceCredential, err error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return
	}

	data, err := credReader.GetCredential(ctx, PluginCredentialType, pluginName, svcEntity)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading cluster global values %s/%s from the vault",
			captenConfigEntityName, globalValuesCredIdentifier)
		return
	}
	return credentials.ParseServiceCredential(data)
}
