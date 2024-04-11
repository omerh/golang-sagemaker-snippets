package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime/types"
)

func v2() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatal(err)
	}

	client := sagemakerruntime.NewFromConfig(cfg)

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

	resp, err := client.InvokeEndpointWithResponseStream(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	stream := resp.GetStream()

	// itterate over the events from the stream
	for event := range stream.Events() {
		switch ev := event.(type) {
		case *types.ResponseStreamMemberPayloadPart:
			fmt.Println(string(ev.Value.Bytes))
		default:
			fmt.Print("default")
		}
	}

}
