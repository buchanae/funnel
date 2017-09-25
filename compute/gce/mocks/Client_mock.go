package mocks

import mock "github.com/stretchr/testify/mock"
import scheduler "github.com/ohsu-comp-bio/funnel/proto/scheduler"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// StartNode provides a mock function with given fields: tplName, serverAddress, nodeID
func (_m *Client) StartNode(tplName string, serverAddress string, nodeID string) error {
	ret := _m.Called(tplName, serverAddress, nodeID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(tplName, serverAddress, nodeID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Templates provides a mock function with given fields:
func (_m *Client) Templates() []scheduler.Node {
	ret := _m.Called()

	var r0 []scheduler.Node
	if rf, ok := ret.Get(0).(func() []scheduler.Node); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]scheduler.Node)
		}
	}

	return r0
}