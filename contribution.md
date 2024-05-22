# Contribution Guidelines
Please read this guide if you plan to contribute to the kad. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues
If you find a bug while working with the kad, please [open an issue on GitHub](https://github.com/kube-tarian/kad/issues/new?labels=kind%2Fbug&template=bug-report.md&title=Bug:) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests
You are more than welcome to open issues in this project to [suggest new features](https://github.com/kube-tarian/kad/issues/new?labels=kind%2Ffeature&template=feature-request.md&title=Feature%20Request:).


## Developing 

Development can be conducted using  GoLang compatible IDE/editor (e.g., Jetbrains GoLand, VSCode).

There are 3 places where you develop new things on kad: config-worker,deployment-worker and agent

### Folder Structure 

kad/
│
├── capten/
│   ├── agent/
│   │   ├── cmd/                # Command-line related code for the agent.
│   │   ├── gin-api-server/     # API server implementation for the agent.
│   │   └── internal/           # Internal packages for the agent.
│   │
│   ├── common-pkg/
│   │   ├── capten-store/       # Capten store related code.
│   │   ├── cert/               # Certificate management code.
│   │   ├── credential/         # Credential management code.
│   │   ├── gerrors/            # Custom error handling.
│   │   ├── k8s/                # Kubernetes related code.
│   │   ├── managed-cluster/    # Managed cluster related code.
│   │   ├── pb/                 # Protocol Buffers generated code.
│   │   ├── plugin-store/       # Plugin store related code.
│   │   ├── plugins/            # Plugins related code.
│   │   ├── postgres/           # PostgreSQL related code.
│   │   ├── temporalclient/     # Temporal client related code.
│   │   ├── vault-cred/         # Vault credential management code.
│   │   └── worker-framework/   # Worker framework related code.
│   │
│   ├── config-worker/
│   │   ├── internal/           # Internal packages for config-worker.
│   │   └── Dockerfile/         # Dockerfile for building config-worker.
│   │
│   ├── database/
│   │   ├── postgres/           # PostgreSQL related code.
│   │       └── migrations/     # Database migration scripts.
│   │
│   └── deployment-worker/
│   │   ├── internal/           # Internal packages for deployment-worker.
│   │   ├── cfg.yaml/           # Configuration file for deployment-worker.
│   │   ├── Dockerfile/         # Dockerfile for building deployment-worker.
│   │   ├── main.go/            # Main Go file for deployment-worker.
│   │   ├── openapi.yaml/       # OpenAPI specification.
│   │   └── model/              # Data models used by deployment-worker.
│   ├── Makefile/               # Makefile for building tasks related to capten.
│
├── charts/
│   ├── kad/                    # Helm charts for KAD application.
│   ├── server/                 # Helm charts for server application.
│
├── dockerfiles/
│   ├── agent/                  # Dockerfile for building the agent application.
│   ├── config-worker/          # Dockerfile for building the config-worker application.
│   ├── deployment-worker/      # Dockerfile for building the deployment-worker application.
│   └── server/                 # Dockerfile for building the server application.
│
├── proto/
│   # Contains .proto files for Protocol Buffers.
│   ├── agent.proto             # Protobuf definitions for agent.
│   ├── capten_plugins.proto    # Protobuf definitions for Capten plugins.
│   ├── cluster_plugins.proto   # Protobuf definitions for cluster plugins.
│   ├── iamoauth.proto          # Protobuf definitions for IAM OAuth.
│   ├── plugin_store.proto      # Protobuf definitions for plugin store.
│   ├── server.proto            # Protobuf definitions for server.
│   └── vault_cred.proto        # Protobuf definitions for Vault credentials.
│
├── .github/
│   ├── workflows/              # GitHub Actions workflows for CI/CD.
│   └── dependabot.yml          # Configuration file for Dependabot.
│
├── server/
│   # Contains the server application code.
│   ├── cmd/                    # Command-line related code for the server.
│   ├── data/                   # Data management code for the server.
│   └── pkg/                    # Package code for the server.
│   ├── Makefile/               # Makefile for building server-specific tasks.
│   ├── docker-compose-postgres.yaml # Docker Compose file for PostgreSQL.
│   ├── server.proto            # Protobuf definitions for server.
│   └── vault_cred.proto        # Protobuf definitions for Vault credentials.
│
├── README.md                   # Main project readme.
└── Makefile                    # Root makefile for building the project.
├── .gitignore                  # Git ignore file.
└── contribution.md             # Contribution guidelines (this file).



## How to Contribute 

You can contribute by adding additional features in agent ,config-worker,deployment-worker and server.You can also add additional rpcs  in the proto files in `./proto` for adding additional feature

We provide a [Makefile](./Makefile) with a few targets that helps build all the parts in a development configuration without a lot of intervention. The more common used targets are:

1. **make gen-protoc**: This command generates the necessary Go code from Protocol Buffers (.proto) files. It creates the required directories and then uses the protoc compiler to generate the code.

2. **make docker-build-server**: This command builds the Docker image for the server application.

3. **make docker-build-kad**: This command triggers the build of Docker images for multiple components: agent, deployment-worker, and config-worker.

4. **make docker-build-agent**: This command builds the Docker image for the agent application.

5. **make docker-build-config**: This command builds the Docker image for the config-worker application.

6. **make docker-build**: This command builds the Docker images for the entire project, including server, agent, deployment-worker, and config-worker.


> **A note on go builds:**
> When running **make docker-build**, the go binaries are built to be run inside a docker container.
> Docker containers are always Linux, regardless of the host OS. 


## General Instructions for contributing Code
This project is written in Golang 

To contribute code.
1. Ensure you are running golang version 1.21 or greater for go module support
2. Set the following environment variables:
    ```
    GO111MODULE=on
    GOFLAGS=-mod=vendor
    ```
3. Fork the project.
4. Clone the project: `git clone https://github.com/[YOUR_USERNAME]/kad && cd kad`
5. kindly refer capten.md file to know the structure of the project.
6. Commit changes *([Please refer the commit message conventions](https://www.conventionalcommits.org/en/v1.0.0/))*
7. Push commits.
8. Open pull request.