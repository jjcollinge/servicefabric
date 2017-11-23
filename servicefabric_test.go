package servicefabric

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetApplications(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleApplications))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

	expected := &ApplicationItemsPage{
		ContinuationToken: nil,
		Items: []ApplicationItem{
			{
				HealthState: "Ok",
				ID:          "TestApplication",
				Name:        "fabric:/TestApplication",
				Parameters: []*AppParameter{
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
				Parameters: []*AppParameter{
					{"Param1", "Value1"},
					{"Param2", "Value2"},
				},
				Status:      "Ready",
				TypeName:    "TestApplication2Type",
				TypeVersion: "1.0.0",
			},
		},
	}

	actual, err := sfClient.GetApplications()
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %+v, want %+v", actual, expected)
	}
}

func TestGetServices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleServices))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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

	actual, err := sfClient.GetServices("TestApplication")
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %+v, want %+v", actual, expected)
	}
}

func TestGetServicesWithNonExistentApplicationReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

	actual, err := sfClient.GetServices("TestApplicationNonExistent")
	if err == nil {
		t.Fatal("Error should have been returned")
	}

	if actual != nil {
		t.Errorf("Got %+v, want nil", actual)
	}
}

func TestGetPartitions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlePartitions))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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

	actual, err := sfClient.GetPartitions("TestApplication", "TestApplication/TestService")
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %+v, want %+v", actual, expected)
	}
}

func TestGetPartitionsReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

	testCases := []struct {
		desc        string
		appName     string
		serviceName string
	}{
		{
			desc:        "With Non Existent Application",
			appName:     "TestApplicationNoneExistent",
			serviceName: "TestApplication/TestService",
		},
		{
			desc:        "With Non Existent Service",
			appName:     "TestApplicationNoneExistent",
			serviceName: "TestApplication/TestServiceNonExistent",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			actual, err := sfClient.GetPartitions(test.appName, test.serviceName)
			if err == nil {
				t.Fatal("Error should have been returned")
			}

			if actual != nil {
				t.Errorf("Got %+v, want nil", actual)
			}
		})
	}
}

func TestGetReplicas(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleReplicas))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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

	actual, err := sfClient.GetReplicas("TestApplication", "TestApplication/TestService", "bce46a8c-b62d-4996-89dc-7ffc00a96902")
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %+v, want %+v", actual, expected)
	}
}

func TestGetReplicasReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

	testCases := []struct {
		desc          string
		appName       string
		serviceName   string
		partitionName string
	}{
		{
			desc:          "With Non Existent Application",
			appName:       "TestApplicationNonExistent",
			serviceName:   "TestApplication/TestService",
			partitionName: "bce46a8c-b62d-4996-89dc-7ffc00a96902",
		},
		{
			desc:          "With Non Existent Service",
			appName:       "TestApplication",
			serviceName:   "TestApplication/TestServiceNonExistent",
			partitionName: "bce46a8c-b62d-4996-89dc-7ffc00a96902",
		},
		{
			desc:          "With Non Existent Partition",
			appName:       "TestApplication",
			serviceName:   "TestApplication/TestService",
			partitionName: "bce46a8c-b62d-4996-89dc-NonExistent",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			actual, err := sfClient.GetReplicas(test.appName, test.serviceName, test.partitionName)
			if err == nil {
				t.Fatal("Error should have been returned")
			}

			if actual != nil {
				t.Errorf("Got %+v, want nil", actual)
			}
		})
	}
}

func TestGetInstances(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleInstances))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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

	actual, err := sfClient.GetInstances("TestApplication", "TestApplication/TestService", "824091ba-fa32-4e9c-9e9c-71738e018312")
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Got %+v, want %+v", actual, expected)
	}
}

func TestGetInstancesReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

	testCases := []struct {
		desc          string
		appName       string
		serviceName   string
		partitionName string
	}{
		{
			desc:          "WithNonExistentApplication",
			appName:       "TestApplicationNonExistent",
			serviceName:   "TestApplication/TestService",
			partitionName: "824091ba-fa32-4e9c-9e9c-71738e018312",
		},
		{
			desc:          "WithNonExistentService",
			appName:       "TestApplication",
			serviceName:   "TestApplication/TestServiceNonExistent",
			partitionName: "824091ba-fa32-4e9c-9e9c-71738e018312",
		},
		{
			desc:          "WithNonExistentPartition",
			appName:       "TestApplication",
			serviceName:   "TestApplication/TestService",
			partitionName: "bce46a8c-b62d-4996-89dc-NonExistent",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			actual, err := sfClient.GetInstances(test.appName, test.serviceName, test.partitionName)
			if err == nil {
				t.Fatal("Error should have been returned")
			}

			if actual != nil {
				t.Errorf("Got %+v, want nil", actual)
			}
		})
	}
}

func TestGetServiceExtension(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleExtensionA))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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
	err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "Test", service.TypeName, &extension)
	if extension == initial {
		t.Error("Extension should have been populated")
	}
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if extension.Test.Value != "value1" {
		t.Errorf("Extension value %q does not equal value1", extension.Test.Value)
	}
	if extension.Test.Key != "key1" {
		t.Errorf("Extension key %q does not equal key1", extension.Test.Key)
	}
}

func TestGetServiceExtensionNoMatchingServiceTypeName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleExtensionA))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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
	err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "MissingKey", service.TypeName, &extension)
	if err != nil {
		t.Fatalf("Should not have thrown: %v", err)
	}
	if extension != initial {
		t.Error("Should have returned default ResponseType")
	}
}

func TestGetServiceExtensionNoMatchingExtensions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleExtensionB))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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
	err := sfClient.GetServiceExtension("TestApplication", "1.0.1", "Test", service.TypeName, &extension)
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}

	if extension != initial {
		t.Error("Should have returned default ResponseType")
	}
}

func TestGetServiceExtensionWrongType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleExtensionA))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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
	err := sfClient.GetServiceExtension("TestApplication", "1.0.0", "Test", service.TypeName, &extension)
	if err != nil {
		t.Fatalf("Exception thrown %v", err)
	}
	if extension != initial {
		t.Error("Should have returned default ResponseType")
	}
}

func TestGetServiceExtensionWithNonExistentApplicationReturnsErrorName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer server.Close()

	sfClient, _ := NewClient(http.DefaultClient, server.URL, "1.0", nil)

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

	err := sfClient.GetServiceExtension("TestApplicationNonExistent", "1.0.0", "Test", service.TypeName, &ResponseType{})
	if err == nil {
		t.Fatal("Error should have thrown")
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
