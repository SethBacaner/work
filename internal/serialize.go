package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	taskDelimiter = ":::"
)

var SerializationError = errors.New("task couldn't be serialized")
var DeserializationError = errors.New("task couldn't be deserialized")

// TaskDefinition will be needed later for custom serialization plus we need it to make sure the task is actually registered
func Serialize(name string, args interface{}, registry Registry) (string, TaskDefinition, error) {
	definition, err := registry.DefinitionFor(name)
	if err != nil {
		return "", TaskDefinition{}, err
	}

	jsonEncoded, err := json.Marshal(args)
	if err != nil {
		return "", TaskDefinition{}, err
	}

	encoded := base64.StdEncoding.EncodeToString(jsonEncoded)

	return formatSerializedTask(name, encoded), *definition, nil
}

func Deserialize(serialized string, registry Registry) (interface{}, TaskDefinition, error) {
	name, serializedArgs, err := getSerializedTaskParts(serialized)
	if err != nil {
		return nil, TaskDefinition{}, err
	}

	definition, err := registry.DefinitionFor(name)
	if err != nil {
		return nil, TaskDefinition{}, DeserializationError
	}

	// TODO: args must make a copy of definition.Args
	// TODO: probably better to have the TaskDefinition do the actual deserialization actually
	args := definition.argsFactory()

	decodedBase64 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(serializedArgs))

	err = json.NewDecoder(decodedBase64).Decode(args)
	if err != nil {
		return nil, TaskDefinition{}, err
	}

	return args, *definition, nil
}

func formatSerializedTask(name string, serializedArgs string) string {
	return fmt.Sprintf("%s%s%s", name, taskDelimiter, serializedArgs)
}

func getSerializedTaskParts(serialized string) (name, serializedArgs string, err error) {
	parts := strings.Split(serialized, taskDelimiter)
	if len(parts) < 2 {
		return name, serializedArgs, DeserializationError
	}

	name = parts[0]
	serializedArgs = parts[1]

	return name, serializedArgs, nil

}
