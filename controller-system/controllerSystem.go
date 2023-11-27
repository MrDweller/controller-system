package controllersystem

import (
	"os"

	"github.com/MrDweller/service-registry-connection/models"
	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
)

type ControllerSystem struct {
	models.SystemDefinition
	ServiceRegistryConnection serviceregistry.ServiceRegistryConnection
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

	return &ControllerSystem{
		SystemDefinition:          systemDefinition,
		ServiceRegistryConnection: serviceRegistryConnection,
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
