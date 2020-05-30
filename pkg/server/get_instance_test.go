package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/akkeris/osb-broker-lib/pkg/broker"
	"github.com/akkeris/osb-broker-lib/pkg/metrics"
	"github.com/akkeris/osb-broker-lib/pkg/rest"

	osb "github.com/akkeris/go-open-service-broker-client/v2"
	prom "github.com/prometheus/client_golang/prometheus"
)

func TestGetInstance(t *testing.T) {
	cases := []struct {
		name            string
		validateFunc    func(string) error
		getInstanceFunc func(req *osb.GetInstanceRequest, c *broker.RequestContext) (*broker.GetInstanceResponse, error)
		response        *broker.GetInstanceResponse
		err             error
	}{
		{
			name: "version validation error",
			validateFunc: func(string) error {
				return errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusPreconditionFailed,
				Description: strPtr("oops"),
			},
		},
		{
			name: "returns errors.New",
			getInstanceFunc: func(req *osb.GetInstanceRequest, c *broker.RequestContext) (*broker.GetInstanceResponse, error) {
				return nil, osb.HTTPStatusCodeError{
					StatusCode:  http.StatusInternalServerError,
					Description: strPtr("oops"),
				}
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusInternalServerError,
				Description: strPtr("oops"),
			},
		},
		{
			name: "returns osb.HTTPStatusCodeError",
			getInstanceFunc: func(req *osb.GetInstanceRequest, c *broker.RequestContext) (*broker.GetInstanceResponse, error) {
				return nil, osb.HTTPStatusCodeError{
					StatusCode:  http.StatusBadGateway,
					Description: strPtr("custom error"),
				}
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusBadGateway,
				Description: strPtr("custom error"),
			},
		},
		{
			name: "OK",
			getInstanceFunc: func(req *osb.GetInstanceRequest, c *broker.RequestContext) (*broker.GetInstanceResponse, error) {
				return &broker.GetInstanceResponse{
					GetInstanceResponse: osb.GetInstanceResponse{
						ServiceID:    "12345",
						PlanID:       "67890",
						DashboardURL: strPtr("my.service.to/ABCDE"),
					}}, nil
			},
			response: &broker.GetInstanceResponse{
				osb.GetInstanceResponse{
					ServiceID:    "12345",
					PlanID:       "67890",
					DashboardURL: strPtr("my.service.to/ABCDE"),
				},
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			validateFunc := defaultValidateFunc
			if tc.validateFunc != nil {
				validateFunc = tc.validateFunc
			}

			reg := prom.NewRegistry()
			osbMetrics := metrics.New()
			reg.MustRegister(osbMetrics)

			api := &rest.APISurface{
				Broker: &FakeBroker{
					validateAPIVersion: validateFunc,
					getInstance:        tc.getInstanceFunc,
				},
				Metrics: osbMetrics,
			}

			s := New(api, reg)
			fs := httptest.NewServer(s.Router)
			defer fs.Close()

			config := defaultClientConfiguration()
			config.URL = fs.URL

			client, err := osb.NewClient(config)
			if err != nil {
				t.Error(err)
			}

			actualResponse, err := client.GetInstance(&osb.GetInstanceRequest{
				InstanceID: "ABCDE",
			})
			if err != nil {
				if tc.err != nil {
					if e, a := tc.err, err; !reflect.DeepEqual(e, a) {
						t.Errorf("Unexpected error; expected %v, got %v", e, a)
						return
					}
					return
				}
				t.Error(err)
				return
			}

			if e, a := &tc.response.GetInstanceResponse, actualResponse; !reflect.DeepEqual(e, a) {
				t.Errorf("Unexpected response\n\nExpected: Expected: %#+v\n\nGot: %#+v", e, a)
			}
		})
	}
}
