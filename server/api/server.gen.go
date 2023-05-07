// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

// AgentRequest defines model for AgentRequest.
type AgentRequest struct {
	CaCrt     []openapi_types.File `json:"ca_crt"`
	ClientCrt []openapi_types.File `json:"client_crt"`
	ClientKey []openapi_types.File `json:"client_key"`
}

// AgentResponse defines model for AgentResponse.
type AgentResponse struct {
	CustomerId string `json:"customer_id"`
	Endpoint   string `json:"endpoint"`
}

// ClimonDeleteRequest defines model for ClimonDeleteRequest.
type ClimonDeleteRequest struct {
	// ClusterName Cluster in which to be deleted, default in-build cluster
	ClusterName *string `json:"cluster_name,omitempty"`
	CustomerId  string  `json:"customer_id"`

	// Namespace Namespace chart to be installed
	Namespace string `json:"namespace"`

	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ReleaseName Release name to be used for install
	ReleaseName string `json:"release_name"`

	// Timeout Timeout for the application installation
	Timeout int `json:"timeout"`
}

// ClimonPostRequest defines model for ClimonPostRequest.
type ClimonPostRequest struct {
	// ChartName Chart name in Repository
	ChartName string `json:"chart_name"`

	// ClusterName Cluster in which to be installed, default in-build cluster
	ClusterName *string `json:"cluster_name,omitempty"`

	// Namespace Namespace chart to be installed
	Namespace string `json:"namespace"`

	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ReleaseName Release name to be used for install
	ReleaseName string `json:"release_name"`

	// RepoName Repository name
	RepoName string `json:"repo_name"`

	// RepoUrl Repository URL
	RepoUrl string `json:"repo_url"`

	// Timeout Timeout for the application installation
	Timeout int `json:"timeout"`

	// Version Version of the chart
	Version *string `json:"version,omitempty"`
}

// ClusterRequest defines model for ClusterRequest.
type ClusterRequest struct {
	ClusterName string `json:"cluster_name"`
	CustomerId  string `json:"customer_id"`
	PluginName  string `json:"plugin_name"`
}

// DeployerDeleteRequest defines model for DeployerDeleteRequest.
type DeployerDeleteRequest struct {
	// ClusterName Cluster in which to be deleted, default in-build cluster
	ClusterName *string `json:"cluster_name,omitempty"`

	// Namespace Namespace chart to be installed
	Namespace string `json:"namespace"`

	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ReleaseName Release name to be used for install
	ReleaseName string `json:"release_name"`

	// Timeout Timeout for the application installation
	Timeout int `json:"timeout"`
}

// DeployerPostRequest defines model for DeployerPostRequest.
type DeployerPostRequest struct {
	// ChartName Chart name in Repository
	ChartName string `json:"chart_name"`

	// ClusterName Cluster in which to be installed, default in-build cluster
	ClusterName *string `json:"cluster_name,omitempty"`

	// Namespace Namespace chart to be installed
	Namespace string `json:"namespace"`

	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ReleaseName Release name to be used for install
	ReleaseName string `json:"release_name"`

	// RepoName Repository name
	RepoName string `json:"repo_name"`

	// RepoUrl Repository URL
	RepoUrl string `json:"repo_url"`

	// Timeout Timeout for the application installation
	Timeout int `json:"timeout"`

	// Version Version of the chart
	Version *string `json:"version,omitempty"`
}

// ProjectDeleteRequest defines model for ProjectDeleteRequest.
type ProjectDeleteRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ProjectName Project name to be created in plugin
	ProjectName string `json:"project_name"`
}

// ProjectPostRequest defines model for ProjectPostRequest.
type ProjectPostRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ProjectName Project name to be created in plugin
	ProjectName string `json:"project_name"`
}

// RepositoryDeleteRequest defines model for RepositoryDeleteRequest.
type RepositoryDeleteRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// RepoName Repository to added to plugin
	RepoName string `json:"repo_name"`
}

// RepositoryPostRequest defines model for RepositoryPostRequest.
type RepositoryPostRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// RepoName Repository to added to plugin
	RepoName string `json:"repo_name"`

	// RepoUrl Repository URL
	RepoUrl string `json:"repo_url"`
}

// Response Configuration request response
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// StoreCredRequest defines model for StoreCredRequest.
type StoreCredRequest struct {
	Credname   *string `json:"credname,omitempty"`
	CustomerId *string `json:"customer_id,omitempty"`
	Password   *string `json:"password,omitempty"`
	Username   *string `json:"username,omitempty"`
}

// DeleteClimonJSONRequestBody defines body for DeleteClimon for application/json ContentType.
type DeleteClimonJSONRequestBody = ClimonDeleteRequest

// PostClimonJSONRequestBody defines body for PostClimon for application/json ContentType.
type PostClimonJSONRequestBody = ClimonPostRequest

// PutClimonJSONRequestBody defines body for PutClimon for application/json ContentType.
type PutClimonJSONRequestBody = ClimonPostRequest

// DeleteConfigatorClusterJSONRequestBody defines body for DeleteConfigatorCluster for application/json ContentType.
type DeleteConfigatorClusterJSONRequestBody = ClusterRequest

// PostConfigatorClusterJSONRequestBody defines body for PostConfigatorCluster for application/json ContentType.
type PostConfigatorClusterJSONRequestBody = ClusterRequest

// DeleteConfigatorProjectJSONRequestBody defines body for DeleteConfigatorProject for application/json ContentType.
type DeleteConfigatorProjectJSONRequestBody = ProjectDeleteRequest

// PostConfigatorProjectJSONRequestBody defines body for PostConfigatorProject for application/json ContentType.
type PostConfigatorProjectJSONRequestBody = ProjectPostRequest

// PutConfigatorProjectJSONRequestBody defines body for PutConfigatorProject for application/json ContentType.
type PutConfigatorProjectJSONRequestBody = ProjectPostRequest

// DeleteConfigatorRepositoryJSONRequestBody defines body for DeleteConfigatorRepository for application/json ContentType.
type DeleteConfigatorRepositoryJSONRequestBody = RepositoryDeleteRequest

// PostConfigatorRepositoryJSONRequestBody defines body for PostConfigatorRepository for application/json ContentType.
type PostConfigatorRepositoryJSONRequestBody = RepositoryPostRequest

// PutConfigatorRepositoryJSONRequestBody defines body for PutConfigatorRepository for application/json ContentType.
type PutConfigatorRepositoryJSONRequestBody = RepositoryPostRequest

// DeleteDeployerJSONRequestBody defines body for DeleteDeployer for application/json ContentType.
type DeleteDeployerJSONRequestBody = DeployerDeleteRequest

// PostDeployerJSONRequestBody defines body for PostDeployer for application/json ContentType.
type PostDeployerJSONRequestBody = DeployerPostRequest

// PutDeployerJSONRequestBody defines body for PutDeployer for application/json ContentType.
type PutDeployerJSONRequestBody = DeployerPostRequest

// PostRegisterAgentMultipartRequestBody defines body for PostRegisterAgent for multipart/form-data ContentType.
type PostRegisterAgentMultipartRequestBody = AgentRequest

// PutRegisterAgentJSONRequestBody defines body for PutRegisterAgent for application/json ContentType.
type PutRegisterAgentJSONRequestBody = AgentRequest

// PostStoreCredJSONRequestBody defines body for PostStoreCred for application/json ContentType.
type PostStoreCredJSONRequestBody = StoreCredRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List of APIs provided by the service
	// (GET /api-docs)
	GetApiDocs(c *gin.Context)
	// deploy the application
	// (DELETE /climon)
	DeleteClimon(c *gin.Context)
	// deploy the application
	// (POST /climon)
	PostClimon(c *gin.Context)
	// deploy the application
	// (PUT /climon)
	PutClimon(c *gin.Context)
	// Delete the application
	// (DELETE /configator/cluster)
	DeleteConfigatorCluster(c *gin.Context)
	// deploy the application
	// (POST /configator/cluster)
	PostConfigatorCluster(c *gin.Context)
	// deploy the application
	// (DELETE /configator/project)
	DeleteConfigatorProject(c *gin.Context)
	// deploy the application
	// (POST /configator/project)
	PostConfigatorProject(c *gin.Context)
	// deploy the application
	// (PUT /configator/project)
	PutConfigatorProject(c *gin.Context)
	// deploy the application
	// (DELETE /configator/repository)
	DeleteConfigatorRepository(c *gin.Context)
	// deploy the application
	// (POST /configator/repository)
	PostConfigatorRepository(c *gin.Context)
	// deploy the application
	// (PUT /configator/repository)
	PutConfigatorRepository(c *gin.Context)
	// deploy the application
	// (DELETE /deployer)
	DeleteDeployer(c *gin.Context)
	// deploy the application
	// (POST /deployer)
	PostDeployer(c *gin.Context)
	// deploy the application
	// (PUT /deployer)
	PutDeployer(c *gin.Context)
	// Register agent
	// (GET /register/agent)
	GetRegisterAgent(c *gin.Context)
	// Register agent
	// (POST /register/agent)
	PostRegisterAgent(c *gin.Context)
	// Register agent
	// (PUT /register/agent)
	PutRegisterAgent(c *gin.Context)
	// Kubernetes readiness and liveness probe endpoint
	// (GET /status)
	GetStatus(c *gin.Context)
	// Delete the application
	// (POST /store/cred)
	PostStoreCred(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetApiDocs operation middleware
func (siw *ServerInterfaceWrapper) GetApiDocs(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiDocs(c)
}

// DeleteClimon operation middleware
func (siw *ServerInterfaceWrapper) DeleteClimon(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteClimon(c)
}

// PostClimon operation middleware
func (siw *ServerInterfaceWrapper) PostClimon(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostClimon(c)
}

// PutClimon operation middleware
func (siw *ServerInterfaceWrapper) PutClimon(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutClimon(c)
}

// DeleteConfigatorCluster operation middleware
func (siw *ServerInterfaceWrapper) DeleteConfigatorCluster(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteConfigatorCluster(c)
}

// PostConfigatorCluster operation middleware
func (siw *ServerInterfaceWrapper) PostConfigatorCluster(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostConfigatorCluster(c)
}

// DeleteConfigatorProject operation middleware
func (siw *ServerInterfaceWrapper) DeleteConfigatorProject(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteConfigatorProject(c)
}

// PostConfigatorProject operation middleware
func (siw *ServerInterfaceWrapper) PostConfigatorProject(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostConfigatorProject(c)
}

// PutConfigatorProject operation middleware
func (siw *ServerInterfaceWrapper) PutConfigatorProject(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutConfigatorProject(c)
}

// DeleteConfigatorRepository operation middleware
func (siw *ServerInterfaceWrapper) DeleteConfigatorRepository(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteConfigatorRepository(c)
}

// PostConfigatorRepository operation middleware
func (siw *ServerInterfaceWrapper) PostConfigatorRepository(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostConfigatorRepository(c)
}

// PutConfigatorRepository operation middleware
func (siw *ServerInterfaceWrapper) PutConfigatorRepository(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutConfigatorRepository(c)
}

// DeleteDeployer operation middleware
func (siw *ServerInterfaceWrapper) DeleteDeployer(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteDeployer(c)
}

// PostDeployer operation middleware
func (siw *ServerInterfaceWrapper) PostDeployer(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostDeployer(c)
}

// PutDeployer operation middleware
func (siw *ServerInterfaceWrapper) PutDeployer(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutDeployer(c)
}

// GetRegisterAgent operation middleware
func (siw *ServerInterfaceWrapper) GetRegisterAgent(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetRegisterAgent(c)
}

// PostRegisterAgent operation middleware
func (siw *ServerInterfaceWrapper) PostRegisterAgent(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostRegisterAgent(c)
}

// PutRegisterAgent operation middleware
func (siw *ServerInterfaceWrapper) PutRegisterAgent(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutRegisterAgent(c)
}

// GetStatus operation middleware
func (siw *ServerInterfaceWrapper) GetStatus(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetStatus(c)
}

// PostStoreCred operation middleware
func (siw *ServerInterfaceWrapper) PostStoreCred(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostStoreCred(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {

	errorHandler := options.ErrorHandler

	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/api-docs", wrapper.GetApiDocs)

	router.DELETE(options.BaseURL+"/climon", wrapper.DeleteClimon)

	router.POST(options.BaseURL+"/climon", wrapper.PostClimon)

	router.PUT(options.BaseURL+"/climon", wrapper.PutClimon)

	router.DELETE(options.BaseURL+"/configator/cluster", wrapper.DeleteConfigatorCluster)

	router.POST(options.BaseURL+"/configator/cluster", wrapper.PostConfigatorCluster)

	router.DELETE(options.BaseURL+"/configator/project", wrapper.DeleteConfigatorProject)

	router.POST(options.BaseURL+"/configator/project", wrapper.PostConfigatorProject)

	router.PUT(options.BaseURL+"/configator/project", wrapper.PutConfigatorProject)

	router.DELETE(options.BaseURL+"/configator/repository", wrapper.DeleteConfigatorRepository)

	router.POST(options.BaseURL+"/configator/repository", wrapper.PostConfigatorRepository)

	router.PUT(options.BaseURL+"/configator/repository", wrapper.PutConfigatorRepository)

	router.DELETE(options.BaseURL+"/deployer", wrapper.DeleteDeployer)

	router.POST(options.BaseURL+"/deployer", wrapper.PostDeployer)

	router.PUT(options.BaseURL+"/deployer", wrapper.PutDeployer)

	router.GET(options.BaseURL+"/register/agent", wrapper.GetRegisterAgent)

	router.POST(options.BaseURL+"/register/agent", wrapper.PostRegisterAgent)

	router.PUT(options.BaseURL+"/register/agent", wrapper.PutRegisterAgent)

	router.GET(options.BaseURL+"/status", wrapper.GetStatus)

	router.POST(options.BaseURL+"/store/cred", wrapper.PostStoreCred)

	return router
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+yazW7jNhDHX4Vge3TitHvzLc0WRbCLbhBveymCgBbHNrcSyR1SWRiB370gqe8Py0pi",
	"d936ZpkiOfPjnzPUSM80UolWEqQ1dPZMTbSGhPmf1yuQ9h6+pmCsu9aoNKAV4Fsj9hih/19YSPxfS4UJ",
	"s3RGF0Iy3NAJtRsNdEaNRSFXdFv8wRDZxl1HsQBp33Ckv2HzupG2E4rwNRUInM7+yt2sWVqb7GFCrbCx",
	"G8EDI0KGOYWS5Wxq8QUi62bLqBqtpIEOrKmxKgF8FNxdtowFybUS0nY0Ni2vjFTp17IXc1s6jL2JRaLk",
	"e4jBQr8S4tRYwEfJEu8QBxOh0B7AjN6EViIk+bYW0ZpYRRZAuB+TTwiHJUtjh+1ikYqYk2y4rpUaguMs",
	"MJpFHWb8njeRaM3QZlYIaSyLY+Bds+k4XQnZ49adbyS+saMvQgzMQE/n+9Dqe2eWpAY4WSrMTeoUqkhA",
	"pbY93ufQ4PvbNRCmdSwiL8F8vIYehbSwAhwQTZVAlW7Dv9Kyirjeg47VhjihBEMw6IdotokV4/1yu1Nm",
	"R9hxq9cnNb+yHqqQ5B60MsKq7k3/ItEWchkn2/+4LBG06h0xX4MdJmn1mGK8s/Mf9x+PtSEm9AnQ+GGa",
	"o/4ZGoha+kH9irXNauyo+h4qWVU8n1RVfZx95rW6d0QfHYgbGh2TqOq8apZUvL5RcilWKcIuNwMawO8m",
	"f53z02vz07EzUi6hc04656RzTjrkTrtD5X4OxOrXCE2HGfo6h9aq0CIEZoG7HRfmHQm2NuHY5JXZszPw",
	"/H9olPvugALZL2hYRRjnwN2PF4Eop3k5hYPJ4u0YvCaMjo9X41GWFZhGns261yNXpUpSh52AMWzVfUw2",
	"ltnUDJ+As/smxWAPHRbPrUK4QeD95xAE/vIjOzPmm8LuxtQA9h/mc/DGWUgMRAjWtKG7W4VcqjbyOWNz",
	"YgCfAMknDZLc/zr/TK7vbonREIllliZpOVV/j3mjR5E56U+XV5dXzhulQTIt6Iy+u7y6fOeWlNm1Rzhl",
	"WlxwFfmLFXjIRRq75XRGfwN7rcV7d4tbxKAKf/vPV1dt1z598IhMmiQMN3RGPwpjXba+vrs1RKN6Em4f",
	"LTY+fTuPhE+ylq2Ml366iEVEH9wg08jXR8IkLgq2rQvRMdRRaBAZGPuL4r4qGilpIdQOK6eP6RcTBg21",
	"X/frR4QlndEfpmVxeJpVhqddNcHtNki6TeNNpiw2q59mADAPh4/GEauCNKP44ESvTMcSu+B6BITVGH6y",
	"ANMufukZ3z74/Jb22YZZhdP8aW14exd9booHvMNwrpWJvifIAcROyIHMwC4/k9xDrjnJhl6zA/0YvWaP",
	"FQei3PkIeVKsc6T7qvYoPE820NZo9iWqM8uRLBtRAMuS5ohAUCuEHoJ3X8HgpKBX2O4bEY4I9mS13MQ6",
	"GBrOUEdCdTGCZ69OhsNC/pLlQHC7XwOeFNyC5e4wcCSQJ6vROsaebX+GuDdEt80RVsI9H0zZKjO0r2x2",
	"n93pPzujr/S3+MBvl+P1b+3an/oN08htJiwzOqfgry+C88HNgb3Z9r5PW0kaW6EZ2ulSYXLBmWX7r3Xt",
	"o81+Xb212z17aX+nX67uf8NjJ/yyuN8n+Hle1h+2x6RRBMYs07h8a9qw8EO6AJRgwRAExoUEYwiTnMTi",
	"CfyFRrUAUnzpWT28iydmoTBcIUwj/+7heYdiixcOB1q41guNkyoweYqOqOvmX0S4v5+pf99Fp3T7sP0n",
	"AAD//64/C09cLQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
