package main

import (
	"flag"
	"log"

	"github.com/jasonlvhit/gocron"
)

var (
	reportTime     = flag.String("backup-time", "10:00:00", "At what time to perform daily backup and send out the report. In format xx:xx:xx")
	slackToken     = flag.String("slack-token", "", "Token for the slack chat.")
	slackChannel   = flag.String("slack-channel", "", "Slack channel to post reports.")
	kubeconfig     = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file.")
	inCluster      = flag.Bool("in-cluster", false, "Specify this if the application is deployed to a k8s cluster. If you specify 'in-cluster' there's no need to specify a 'kubeconfig' path.")
	computeConfig  = flag.String("compute-config", "", "GKE ONLY. Path to a base64 encoded default credentials json file.")
	computeProject = flag.String("project", "wattx-infra", "GCE project")
	computeZone    = flag.String("zone", "europe-west1-c", "GCE zone")
	testRun        = flag.Bool("test-run", false, "Specify this if you want just to run the backup once and exit")
)

func main() {
	flag.Parse()

	log.Println("Running WATTx Backup Tool....")

	compute, err := NewComputeClient(*computeConfig, *computeProject, *computeZone)
	if err != nil {
		log.Fatal(err)
	}

	kube, err := NewKubeClient(*inCluster, *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	bot := NewSlackBot(*slackToken, *slackChannel)

	snapshooter := NewSnapshooter(bot, kube, compute)

	if *testRun {
		snapshooter.Exec()
	} else {
		gocron.Every(1).Day().At(*reportTime).Do(snapshooter.Exec)
		<-gocron.Start()
	}
}
