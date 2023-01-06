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

package msteam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

var msteamsErrMsg = `
%s

You need to set the MS teams webhook URL,
using --webhookURL, or using environment variables:

export KW_MSTEAMS_WEBHOOKURL=webhook_url

Command line flags will override environment variables

`

var msTeamsColors = map[string]string{
	"Normal":  "2DC72D",
	"Warning": "DEFF22",
	"Danger":  "8C1A1A",
}

// Constants for Sending a Card
const (
	messageType = "MessageCard"
	context     = "http://schema.org/extensions"
)

// TeamsMessageCard is for the Card Fields to send in Teams
// The Documentation is in https://docs.microsoft.com/en-us/outlook/actionable-messages/card-reference#card-fields
type T