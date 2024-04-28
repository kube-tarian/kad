package cert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRootCerts(t *testing.T) {
	certInfo, err := GenerateRootCerts()
	assert.NoError(t, err)
	t.Log(certInfo)
}
