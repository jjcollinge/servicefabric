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
		false)
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
		t.Error("actual should equal expected")
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
		t.Error("actual should equal expected")
	}
}

func TestGetServicesWithNonExistentApplicationReturnsDefaultTypeAndError(t *testing.T) {
	sfClient := setupClient()
	expected := &ServiceItemsPage{}
	actual, err := sfClient.GetServices("TestApplicationNonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetServicesWithNonExistentApplicationReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	sfClient := setupClient()
	notExpected := &ServiceItemsPage{
		ContinuationToken: nil,
		Items: []ServiceItem{
			ServiceItem{
				HasPersistedState: true,
			},
		},
	}
	actual, err := sfClient.GetServices("TestApplicationNonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
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

func TestGetPartitionsWithNonExistentApplicationReturnsDefaultTypeAndError(t *testing.T) {
	sfClient := setupClient()
	expected := &PartitionItemsPage{}
	actual, err := sfClient.GetPartitions("TestApplicationNoneExistent", "TestApplication/TestService")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetPartitionsWithNonExistentApplicationReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	sfClient := setupClient()
	notExpected := &PartitionItemsPage{
		ContinuationToken: nil,
		Items: []PartitionItem{
			PartitionItem{
				HealthState: "ERROR",
			},
		},
	}
	actual, err := sfClient.GetPartitions("TestApplicationNoneExistent", "TestApplication/TestService")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
	}
}

func TestGetPartitionsWithNonExistentServiceReturnsDefaultTypeAndError(t *testing.T) {
	sfClient := setupClient()
	expected := &PartitionItemsPage{}
	actual, err := sfClient.GetPartitions("TestApplicationNoneExistent", "TestApplication/TestServiceNonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetPartitionsWithNonExistentServiceReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	sfClient := setupClient()
	notExpected := &PartitionItemsPage{
		ContinuationToken: nil,
		Items: []PartitionItem{
			PartitionItem{
				HealthState: "ERROR",
			},
		},
	}
	actual, err := sfClient.GetPartitions("TestApplicationNoneExistent", "TestApplication/TestServiceNonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
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
		t.Error("actual should be equal to expected")
	}
}

func TestGetReplicasWithNonExistentApplicationReturnsDefaultTypeAndError(t *testing.T) {
	expected := &ReplicaItemsPage{}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplicationNonExistent", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-7ffc00a96902")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetReplicasWithNonExistentApplicationReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	notExpected := &ReplicaItemsPage{
		ContinuationToken: nil,
		Items: []ReplicaItem{
			ReplicaItem{
				ReplicaItemBase: nil,
				ID:              "00001",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplicationNonExistent", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-7ffc00a96902")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
	}
}

func TestGetReplicasWithNonExistentServiceReturnsDefaultTypeAndError(t *testing.T) {
	expected := &ReplicaItemsPage{}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplication", "TestApplication/TestServiceNonExistent", "bce46a8c-b62d-4996-89dc-7ffc00a96902")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetReplicasWithNonExistentServiceReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	notExpected := &ReplicaItemsPage{
		ContinuationToken: nil,
		Items: []ReplicaItem{
			ReplicaItem{
				ReplicaItemBase: nil,
				ID:              "00001",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplication", "TestApplication/TestServiceNonExistent", "bce46a8c-b62d-4996-89dc-7ffc00a96902")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
	}
}

func TestGetReplicasWithNonExistentPartitionReturnsDefaultTypeAndError(t *testing.T) {
	expected := &ReplicaItemsPage{}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplication", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-NonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetReplicasWithNonExistentPartitionReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	notExpected := &ReplicaItemsPage{
		ContinuationToken: nil,
		Items: []ReplicaItem{
			ReplicaItem{
				ReplicaItemBase: nil,
				ID:              "00001",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetReplicas("TestApplication", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-NonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
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
		t.Error("actual should be equal to expected")
	}
}

func TestGetInstancesWithNonExistentApplicationReturnsDefaultTypeAndError(t *testing.T) {
	expected := &InstanceItemsPage{}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplicationNonExistent", "TestApplication/TestService", "824091ba-fa32-4e9c-9e9c-71738e018312")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetInstancesWithNonExistentApplicationReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	notExpected := &InstanceItemsPage{
		ContinuationToken: nil,
		Items: []InstanceItem{
			InstanceItem{
				ReplicaItemBase: nil,
				ID:              "00001",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplicationNonExistent", "TestApplication/TestService", "824091ba-fa32-4e9c-9e9c-71738e018312")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
	}
}

func TestGetInstancesWithNonExistentServiceReturnsDefaultTypeAndError(t *testing.T) {
	expected := &InstanceItemsPage{}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplication", "TestApplication/TestServiceNonExistent", "824091ba-fa32-4e9c-9e9c-71738e018312")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetInstancesWithNonExistentServiceReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	notExpected := &InstanceItemsPage{
		ContinuationToken: nil,
		Items: []InstanceItem{
			InstanceItem{
				ReplicaItemBase: nil,
				ID:              "00001",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplication", "TestApplication/TestServiceNonExistent", "824091ba-fa32-4e9c-9e9c-71738e018312")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
	}
}

func TestGetInstancesWithNonExistentPartitionReturnsDefaultTypeAndError(t *testing.T) {
	expected := &InstanceItemsPage{}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplication", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-NonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(expected, actual)
	if !isEqual {
		t.Error("actual should equal expected")
	}
}

func TestGetInstancesWithNonExistentPartitionReturnsDefaultTypeAndErrorOnly(t *testing.T) {
	notExpected := &InstanceItemsPage{
		ContinuationToken: nil,
		Items: []InstanceItem{
			InstanceItem{
				ReplicaItemBase: nil,
				ID:              "00001",
			},
		},
	}
	sfClient := setupClient()
	actual, err := sfClient.GetInstances("TestApplication", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-NonExistent")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
	isEqual := reflect.DeepEqual(notExpected, actual)
	if isEqual {
		t.Error("actual should not equal notExpected")
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
	var extension, initial ResponseType
	err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "Test", service, &extension)
	if extension == initial {
		t.Error("Extension should have been populated")
	}
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
	var extension, initial ResponseType
	err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "MissingKey", service, &extension)
	if extension != initial {
		t.Error("Should have returned default ResponseType")
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
	var extension, initial ResponseType
	_ = sfClient.GetServiceExtension("TestApplication", "1.0.1", "Test", service, &extension)
	if extension != initial {
		t.Error("Should have returned default ResponseType")
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
	var extension, initial WrongType
	_ = sfClient.GetServiceExtension("TestApplication", "1.0.0", "Test", service, &extension)
	if extension != initial {
		t.Error("Should have returned default ResponseType")
	}
}

func TestGetServiceExtensionWithNonExistentApplicationReturnsDefaultTypeAndErrorName(t *testing.T) {
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
	err := sfClient.GetServiceExtension("TestApplicationNonExistent", "1.0.0", "Test", service, &ResponseType{})
	if err == nil {
		t.Error("Error should have thrown")
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
