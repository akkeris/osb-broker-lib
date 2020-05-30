module github.com/akkeris/osb-broker-lib

replace github.com/akkeris/go-open-service-broker-client/v2 => /Workspace/golang/src/github.com/akkeris/go-open-service-broker-client/v2

go 1.13

require (
	github.com/akkeris/go-open-service-broker-client/v2 v2.0.0-00010101000000-000000000000
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorilla/mux v1.7.4
	github.com/prometheus/client_golang v1.6.0
)
