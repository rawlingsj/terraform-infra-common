package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"

	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"google.golang.org/api/idtoken"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/terraform-infra-common/pkg/httpmetrics"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type PublishConfig struct {
	ProjectID  string `envconfig:"GCP_PROJECT_ID" required:"true"`
	IngressURI string `envconfig:"EVENT_INGRESS_URI" required:"true"`
}

type Options struct {
	source   string
	ceclient cloudevents.Client
}

func getPublishConfig(ctx context.Context) *PublishConfig {
	var c PublishConfig
	err := envconfig.Process("", &c)
	if err != nil {
		clog.FromContext(ctx).Error(err.Error())
		os.Exit(1)
	}
	return &c
}

func InitClient(ctx context.Context) Options {
	log := clog.FromContext(ctx)
	config := getPublishConfig(ctx)

	c, err := idtoken.NewClient(ctx, config.IngressURI)
	if err != nil {
		log.Fatalf("failed to create idtoken client: %v", err) //nolint:gocritic
	}
	ceclient, err := cloudevents.NewClientHTTP(
		cloudevents.WithTarget(config.IngressURI),
		cehttp.WithClient(http.Client{Transport: httpmetrics.WrapTransport(c.Transport)}))
	if err != nil {
		log.Fatalf("failed to create cloudevents client: %v", err)
	}

	return Options{
		ceclient: ceclient,
		source:   getSource(ctx),
	}
}

func (o Options) Publish(ctx context.Context, data []byte, eventType, subject string, extensions map[string]interface{}) error {
	// Create a CloudEvent
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetSource(o.source)
	event.SetSubject(subject)
	event.SetType(eventType)
	event.SetTime(time.Now())
	for k, v := range extensions {
		event.SetExtension(k, v)
	}
	event.SetSpecVersion(cloudevents.VersionV1)
	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return fmt.Errorf("failed to set cloudevent data: %v\n", err)
	}

	const retryDelay = 10 * time.Millisecond
	const maxRetry = 3
	rctx := cloudevents.ContextWithRetriesExponentialBackoff(context.WithoutCancel(ctx), retryDelay, maxRetry)
	if ceresult := o.ceclient.Send(rctx, event); cloudevents.IsUndelivered(ceresult) || cloudevents.IsNACK(ceresult) {
		clog.FromContext(ctx).Errorf("failed to deliver event: %v", ceresult)
		return ceresult
	}

	return nil
}

// todo this seems like it stopped being able to get the internal ip after switching to private VPC
func getSource(ctx context.Context) string {
	s, err := metadata.InternalIP()
	if err != nil {
		clog.FromContext(ctx).Warnf("failed to get GCP Cloud Run IP, falling back to unknown source: %s", err.Error())
		return "unknown"
	}
	return fmt.Sprintf("https://%s", s)
}

// ConvertToNDJSON takes a byte slice containing a JSON array and converts it to NDJSON format.
func ConvertToNDJSON[T any](items []T) ([]byte, error) {
	var ndjsonData []byte
	lastIndex := len(items) - 1 // Get the index of the last item

	for i, item := range items {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("error marshalling item to JSON: %v", err)
		}
		ndjsonData = append(ndjsonData, itemJSON...)
		if i < lastIndex {
			ndjsonData = append(ndjsonData, '\n') // Only add newline if it's not the last item
		}
	}

	return ndjsonData, nil
}
