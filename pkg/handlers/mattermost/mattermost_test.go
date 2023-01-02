/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mattermost

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestMattermostInit(t *testing.T) {
	s := &Mattermost{}
	expectedError := fmt.Errorf(mattermostErrMsg, "Missing Mattermost channel, url or username")

	var Tests = []struct {
		mattermost config.Mattermost
		err        error
	}{
		{config.Mattermost{Url: "foo", Channel: "bar", Username: "username"}, nil},
		{config.Mattermost{Url: "foo", Channel: "bar"}, expectedError},
		{config.Mattermost{Url: "foo", Username: "username"}, expectedError},
		{config.Mattermost{Channel: "foo", Username: "username"}, expectedError},
		{config.Mattermost{Url: "foo"}, expectedError},
		{config.Mattermost{Channel: "bar"}, expectedError},
		{config.Mattermost{Username: "bar"}, expectedError},
		{config.Mattermost{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Mattermost = tt.mattermost
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
