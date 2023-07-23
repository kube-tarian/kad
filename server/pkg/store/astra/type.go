package astra

import pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"

const (
	keyspace                        = "capten"
	createKeyspaceQuery             = "CREATE KEYSPACE IF NOT EXISTS capten WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};"
	createClusterEndpointTableQuery = "CREATE TABLE IF NOT EXISTS capten.cluster_endpoint (cluster_id uuid, org_id uuid, cluster_name text, endpoint text, PRIMARY KEY (cluster_id, org_id));"
	createOrgClusterTableQuery      = "CREATE TABLE IF NOT EXISTS capten.org_cluster (org_id uuid, cluster_ids set<uuid>, PRIMARY KEY (org_id));"
)

var (
	UuidSetSpec = &pb.TypeSpec{
		Spec: &pb.TypeSpec_Set_{
			Set: &pb.TypeSpec_Set{
				Element: &pb.TypeSpec{
					Spec: &pb.TypeSpec_Basic_{
						Basic: pb.TypeSpec_UUID,
					},
				},
			},
		},
	}
)
