package instance

import (
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type Instance struct {
	ProjectId      string
	Zone           string
	Name           string
	ComputeService *compute.Service
}

func (i *Instance) GetStatus() string {
	c, _ := i.ComputeService.Instances.Get(i.ProjectId, i.Zone, i.Name).Do()
	return c.Status
}

func (i *Instance) Start() error {
	resp, err := i.ComputeService.Instances.Start(i.ProjectId, i.Zone, i.Name).Do()
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", resp)
	return nil
}

func (i *Instance) Stop() error {
	resp, err := i.ComputeService.Instances.Stop(i.ProjectId, i.Zone, i.Name).Do()
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", resp)
	return nil
}

func New(pjId, zone, name, credFileb64 string) *Instance {
	f, err := base64.StdEncoding.DecodeString(credFileb64)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := compute.NewService(nil, option.WithCredentialsJSON(f))
	if err != nil {
		log.Fatal(err)
	}

	return &Instance{
		ProjectId:      pjId,
		Zone:           zone,
		Name:           name,
		ComputeService: srv,
	}
}
