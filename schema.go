package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"google.golang.org/protobuf/encoding/protojson"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

func dump(rpc *metal.Client) error {
	var req pb.ReadSchemaRequest

	resp, err := rpc.Metal.ReadSchema(rpc.Context(), &req)
	if err != nil {
		return err
	}

	doc := resp.GetSchema()

	if !jsonFlag { // parse json as yaml and re-marshal
		var obj map[string]interface{}

		b, err := protojson.MarshalOptions{
			UseProtoNames: true,
		}.Marshal(doc)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(b, &obj); err != nil {
			return err
		}

		return yaml.NewEncoder(os.Stdout).Encode(obj)
	}

	b, err := protojson.MarshalOptions{
		Multiline:     true,
		UseProtoNames: true,
	}.Marshal(doc)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func load(rpc *metal.Client, filename string) error {
	fin, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fin.Close()

	data, err := io.ReadAll(fin)
	if err != nil {
		return err
	}

	var doc pb.Schema

	switch filepath.Ext(filename) {
	case ".json":
		if err := protojson.Unmarshal(data, &doc); err != nil {
			return err
		}
	case ".yaml", ".yml":
		var jsonMap map[string]interface{}

		if err := yaml.Unmarshal(data, &jsonMap); err != nil {
			return err
		}

		jsonData, err := json.Marshal(jsonMap)
		if err != nil {
			return err
		}

		if err := protojson.Unmarshal(jsonData, &doc); err != nil {
			return err
		}

	default:
		return errors.New("unknown file type")
	}

	req := pb.CreateSchemaRequest_builder{
		Schema: &doc,
	}.Build()

	_, err = rpc.Metal.CreateSchema(rpc.Context(), req)

	return err
}
