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

package sqs

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	QueueUrl string
	Region   string
	svc      *sqs.SQS
}

func Load(config map[string]string) *SQS {
	result := SQS{}

	if queueUrl, ok := config["queueUrl"]; ok {
		result.QueueUrl = queueUrl
	}

	if region, ok := config["region"]; ok {
		result.Region = region
	} else {
		result.Region = "eu-central-1"
	}

	result.svc = sqs.New(session.New(), &aws.Config{Region: aws.String(result.Region)})
	return &result
}

func (s *SQS) Read() []*sqs.Message {
	log.Println("sqs: receiving messages")
	response, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.QueueUrl),
		MessageAttributeNames: []*string{aws.String("All")},
		WaitTimeSeconds:       aws.Int64(20),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return response.Messages
}

func (s *SQS) Delete(message *sqs.Message) {
	log.Println("sqs: deleting message", message.ReceiptHandle)
	_, err := s.svc.DeleteMessage(&sqs.DeleteMessageInput{QueueUrl: aws.String(s.QueueUrl), ReceiptHandle: message.ReceiptHandle})
	if err != nil {
		log.Println(err)
	}
}
