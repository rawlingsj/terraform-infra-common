package pubsub

import (
	"context"
	"testing"

	"github.com/cloudevents/sdk-go/v2/protocol"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// MockClient is a mock type for the Client interface
type MockClient struct {
	mock.Mock
}

// Send is the mock method for cloudevents.Client's Send method
func (m *MockClient) Send(ctx context.Context, event cloudevents.Event) cloudevents.Result {
	args := m.Called(ctx, event)
	return args.Get(0).(cloudevents.Result)
}

// Request mocks the Request method required by the cloudevents.Client interface
func (m *MockClient) Request(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	args := m.Called(ctx, event)
	return args.Get(0).(*cloudevents.Event), args.Get(1).(cloudevents.Result)
}

// StartReceiver mocks the StartReceiver method required by the cloudevents.Client interface
func (m *MockClient) StartReceiver(ctx context.Context, fn interface{}) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestPublish(t *testing.T) {
	mockClient := new(MockClient)
	testData := []byte(`{"hello":"world"}`)

	// Configure the mock to expect any context and any event of type event.Event
	mockClient.On("Send", mock.Anything, mock.AnythingOfType("event.Event")).Return(protocol.ResultACK)

	opts := Options{
		source:   "test-source",
		ceclient: mockClient,
	}

	err := opts.Publish(context.Background(), testData, "test-type", "test-subject", nil)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// Example struct for testing
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestConvertToNDJSON(t *testing.T) {
	// Setup test data
	people := []Person{
		{Name: "James", Age: 30},
		{Name: "Felix", Age: 25},
	}

	// Expected NDJSON output
	expectedNDJSON := `{"name":"James","age":30}
{"name":"Felix","age":25}`

	// Perform the function under test
	actualNDJSON, err := ConvertToNDJSON(people)

	// Assert no error and correct output
	assert.NoError(t, err, "ConvertToNDJSON should not produce an error")
	assert.Equal(t, expectedNDJSON, string(actualNDJSON), "Output should match the expected NDJSON format")
}
