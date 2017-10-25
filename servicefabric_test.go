package servicefabric

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func setupClient() Client {
	httpClient := &mockHTTPClient{}
	sfClient, _ := NewClient(
		httpClient,
		"Test",
		"1.0",
		"",
		"",
		"")
	return sfClient
}

func TestGetApplications(t *testing.T) {
	expected := &ApplicationItemsPage{
		ContinuationToken: nil,
		Items: []ApplicationItem{
			{
				HealthState: "Ok",
				ID:          "TestApplication",
				Name:        "fabric:/TestApplication",
				Parameters: []*struct {
					Key   string `json:"Key"`
					Value string `json:"Value"`
				}{

					{"Param1", "Value1"},

					{"Param2", "Value2"},
				},
				Status:      "Ready",
				TypeName:    "TestApplicationType",
				TypeVersion: "1.0.0",
			},
			{
				HealthState: "Ok",
				ID:          "TestApplication2",
				Name:        "fabric:/TestApplication2",
				Parameters: []*struct {
					Key   string `json:"Key"`
					Value string `json:"Value"`
				}{

					{"Param1", "Value1"},

					{"Param2", "Value2"},
				},
				Status:      "Ready",
				TypeName:    "TestApplication2Type",
				TypeVersion: "1.0.0",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetApplications()
	if err != nil {
		t.Errorf("Exception thrown %v", err)
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual != expected")
	}
}

func TestGetServices(t *testing.T) {
	expected := &ServiceItemsPage{
		ContinuationToken: nil,
		Items: []ServiceItem{
			{
				HasPersistedState: true,
				HealthState:       "Ok",
				ID:                "TestApplication/TestService",
				IsServiceGroup:    false,
				ManifestVersion:   "1.0.0",
				Name:              "fabric:/TestApplication/TestService",
				ServiceKind:       "Stateful",
				ServiceStatus:     "Active",
				TypeName:          "TestServiceType",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetServices("TestApplication")
	if err != nil {
		t.Errorf("Exception thrown %v", err)
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual != expected")
	}
}

func TestGetPartitions(t *testing.T) {
	expected := &PartitionItemsPage{
		ContinuationToken: nil,
		Items: []PartitionItem{
			{
				CurrentConfigurationEpoch: struct {
					ConfigurationVersion string `json:"ConfigurationVersion"`
					DataLossVersion      string `json:"DataLossVersion"`
				}{
					ConfigurationVersion: "12884901891",
					DataLossVersion:      "131496928071680379",
				},
				HealthState:       "Ok",
				MinReplicaSetSize: 3,
				PartitionInformation: struct {
					HighKey              string `json:"HighKey"`
					ID                   string `json:"Id"`
					LowKey               string `json:"LowKey"`
					ServicePartitionKind string `json:"ServicePartitionKind"`
				}{
					HighKey:              "9223372036854775807",
					ID:                   "bce46a8c-b62d-4996-89dc-7ffc00a96902",
					LowKey:               "-9223372036854775808",
					ServicePartitionKind: "Int64Range",
				},
				PartitionStatus:      "Ready",
				ServiceKind:          "Stateful",
				TargetReplicaSetSize: 3,
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetPartitions("TestApplication", "TestApplication/TestService")
	if err != nil {
		t.Errorf("Exception thrown %v", err)
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual != expected")
	}
}

func TestGetReplicas(t *testing.T) {
	expected := &ReplicaItemsPage{
		ContinuationToken: nil,
		Items: []ReplicaItem{
			{
				ReplicaItemBase: &ReplicaItemBase{
					Address:                      "{\"Endpoints\":{\"\":\"localhost:30001+bce46a8c-b62d-4996-89dc-7ffc00a96902-131496928082309293\"}}",
					HealthState:                  "Ok",
					LastInBuildDurationInSeconds: "1",
					NodeName:                     "_Node_0",
					ReplicaRole:                  "Primary",
					ReplicaStatus:                "Ready",
					ServiceKind:                  "Stateful",
				},
				ID: "131496928082309293",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplication", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-7ffc00a96902")
	if err != nil {
		t.Errorf("Exception thrown %v", err)
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual != expected")
	}
}

func TestGetInstances(t *testing.T) {
	expected := &InstanceItemsPage{
		ContinuationToken: nil,
		Items: []InstanceItem{
			{
				ReplicaItemBase: &ReplicaItemBase{
					Address:                      "{\"Endpoints\":{\"\":\"http:\\/\\/localhost:8081\"}}",
					HealthState:                  "Ok",
					LastInBuildDurationInSeconds: "3",
					NodeName:                     "_Node_0",
					ReplicaStatus:                "Ready",
					ServiceKind:                  "Stateless",
				},
				ID: "131497042182378182",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplication", "TestApplication/TestService", "824091ba-fa32-4e9c-9e9c-71738e018312")
	if err != nil {
		t.Errorf("Exception thrown %v", err)
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual != expected")
	}
}

func TestGetServiceExtension(t *testing.T) {
	sfClient := setupClient()
	service := &ServiceItem{
		HasPersistedState: true,
		HealthState:       "Ok",
		ID:                "TestApplication/TestService",
		IsServiceGroup:    false,
		ManifestVersion:   "1.0.0",
		Name:              "fabric:/TestApplication/TestService",
		ServiceKind:       "Stateful",
		ServiceStatus:     "Active",
		TypeName:          "Test",
	}
	actual, err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "Test", service, &ResponseType{})
	extension := actual.(*ResponseType)
	if err != nil {
		t.Errorf("Exception thrown %v", err)
	}
	if extension.Test.Value != "value1" {
		t.Error("Extension value does not equal value1")
	}
	if extension.Test.Key != "key1" {
		t.Error("Extension key does not equal key1")
	}
}

func TestGetServiceExtensionNoMatchingServiceTypeName(t *testing.T) {
	sfClient := setupClient()
	service := &ServiceItem{
		HasPersistedState: true,
		HealthState:       "Ok",
		ID:                "TestApplication/TestService",
		IsServiceGroup:    false,
		ManifestVersion:   "1.0.0",
		Name:              "fabric:/TestApplication/TestService",
		ServiceKind:       "Stateful",
		ServiceStatus:     "Active",
		TypeName:          "Test1",
	}
	extensions, err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "MissingKey", service, &ResponseType{})
	if extensions != nil {
		t.Error("Should have returned nil interface as no matching extensions")
	}
	if err != nil {
		t.Error("Should not have thrown")
	}
}

func TestGetServiceExtensionNoMatchingExtensions(t *testing.T) {
	sfClient := setupClient()
	service := &ServiceItem{
		HasPersistedState: true,
		HealthState:       "Ok",
		ID:                "TestApplication/TestService",
		IsServiceGroup:    false,
		ManifestVersion:   "1.0.0",
		Name:              "fabric:/TestApplication/TestService",
		ServiceKind:       "Stateful",
		ServiceStatus:     "Active",
		TypeName:          "Test",
	}
	extensions, err := sfClient.GetServiceExtension("TestApplication", "1.0.1", "Test", service, &ResponseType{})
	if extensions != nil {
		t.Error("Should have returned nil interface as no matching extensions")
	}
	if err != nil {
		t.Error("Should not have thrown")
	}
}

func TestGetServiceExtensionWrongType(t *testing.T) {
	sfClient := setupClient()
	service := &ServiceItem{
		HasPersistedState: true,
		HealthState:       "Ok",
		ID:                "TestApplication/TestService",
		IsServiceGroup:    false,
		ManifestVersion:   "1.0.0",
		Name:              "fabric:/TestApplication/TestService",
		ServiceKind:       "Stateful",
		ServiceStatus:     "Active",
		TypeName:          "Test",
	}
	actual, _ := sfClient.GetServiceExtension("TestApplication", "1.0.0", "Test", service, &WrongType{})
	extension, ok := actual.(*ResponseType)
	if ok {
		t.Error("Type assertion should have failed")
	}
	if extension != nil {
		t.Error("Extension should have been nil")
	}
}

type ResponseType struct {
	XMLName xml.Name `xml:"Tests"`
	Test    struct {
		XMLName xml.Name `xml:"Test"`
		Key     string   `xml:"Key,attr"`
		Value   string   `xml:",chardata"`
	}
}

type WrongType struct {
	Prop1 string `json:"Prop1"`
	Prop2 string `json:"Prop2"`
}
