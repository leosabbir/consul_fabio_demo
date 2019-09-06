package utility

import (
	"fmt"
	"log"
	"os"
	"time"

	consul "github.com/hashicorp/consul/api"
)

//---------------------------------------------------------------------------------------

const (
	// DefaultConsulHTTPAddr points to default location where consul is running
	DefaultConsulHTTPAddr = "localhost:8500"
	// ServiceCheckTTL time for which consul will consider the service is alive. Service has
	// to update TTL to keep itself alive
	ServiceCheckTTL = 30 * time.Second
	// ServiceDeregisterAfter time after which consul will deregister critical service
	ServiceDeregisterAfter = 24 * time.Hour
) // const

//---------------------------------------------------------------------------------------

var consulClient Client
var localIPAddress string

//---------------------------------------------------------------------------------------

func init() {
	localIPAddress = os.Getenv("LOCAL_IP")
	if localIPAddress == "" {
		log.Fatalf("init: Environment variable LOCAL_IP is required to locate service.")
	}
	if err := createConsulClient(); err != nil {
		log.Fatalf("init: Error creating the consul client to connect to consul agent: %s", err)
	}
} // init

//---------------------------------------------------------------------------------------

// Client provides an interface for getting data out of Consul
type Client interface {
	// SErvices Get all services from consul
	Services(string, string) ([]*consul.ServiceEntry, *consul.QueryMeta, error)
	// Service Get single services from consul
	Service(string, string) (string, error)
	// Register registers a service to consul agent with ID provided
	RegisterWithID(string, string, int, *[]string) (string, error)
	// Register registers a service to consul with machineID as ID
	Register(string, int, *[]string) (string, error)
	// DeRegister deregisters a service from consul agent
	DeRegister(string) error
} // Client

//---------------------------------------------------------------------------------------

type client struct {
	TTL    time.Duration
	consul *consul.Client
} // client

//---------------------------------------------------------------------------------------

// GetConsulClient returns a Client interface for given consul address to interact
// (register, deregister, get service) with consul agent
func GetConsulClient() *Client {
	return &consulClient
} // GetConsulClient

//---------------------------------------------------------------------------------------

// createConsulClient creates an instance of ConsulClient
func createConsulClient() error {
	config := consul.DefaultConfig()
	config.Address = GetConsulAddress()
	c, err := consul.NewClient(config)
	if err != nil {
		return err
	}
	consulClient = &client{
		TTL:    ServiceCheckTTL,
		consul: c,
	}
	return nil
} // createConsulClient

//---------------------------------------------------------------------------------------

// Register a service with consul local agent
func (c *client) RegisterWithID(name, serviceID string, port int, tags *[]string) (string, error) {
	log.Printf("Register: Registering %s to consul agent", name)
	reg := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Address: localIPAddress,
		Port:    port,
		Tags:    *tags,
		Check: &consul.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", localIPAddress, port),
			Interval:                       ServiceCheckTTL.String(),
			DeregisterCriticalServiceAfter: ServiceDeregisterAfter.String(),
		},
	}
	err := c.consul.Agent().ServiceRegister(reg)
	if err == nil {
		log.Printf("Register: Registered %s to consul agent", name)
		return serviceID, nil
	}
	log.Printf("Register: ERROR Could not register %s to consul agent: %s", name, err)
	return "", err
} // Register

//---------------------------------------------------------------------------------------

// Register registers a service to consul agent with ID provided
func (c *client) Register(name string, port int, tags *[]string) (string, error) {
	serviceID := fmt.Sprintf("%s_%s", name, localIPAddress)
	return c.RegisterWithID(name, serviceID, port, tags)
}

//---------------------------------------------------------------------------------------

// DeRegister a service with consul local agent
func (c *client) DeRegister(serviceID string) error {
	log.Printf("DeRegistering %s from consul agent", serviceID)
	return c.consul.Agent().ServiceDeregister(serviceID)
} // DeRegister

//---------------------------------------------------------------------------------------

// Service return a service
func (c *client) Service(service, tag string) (string, error) {
	addrs, _, err := c.Services(service, tag)
	if err != nil {
		return "", err
	}
	serviceEntry := addrs[0]
	serviceURL := fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port)
	return serviceURL, nil
} // Service

//---------------------------------------------------------------------------------------

// Services returns all services
func (c *client) Services(service, tag string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	passingOnly := true
	addrs, meta, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
} // Services

//---------------------------------------------------------------------------------------

// GetConsulAddress returns the address of the consul agent
func GetConsulAddress() string {
	consulHTTPAddr := os.Getenv(consul.HTTPAddrEnvName)
	if consulHTTPAddr == "" {
		return fmt.Sprintf("%s:8500", localIPAddress)
	}
	if consulHTTPAddr == "" {
		return DefaultConsulHTTPAddr
	}
	return consulHTTPAddr
} // GetConstulAddress

//---------------------------------------------------------------------------------------
