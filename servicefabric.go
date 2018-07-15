package servicefabric

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/ido50/requests"
	"github.com/pkg/errors"
)

// DefaultAPIVersion is a default Service Fabric REST API version
const DefaultAPIVersion = "6.0"

var ErrApplicationNotFound = errors.New("service fabric application not found")
var ErrApplicationNotExists = errors.New("service fabric application does not exist")

// Client for Service Fabric.
type ServiceFabricClient struct {
	// endpoint Service Fabric cluster management endpoint
	endpoint string
	// apiVersion Service Fabric API version
	apiVersion string
	// httpClient HTTP client
	httpClient *requests.HTTPClient
}

func NewServiceFabricClient(httpClient *requests.HTTPClient, endpoint, apiVersion string) (*ServiceFabricClient, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint missing for httpClient configuration")
	}
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}

	return &ServiceFabricClient{
		endpoint:   endpoint,
		apiVersion: apiVersion,
		httpClient: httpClient,
	}, nil
}

func (c ServiceFabricClient) GetApplications() (*ApplicationItemsPage, error) {
	var aggregateAppItemsPages ApplicationItemsPage
	var continueToken string
	for {
		res, _, err := c.getHTTP("Applications/", withContinue(continueToken))
		if err != nil {
			return nil, err
		}

		var appItemsPage ApplicationItemsPage
		err = json.Unmarshal(res, &appItemsPage)
		if err != nil {
			return nil, fmt.Errorf("could not deserialise JSON response: %+v", err)
		}

		aggregateAppItemsPages.Items = append(aggregateAppItemsPages.Items, appItemsPage.Items...)

		continueToken = getString(appItemsPage.ContinuationToken)
		if continueToken == "" {
			break
		}
	}
	return &aggregateAppItemsPages, nil
}

func (c ServiceFabricClient) GetApplication(appName string) (*ApplicationItem, error) {
	var app *ApplicationItem

	res,status, err := c.getHTTP("Applications/"+appName, withParam("api-version",c.apiVersion))

	if status == http.StatusNoContent{
		return nil,ErrApplicationNotExists
	}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &app)
	return app, err
}

func (c ServiceFabricClient) GetDeployment(deploymentName string) (interface{}, error) {
	var deployment interface{}

	res, status, err := c.getHTTP("ComposeDeployments/"+deploymentName, withParam("api-version",c.apiVersion))

	if status == http.StatusNoContent{
		return nil, ErrApplicationNotExists
	}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &deployment)
	return deployment, err
}


func (c ServiceFabricClient) GetServices(appName string) (*ServiceItemsPage, error) {
	var aggregateServiceItemsPages ServiceItemsPage
	var continueToken string
	for {
		res,_, err := c.getHTTP("Applications/"+appName+"/$/GetServices", withContinue(continueToken))
		if err != nil {
			return nil, err
		}

		var servicesItemsPage ServiceItemsPage
		err = json.Unmarshal(res, &servicesItemsPage)
		if err != nil {
			return nil, fmt.Errorf("could not deserialise JSON response: %+v", err)
		}

		aggregateServiceItemsPages.Items = append(aggregateServiceItemsPages.Items, servicesItemsPage.Items...)

		continueToken = getString(servicesItemsPage.ContinuationToken)
		if continueToken == "" {
			break
		}
	}
	return &aggregateServiceItemsPages, nil
}

func (c ServiceFabricClient) GetClusterHealth() (bool, error) {
	res, err := c.getHTTPRaw("$/GetClusterHealth?api-version=6.0&")
	if err != nil {
		return false, fmt.Errorf("error getting cluster health")
	}

	return res == http.StatusOK, nil
}

func (c ServiceFabricClient) DeleteService(serviceId string) error {
	_ ,_, err := c.postHTTP("Services/" + serviceId + "/$/Delete",[]byte{}, withParam("api-version",c.apiVersion))
	if err != nil {
		return errors.Wrap(err,"failed deleting service")
	}

	return nil
}

func (c ServiceFabricClient) DeleteApplication(applicationId string) error {
	_ ,status, err := c.postHTTP("Applications/" + applicationId + "/$/Delete",[]byte{},withParam("api-version",c.apiVersion))

	if err != nil {
		// handle unexpected status
		if status > 200 && status < 300 {
			return nil
		}

		if status == http.StatusNotFound {
			return ErrApplicationNotFound
		}

		return errors.Wrap(err,"failed deleting application")
	}

	return nil
}

func (c ServiceFabricClient) DeleteComposeDeployment(deploymentName string) error{
	_ ,status, err := c.postHTTP("ComposeDeployments/" + deploymentName+ "/$/Delete",[]byte{},withParam("api-version",c.apiVersion))
	if err != nil {
		// handle unexpected status
		if status > 200 && status < 300 {
			return nil
		}

		if status == http.StatusNotFound {
			return ErrApplicationNotFound
		}
		return errors.Wrap(err,"failed deleting compose deployment")
	}

	return nil

}
func (c ServiceFabricClient) GetServiceExtension(appType, applicationVersion, serviceTypeName, extensionKey string, response interface{}) error {
	res,_, err := c.getHTTP("ApplicationTypes/"+appType+"/$/GetServiceTypes", withParam("ApplicationTypeVersion", applicationVersion))
	if err != nil {
		return fmt.Errorf("error requesting service extensions: %v", err)
	}

	var serviceTypes []ServiceType
	err = json.Unmarshal(res, &serviceTypes)
	if err != nil {
		return fmt.Errorf("could not deserialise JSON response: %+v", err)
	}

	for _, serviceTypeInfo := range serviceTypes {
		if serviceTypeInfo.ServiceTypeDescription.ServiceTypeName == serviceTypeName {
			for _, extension := range serviceTypeInfo.ServiceTypeDescription.Extensions {
				if strings.EqualFold(extension.Key, extensionKey) {
					err = xml.Unmarshal([]byte(extension.Value), &response)
					if err != nil {
						return fmt.Errorf("could not deserialise extension's XML value: %+v", err)
					}
					return nil
				}
			}
		}
	}
	return nil
}

func (c ServiceFabricClient) GetServiceExtensionMap(service *ServiceItem, app *ApplicationItem, extensionKey string) (map[string]string, error) {
	extensionData := ServiceExtensionLabels{}
	err := c.GetServiceExtension(app.TypeName, app.TypeVersion, service.TypeName, extensionKey, &extensionData)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{}
	if extensionData.Label != nil {
		for _, label := range extensionData.Label {
			labels[label.Key] = label.Value
		}
	}

	return labels, nil
}

func (c ServiceFabricClient) GetProperties(name string) (bool, map[string]string, error) {
	nameExists, err := c.nameExists(name)
	if err != nil {
		return false, nil, err
	}

	if !nameExists {
		return false, nil, nil
	}

	properties := make(map[string]string)

	var continueToken string
	for {
		res,_, err := c.getHTTP("Names/"+name+"/$/GetProperties", withContinue(continueToken), withParam("IncludeValues", "true"))
		if err != nil {
			return false, nil, err
		}

		var propertiesListPage PropertiesListPage
		err = json.Unmarshal(res, &propertiesListPage)
		if err != nil {
			return false, nil, fmt.Errorf("could not deserialise JSON response: %+v", err)
		}

		for _, property := range propertiesListPage.Properties {
			if property.Value.Kind != "String" {
				continue
			}
			properties[property.Name] = property.Value.Data
		}

		continueToken = propertiesListPage.ContinuationToken
		if continueToken == "" {
			break
		}
	}

	return true, properties, nil
}

func (c ServiceFabricClient) nameExists(propertyName string) (bool, error) {
	res, err := c.getHTTPRaw("Names/" + propertyName)
	// Get http will return error for any non 200 response code.
	if err != nil {
		return false, err
	}

	return res == http.StatusOK, nil
}

func (c ServiceFabricClient) getHTTP(basePath string, paramsFuncs ...queryParamsFunc) ([]byte,int, error) {
	if c.httpClient == nil {
		return nil,0, errors.New("invalid http client provided")
	}

	var text interface{}
	var status int
	url := c.getURL(basePath, paramsFuncs...)
	err := c.httpClient.
		NewRequest("GET", url).
		Into(&text).
		StatusInto(&status).
		Run()

	if err != nil {
		return nil,status, fmt.Errorf("failed connecting to Service Fabric server, status code %d: %s",status, err)
	}

	b, err := json.Marshal(text)
	return b, status, err

}

func (c ServiceFabricClient) getHTTPRaw(basePath string) (int, error) {
	if c.httpClient == nil {
		return -1, fmt.Errorf("invalid http client provided")
	}

	url := c.getURL(basePath)

	var text string
	var status int
	err := c.httpClient.NewRequest("GET", url).Into(&text).
		StatusInto(&status).
		Run()
	if err != nil {
		return -1, fmt.Errorf("failed to connect to Service Fabric server: %s", err)
	}
	return status, nil
}

func (c ServiceFabricClient) getURL(basePath string, paramsFuncs ...queryParamsFunc) string {
	params := []string{"api-version=" + c.apiVersion}

	for _, paramsFunc := range paramsFuncs {
		params = paramsFunc(params)
	}

	return fmt.Sprintf("/%s?%s",basePath, strings.Join(params, "&"))
}

func (c ServiceFabricClient) postHTTP(basePath string,body []byte,paramsFuncs ...queryParamsFunc)([]byte,int,error){
	if c.httpClient == nil {
		return nil,0, errors.New("invalid http client provided")
	}

	url := c.getURL(basePath,paramsFuncs...)
	var responseBody interface{}
	var status int
	err := c.
		httpClient.
		NewRequest("POST", url).
		Into(&responseBody).
		StatusInto(&status).
		Run()

	if err != nil {
		return nil,status, fmt.Errorf("failed connecting to Service Fabric server, status code %d: %s", status, err)
	}

	if responseBody != nil{
		b, err := json.Marshal(responseBody)
		return b, status, err
	}

	return []byte{},status,nil

}

func getString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}