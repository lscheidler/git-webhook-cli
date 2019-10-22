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

package main

import (
	"flag"
	"log"
	"time"

	awssqs "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/lscheidler/git-webhook-plugin"
	"github.com/lscheidler/git-webhook-plugin-bitbucket"

	"github.com/lscheidler/git-webhook-cli/jenkins"
	"github.com/lscheidler/git-webhook-cli/sqs"
)

func main() {
	jenkinsUrl := flag.String("jenkins-url", "http://localhost:8080", "set jenkins url")
	keep := flag.Bool("keep", false, "keep sqs message")
	sqsUrl := flag.String("sqs-url", "", "set sqs url")
	flag.Parse()

	jnkns := jenkins.Init(*jenkinsUrl)

	for {
		sqs := sqs.Load(map[string]string{"queueUrl": *sqsUrl})
		for _, message := range sqs.Read() {
			event := getEvent(message)
			if err := jnkns.Process(event); err == nil && !*keep {
				sqs.Delete(message)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func getEvent(message *awssqs.Message) gitWebhookPlugin.GitWebhookPlugin {
	attributes := make(map[string]string)
	for key, value := range message.MessageAttributes {
		attributes[key] = *value.StringValue
	}
	bitbkt := bitbucket.Init()
	if bitbkt.IsBitbucketWebhookRequest(attributes) {
		if bitbkt.ValidBody(bitbkt.Attributes()["x-event-key"], []byte(*message.Body)) {
			return bitbkt
		}
	} else {
		log.Println("Webhook unknown", *message.Body)
	}
	return nil
}
