package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"

	compute "google.golang.org/api/compute/v1"
)

const (
	statusSuccess = iota
	statusFailure

	confTmpName = ".compute-key-dec"
)

// ComputeClient abstracts google compute engine operations
type ComputeClient struct {
	computeSrv *compute.Service
	ctx        context.Context
	project    string
	zone       string
}

// SnapshotStatus represents the status of a single snapshot operation
type SnapshotStatus struct {
	Disk          string
	status        int
	failureReason string
}

// Text returns a chat-friendly text version of the status
func (s *SnapshotStatus) Text() string {
	switch s.status {
	case statusSuccess:
		return ":green_heart:"
	case statusFailure:
		return fmt.Sprintf(":negative_squared_cross_mark: - %s", s.failureReason)
	default:
		log.Printf("unknown status: %d", s.status)
		return string(s.status)
	}
}

// NewComputeClient is a constructor for a compute client.
func NewComputeClient(configPath, project, zone string) (*ComputeClient, error) {
	err := prepareCredentials(configPath)
	if err != nil {
		return nil, err
	}
	// If the variable not provided - try to find a base64 encoded secret from k8s

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return nil, err
	}

	computeService, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return &ComputeClient{
		computeSrv: computeService,
		ctx:        ctx,
		project:    project,
		zone:       zone,
	}, nil
}

// SnapshotDisks runs snapshot for a slice of disk names and returns a corresponding slice of snapshot statuses.
func (c *ComputeClient) SnapshotDisks(disks []string) ([]SnapshotStatus, error) {
	report := make([]SnapshotStatus, 0)

	n := time.Now()
	snapshotPostfix := fmt.Sprintf("%d-%d-%d-%d", n.Year(), n.Month(), n.Day(), n.Unix())

	for _, disk := range disks {
		rb := &compute.Snapshot{
			Name: fmt.Sprintf("bu-%s-%s", disk[0:31], snapshotPostfix),
		}
		_, err := c.computeSrv.Disks.CreateSnapshot(c.project, c.zone, disk, rb).Context(c.ctx).Do()
		s := SnapshotStatus{Disk: disk, status: statusSuccess}
		if err != nil {
			s.status = statusFailure
			s.failureReason = err.Error()
		}

		report = append(report, s)
	}

	return report, nil
}

func prepareCredentials(configPath string) error {
	if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
		log.Printf("Using %s as gce default credentials.", path)
		return nil
	}

	log.Printf("GOOGLE_APPLICATION_CREDENTIALS env variable is not provided. Trying to read the base64 encoded config file %s", configPath)
	encodedConf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	conf, err := base64.StdEncoding.DecodeString(string(encodedConf))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(confTmpName, conf, 0644)
	if err != nil {
		return err
	}

	// This is a weird way of passing the credentials.
	// The problem is that the GKE is not injecting default credentials into the pods.
	return os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", confTmpName)
}
