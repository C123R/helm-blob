module github.com/PTC-Global/helm-blob

go 1.13

require (
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/aws/aws-sdk-go v1.20.6 // indirect
	github.com/spf13/cobra v0.0.7
	gocloud.dev v0.19.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/sys v0.0.0-20191022100944-742c48ecaeb7 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	helm.sh/helm/v3 v3.1.2
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.2.0+incompatible
