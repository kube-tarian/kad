package astra

/*
//This is integration test
func TestAstra_RegisterEndpoint(t *testing.T) {
	os.Setenv("CONFIG_PATH", "../../../")
	_, err := config.New()
	dbObj, err := New()
	assert.NoError(t, err, "failed to create ")
	//orgId := uuid.Must(uuid.NewRandom())
	orgId := "fdbb62be-baf8-4637-bb56-b7de46ec6520"
	//err = dbObj.RegisterCluster(orgId, "test23", "127.0.0.1")
	//assert.NoError(t, err)
	clusters, err := dbObj.GetClusters(orgId)
	assert.NoError(t, err)
	fmt.Println(clusters)
	err = dbObj.UpdateCluster(orgId, "test23", "127.0.0.5")
	assert.NoError(t, err)
	endpoint, err := dbObj.GetClusterEndpoint(orgId, "test23")
	assert.NoError(t, err)
	assert.Equal(t, endpoint, "127.0.0.5")
	err = dbObj.DeleteCluster(orgId, "test23")
	assert.NoError(t, err)
}
*/
