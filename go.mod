module github.com/tidal-engineering/terraform-provider-spinnaker

go 1.14

require (
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/pflag v1.0.5
	github.com/spinnaker/spin v1.30.0
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
