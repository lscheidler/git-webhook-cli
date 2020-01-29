/*
Copyright 2019 Lars Eric Scheidler

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

package jenkins

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lscheidler/git-webhook-plugin"
)

type Jenkins struct {
	Url string
}

func Init(url string) *Jenkins {
	return &Jenkins{Url: url}
}

func (j *Jenkins) Process(plugin gitWebhookPlugin.GitWebhookPlugin) error {
	uri := fmt.Sprintf("%s/git/notifyCommit?url=%s&branches=%s", j.Url, plugin.GitUrl(), strings.Join(plugin.GitBranches(), ","))
	log.Println("jenkins:", uri)
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("jenkins:", err)
		return err
	} else {
		log.Println("jenkins:", resp)
		return nil
	}
}
