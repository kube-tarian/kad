package argocd

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type MockServer struct {
	folderIDRecord []byte
	secretIDResp   [][]byte
	passwordResp   []byte
}

func NewMock() *MockServer {
	return &MockServer{}
}

func (m *MockServer) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/cluster.ClusterService/Create", m.handleAddClusterPost).Methods(http.MethodPost)
	return router
}

func (m *MockServer) handleAddClusterPost(res http.ResponseWriter, req *http.Request) {
	fmt.Println("hi i entered handleAddClusterPost")
}
