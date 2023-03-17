package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"io"
	"net/http"
)

func (a *APIHandler) setFailedResponse(c *gin.Context, msg string, err error) {
	a.log.Errorf("failed to submit job", err)
	c.IndentedJSON(http.StatusInternalServerError, &api.Response{
		Status:  "FAILED",
		Message: fmt.Sprintf("%s, %v", msg, err),
	})
}

func (a *APIHandler) getFileContent(c *gin.Context, fileInfo map[string]string) (map[string]string, error) {
	fileContentMap := make(map[string]string)
	// 10 << 20 to have 10 mb max as file size
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("file contents is greater than 10mb, err: %v", err)
	}

	for fileKey, fileName := range fileInfo {
		// Get handler for filename, size and headers
		file, handler, err := c.Request.FormFile(fileKey)
		if err != nil {
			return nil, fmt.Errorf("failed to download file %s, err %v", fileName, err)
		}

		if fileName != handler.Filename {
			return nil, fmt.Errorf("faile name doesnt match expected: %s, got %s", fileName, handler.Filename)
		}

		fileContents, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read content of file %s, err %v", handler.Filename, err)
		}

		fileContentMap[handler.Filename] = string(fileContents)
		file.Close()
	}

	return fileContentMap, nil
}
