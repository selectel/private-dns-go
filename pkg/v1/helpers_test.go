package v1

import "net/http"

const (
	apiErrorJSON         = `{"error":{"code":404,"message":"Not found"}}`
	apiRecordListJSON    = `{"records":[{"type":"A","domain":"sub.example.com.","ttl":-1,"values":["192.168.0.1"],"generated":true}]}`
	apiServicesListJSON  = `{"services":[{"id":"test_service","project":"test_project","network_id":"test_network","resources":{"ram":50,"cpu":0.1},"high_availability":true}]}`
	apiServiceDetailJSON = `{"service":{"id":"test_service","project":"test_project","network_id":"test_network","resources":{"ram":50,"cpu":0.1},"high_availability":true, "addresses": [{"cidr": "192.168.0.0/24", "address": "192.168.0.1"}]}}`
	apiZonesListJSON     = `{"zones":[{"id":"test_zone","name":"example.com.","domain":"example.com.","project":"test_project","reserved_by":"","ttl":3600,"serial_number":9}]}`
	apiZonesDetailsJSON  = `{"zone":{"id":"test_zone","name":"example.com.","domain":"example.com.","project":"test_project","reserved_by":"","ttl":3600,"serial_number":0,"records":[],"bindings":[]}}`
)

type testHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (client *testHTTPClient) Do(req *http.Request) (*http.Response, error) {
	client.request = req

	return client.response, client.err
}
