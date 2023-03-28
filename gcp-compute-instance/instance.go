package instance

import (
	"encoding/base64"
	"log"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type InstanceController interface {
	Start() error
	Stop() error
	GetStatus() string
}

type instance struct {
	ProjectId      string
	Zone           string
	Name           string
	ComputeService *compute.Service
}

func (i *instance) GetStatus() string {
	c, _ := i.ComputeService.Instances.Get(i.ProjectId, i.Zone, i.Name).Do()
	return c.Status
}

func (i *instance) Start() error {
	_, err := i.ComputeService.Instances.Start(i.ProjectId, i.Zone, i.Name).Do()
	if err != nil {
		return err
	}

	return nil
}

func (i *instance) Stop() error {
	_, err := i.ComputeService.Instances.Stop(i.ProjectId, i.Zone, i.Name).Do()
	if err != nil {
		return err
	}

	return nil
}

func New(pjId, zone, name, credFileb64 string) InstanceController {
	f, err := base64.StdEncoding.DecodeString(credFileb64)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := compute.NewService(nil, option.WithCredentialsJSON(f))
	if err != nil {
		log.Fatal(err)
	}

	return &instance{
		ProjectId:      pjId,
		Zone:           zone,
		Name:           name,
		ComputeService: srv,
	}
}
