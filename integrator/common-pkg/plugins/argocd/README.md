# argo-cd helm chart

    - https://github.com/argoproj/argo-helm/tree/main/charts/argo-cd
    - https://blog.knoldus.com/how-to-create-applications-in-argocd/

# Create app using argo-cd

    - https://argo-cd.readthedocs.io/en/release-1.8/user-guide/commands/argocd_app_create/
    - https://blog.knoldus.com/how-to-create-applications-in-argocd/
    - https://argo-cd.readthedocs.io/en/stable/user-guide/helm/

# argo-cd go sdk

    - https://pkg.go.dev/github.com/argoproj/argo-cd/pkg/apiclient

# Argo-cd API documentation

    - https://argo-cd.readthedocs.io/en/release-1.8/developer-guide/api-docs/

```
POST /api/v1/applications HTTP/1.1
Host: localhost:8080
User-Agent: Mozilla/5.0 (X11; Linux aarch64; rv:106.0) Gecko/20100101 Firefox/106.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
Content-Type: application/json
Content-Length: 340
Origin: https://localhost:8080
Connection: keep-alive
Referer: https://localhost:8080/applications?new=%7B%22apiVersion%22%3A%22argoproj.io%2Fv1alpha1%22%2C%22kind%22%3A%22Application%22%2C%22metadata%22%3A%7B%22name%22%3A%22demo%22%7D%2C%22spec%22%3A%7B%22destination%22%3A%7B%22name%22%3A%22%22%2C%22namespace%22%3A%22default%22%2C%22server%22%3A%22https%3A%2F%2Fkubernetes.default.svc%22%7D%2C%22source%22%3A%7B%22path%22%3A%22.%2Fsimple-app%22%2C%22repoURL%22%3A%22https%3A%2F%2Fgithub.com%2FJasmine-Harit%2Fgitops-certification-examples.git%22%2C%22targetRevision%22%3A%22HEAD%22%7D%2C%22project%22%3A%22default%22%7D%7D
Cookie: argocd.token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjpsb2dpbiIsImV4cCI6MTY2NzkyNTcyNywibmJmIjoxNjY3ODM5MzI3LCJpYXQiOjE2Njc4MzkzMjcsImp0aSI6Ijc0OWIxZjFkLTcwMmItNDQxOC1iOTI4LTA3ZmFiNTk2MzkxMiJ9.Cp8ZabyUYiuTeyKtMZkB5R9QVGUMWlW72jhxB9lUUlw
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin

{
	"Request Cookies": {
		"argocd.token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjpsb2dpbiIsImV4cCI6MTY2NzkyNTcyNywibmJmIjoxNjY3ODM5MzI3LCJpYXQiOjE2Njc4MzkzMjcsImp0aSI6Ijc0OWIxZjFkLTcwMmItNDQxOC1iOTI4LTA3ZmFiNTk2MzkxMiJ9.Cp8ZabyUYiuTeyKtMZkB5R9QVGUMWlW72jhxB9lUUlw"
	}
}

{"apiVersion":"argoproj.io/v1alpha1","kind":"Application","metadata":{"name":"demo"},"spec":{"destination":{"name":"","namespace":"default","server":"https://kubernetes.default.svc"},"source":{"path":"./simple-app","repoURL":"https://github.com/Jasmine-Harit/gitops-certification-examples.git","targetRevision":"HEAD"},"project":"default"}}
```

$ curl $ARGOCD_SERVER/api/v1/applications -H "Authorization: Bearer $ARGOCD_TOKEN" 
{"metadata":{"selfLink":"/apis/argoproj.io/v1alpha1/namespaces/argocd/applications","resourceVersion":"37755"},"items":...}



# Instructions

```
# Configure port-forward for kind environment
kubectl port-forward services/argocd-server -n default 8080:443

ARGOCD_SERVER=https://localhost:8080

# Open swagger UI in browser
https://localhost:8080/swagger-ui

# Get bearer token
ARGOCD_TOKEN=$(curl -s $ARGOCD_SERVER/api/v1/session -d $'{"username":"admin","password":"5gN7ue6gqXjX7blE"}' -k |jq .token | sed s/\"//g)

# Get installed applications
curl $ARGOCD_SERVER/api/v1/applications -H "Authorization: Bearer $ARGOCD_TOKEN" -k

# Install application
curl -X POST -k $ARGOCD_SERVER/api/v1/applications -H "Authorization: Bearer $ARGOCD_TOKEN" -H "Content-Type: application/json" -d '
{
    "metadata":{"name":"demo"},
    "spec":{
        "destination":{
            "name":"",
            "namespace":"default",
            "server":"https://kubernetes.default.svc"
        },
        "source":{
            "path":"./simple-app",
            "repoURL":"https://github.com/Jasmine-Harit/gitops-certification-examples.git",
            "targetRevision":"HEAD"
        },
        "syncPolicy":{
            "automated":{
                "prune":false,
                "selfHeal":false
            }
        },
        "project":"default"
    }
}
'

# Delete application
curl -X DELETE -k $ARGOCD_SERVER/api/v1/applications/{name}

For example,

curl -X DELETE -k $ARGOCD_SERVER/api/v1/applications/demo

```

```
We have exposed below APIs to interact with ArgoCD:
Cluster Add
Cluster Delete
Repo Add
Repo Delete
RepoCreds Add
RepoCreds Delete

Sample Payload for each of above APIs can be found in test files.
```
