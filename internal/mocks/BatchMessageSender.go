// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"
import sqs "github.com/aws/aws-sdk-go/service/sqs"

// BatchMessageSender is an autogenerated mock type for the BatchMessageSender type
type BatchMessageSender struct {
	mock.Mock
}

// SendMessageBatch provides a mock function with given fields: _a0
func (_m *BatchMessageSender) SendMessageBatch(_a0 *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	ret := _m.Called(_a0)

	var r0 *sqs.SendMessageBatchOutput
	if rf, ok := ret.Get(0).(func(*sqs.SendMessageBatchInput) *sqs.SendMessageBatchOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqs.SendMessageBatchOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*sqs.SendMessageBatchInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
