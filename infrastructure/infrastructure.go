package infrastructure

import (
	"context"
	"fmt"
	"github.com/santoshbachar/navik/constants"
	"gopkg.in/yaml.v2"
	"io"
	"os"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

type Instance struct {
	Name         string `yaml:"name"`
	ProjectID    string `yaml:"projectId"`
	Zone         string `yaml:"zone"`
	InstanceName string `yaml:"instanceName"`
	MachineType  string `yaml:"machineType"`
	SourceImage  string `yaml:"sourceImage"`
	NetworkName  string `yaml:"networkName"`
}

type Infra struct {
	Infrastructure struct {
		Vendor struct {
			GCloud struct {
				Instances []Instance `yaml:"instances"`
			} `yaml:"gcloud"`
		} `yaml:"vendor"`
	} `yaml:"infrastructure"`
}

func Provision() {
	var infra Infra

	dat, err := os.ReadFile(constants.ResourceDir + "navik.infrastructure.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(dat), &infra)
	if err != nil {
		panic(err)
	}

	//fmt.Println("entire playbooks: ", playbooks[0].Tasks)

	gCloud := infra.Infrastructure.Vendor.GCloud

	defaultConfig := gCloud.Instances[0]

	w := io.Writer()
	createInstance(w, defaultConfig.ProjectID, defaultConfig.Zone, defaultConfig.InstanceName, defaultConfig.MachineType, defaultConfig.SourceImage, defaultConfig.NetworkName)

}

func createInstance(w io.Writer, projectID, zone, instanceName, machineType, sourceImage, networkName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// machineType := "n1-standard-1"
	// sourceImage := "projects/debian-cloud/global/images/family/debian-10"
	// networkName := "global/networks/default"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(100),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String(networkName),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %v", err)
	}

	fmt.Fprintf(w, "Instance created\n")

	return nil
}
