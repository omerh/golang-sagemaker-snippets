package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

func v1() {
	// aws session
	mySession := session.Must(session.NewSession())

	// Create a SageMakerRuntime client from just a session.
	svc := sagemakerruntime.New(mySession, aws.NewConfig().WithRegion("us-west-2"))

	// build payload body request for generation
	body := map[string]interface{}{
		"inputs": "<s>[INST]how will ai change the wolds[/INST]",
		"parameters": map[string]int{
			"max_new_tokens": 2000,
		},
	}

	// marshal the body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// Create an amazon sagemaker invoke with stremign
	input := &sagemakerruntime.InvokeEndpointWithResponseStreamInput{
		EndpointName:           aws.String("jumpstart-dft-meta-textgeneration-llma2-7b-chat"),
		ContentType:            aws.String("application/json"),
		InferenceComponentName: aws.String("meta-textgeneration-llama-2-7b-f-20240410-132745"),
		Accept:                 aws.String("application/json"),
		Body:                   jsonBody,
	}

	// send the request to the endpoint
	resp, err := svc.InvokeEndpointWithResponseStream(input)
	if err != nil {
		panic(err)
	}

	// request strem
	stream := resp.GetStream()

	// itterate over the events from the stream
	for event := range stream.Events() {
		switch v := event.(type) {
		case *sagemakerruntime.PayloadPart:
			/// chunk streamed
			fmt.Print(string(v.Bytes))
		default:
			fmt.Println("done")
		}
	}
}
