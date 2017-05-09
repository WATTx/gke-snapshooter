package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

// Snapshooter is responsible for orchestration of the whole process of snapshotting:
// - getting the disks names (using kube)
// - doing snapshots (using compute)
// - reporting to slack (using bot)
//
// It handles errors by printing them and sending to a slack channel at the same time.
type Snapshooter struct {
	bot     *SlackBot
	kube    *KubeClient
	compute *ComputeClient
}

// NewSnapshooter is a constructor for a new Snapshooter.
func NewSnapshooter(bot *SlackBot, kube *KubeClient, compute *ComputeClient) *Snapshooter {
	return &Snapshooter{
		bot:     bot,
		kube:    kube,
		compute: compute,
	}
}

// Exec runs the whole snapshotting process.
func (s *Snapshooter) Exec() {
	log.Println("Snapshotting...")

	disks, err := s.kube.GetDisks()
	if err != nil {
		s.postFailure(err)
		return
	}

	status, err := s.compute.SnapshotDisks(disks)
	if err != nil {
		s.postFailure(err)
		return
	}

	report, err := renderReport(status)
	if err != nil {
		s.postFailure(err)
		return
	}

	err = s.bot.Post(report)
	if err != nil {
		s.postFailure(err)
		return
	}
}

func (s *Snapshooter) postFailure(snapshotErr error) {
	msg := fmt.Sprintf("Failure during the snapshot: %s", snapshotErr)
	log.Println(msg)
	err := s.bot.Post(msg)
	if err != nil {
		log.Printf("CRITICAL: can't inform about errors: %s", err)
	}
}

func renderReport(status []SnapshotStatus) (string, error) {
	tpl, err := template.New("report.tpl").ParseFiles("./report.tpl")
	if err != nil {
		return "", err
	}

	var report bytes.Buffer
	err = tpl.Execute(&report, status)
	if err != nil {
		return "", err
	}

	return report.String(), nil
}
