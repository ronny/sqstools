package sqstools_test

import (
	"testing"

	"github.com/ronny/sqstools/internal/sqstools"
	"github.com/stretchr/testify/assert"
)

func TestSQSQueue(t *testing.T) {
	var testCases = []struct {
		name              string
		url               string
		expectedRegion    string
		expectedAccountID string
		expectedQueueName string
	}{
		{
			name:              "valid",
			url:               "https://sqs.ap-southeast-2.amazonaws.com/1234567890/MyTestQueue",
			expectedRegion:    "ap-southeast-2",
			expectedAccountID: "1234567890",
			expectedQueueName: "MyTestQueue",
		},
		{
			name:              "invalid",
			url:               "foobar",
			expectedRegion:    "",
			expectedAccountID: "",
			expectedQueueName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := sqstools.NewSQSQueue(tc.url)
			assert.Equal(t, tc.expectedRegion, q.Region())
			assert.Equal(t, tc.expectedAccountID, q.AccountID())
			assert.Equal(t, tc.expectedQueueName, q.QueueName())
			assert.Equal(t, tc.url, q.URL())
		})
	}
}
