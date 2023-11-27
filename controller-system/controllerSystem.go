package controllersystem

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	orchestrator_models "github.com/MrDweller/orchestrator-connection/models"
	"github.com/MrDweller/orchestrator-connection/orchestrator"
	"github.com/MrDweller/service-registry-connection/models"

	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
)

type ControllerSystem struct {
	models.SystemDefinition
	ServiceRegistryConnection serviceregistry.ServiceRegistryConnection
	OrchestrationConnection   orchestrator.OrchestratorConnection
}

func NewControllerSystem(address string, port int, systemName string, serviceRegistryAddress string, serviceRegistryPort int) (*ControllerSystem, error) {
	systemDefinition := models.SystemDefinition{
		Address:    address,
		Port:       port,
		SystemName: systemName,
	}

	serviceRegistryConnection, err := serviceregistry.NewConnection(serviceregistry.ServiceRegistry{
		Address: serviceRegistryAddress,
		Port:    serviceRegistryPort,
	}, serviceregistry.SERVICE_REGISTRY_ARROWHEAD_4_6_1, models.CertificateInfo{
		CertFilePath: os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:  os.Getenv("KEY_FILE_PATH"),
		Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
	})
	if err != nil {
		return nil, err
	}

	serviceQueryResult, err := serviceRegistryConnection.Query(models.ServiceDefinition{
		ServiceDefinition: "orchestration-service",
	})
	if err != nil {
		return nil, err
	}

	serviceQueryData := serviceQueryResult.ServiceQueryData[0]

	orchestrationConnection, err := orchestrator.NewConnection(orchestrator.Orchestrator{
		Address: serviceQueryData.Provider.Address,
		Port:    serviceQueryData.Provider.Port,
	}, orchestrator.ORCHESTRATION_ARROWHEAD_4_6_1, orchestrator_models.CertificateInfo{
		CertFilePath: os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:  os.Getenv("KEY_FILE_PATH"),
		Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
	})
	if err != nil {
		return nil, err
	}

	return &ControllerSystem{
		SystemDefinition:          systemDefinition,
		ServiceRegistryConnection: serviceRegistryConnection,
		OrchestrationConnection:   orchestrationConnection,
	}, nil
}

func (controllerSystem *ControllerSystem) StartControllerSystem() error {
	_, err := controllerSystem.ServiceRegistryConnection.RegisterSystem(controllerSystem.SystemDefinition)
	if err != nil {
		return err
	}
	return nil
}

func (controllerSystem *ControllerSystem) StopControllerSystem() error {
	err := controllerSystem.ServiceRegistryConnection.UnRegisterSystem(controllerSystem.SystemDefinition)
	if err != nil {
		return err
	}
	return nil
}

func (controllerSystem *ControllerSystem) SendControll(requestedService orchestrator_models.ServiceDefinition, controll map[string]any) error {
	orchestrationResponse, err := controllerSystem.OrchestrationConnection.Orchestration(requestedService, orchestrator_models.SystemDefinition{
		Address:    controllerSystem.Address,
		Port:       controllerSystem.Port,
		SystemName: controllerSystem.SystemName,
	})
	if err != nil {
		return err
	}

	if len(orchestrationResponse.Response) <= 0 {
		return errors.New("found no providers")
	}
	provider := orchestrationResponse.Response[0]

	payload, err := json.Marshal(controll)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://"+provider.Provider.Address+":"+strconv.Itoa(provider.Provider.Port)+provider.ServiceUri, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := controllerSystem.getClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return errors.New(errorString)
	}

	fmt.Println(string(body))

	return nil

}

func (controllerSystem *ControllerSystem) getClient() (*http.Client, error) {
	// cert, err := tls.LoadX509KeyPair(orchestrator.CertFilePath, orchestrator.KeyFilePath)
	// if err != nil {
	// 	return nil, err
	// }

	// // Load truststore.p12
	// truststoreData, err := os.ReadFile(orchestrator.Truststore)
	// if err != nil {
	// 	return nil, err

	// }

	// // Extract the root certificate(s) from the truststore
	// pool := x509.NewCertPool()
	// if ok := pool.AppendCertsFromPEM(truststoreData); !ok {
	// 	return nil, err
	// }

	client := &http.Client{
		Timeout: 5 * time.Second,
		// Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{
		// 		Certificates:       []tls.Certificate{cert},
		// 		RootCAs:            pool,
		// 		InsecureSkipVerify: false,
		// 	},
		// },
	}
	return client, nil
}
