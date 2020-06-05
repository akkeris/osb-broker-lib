module github.com/akkeris/osb-broker-lib

replace github.com/akkeris/go-open-service-broker-client/v2 => /Workspace/golang/src/github.com/akkeris/go-open-service-broker-client/v2

go 1.13

require (
	github.com/akkeris/go-open-service-broker-client/v2 v2.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.4
	github.com/prometheus/client_golang v1.6.0
	k8s.io/klog v1.0.0
)
