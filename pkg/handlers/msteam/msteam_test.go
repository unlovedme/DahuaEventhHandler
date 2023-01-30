package msteam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

// Tests the Init() function
func TestInit(t *testing.T) {
	s := &MSTeams{}
	expectedError := fmt.Errorf(msteamsErrMsg, "Missing MS teams webhook URL")

	var Tests = []struct {
		ms  config.MSTeams
		err error
	}{
		{config.MSTeams{WebhookURL: "somepath"}, nil},
		{config.MSTeams{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.MSTeams = tt.ms
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}

// Tests ObjectCreated() by passing v1.Pod
func TestObjectCreated(t *testing.T) {
	expectedCard := TeamsMessageCard{
		Type:       messageType,
		Context:    context,
		ThemeColor: msTeamsColors["Normal"],
		Summary:    "kubewatch notification received",
		Title:      "kubewatch",
		Text:       "",
		Sections: []TeamsMessageCardSection{
			{
				ActivityTitle: "A `pod` in namespace `new` has been `Created`:\n`foo`",
				Markdown:      true,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http