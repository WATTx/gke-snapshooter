# GKE Snapshooter
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FWATTx%2Fgke-snapshooter.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FWATTx%2Fgke-snapshooter?ref=badge_shield)


GKE snapshooter creates snapshots of k8s persistent volumes. It has the following logic:

- query the kubernetes API for PersitentVolumes
- creates a snapshot for each disk via the google cloud compute api
- sends out a slack message about the status of the snapshots

*IT DOESN'T DELETE* any snapshots, since they are based on each other: https://cloud.google.com/compute/docs/disks/create-snapshots

The ability to clean unnecessary snapshots would be useful to add in the future.


## Running locally

Build and run the app:
```
go build .
./gke-snapshooter --help
```

- use `--test-run` flag to run it just once and exit
- you need to provide application default credentials:
    - ensure that you have those credentials: `gcloud auth application-default login`
    - specify the path to the credentials from the previous step, for example: `export GOOGLE_APPLICATION_CREDENTIALS=/home/<USER>/.config/gcloud/application_default_credentials.json`
    - you should have working `kubeconfig`
    
Example of locally running application:
```
./gke-snapshooter -kubeconfig=/home/me/.kube/config -slack-token="<TOKEN>" -slack-channel=mytestchannel --test-run
```


## Deploying

Build the image:
```
docker build -t eu.gcr.io/wattx-infra/gke-snapshooter:{version} .
```

Push the image to the docker registry:
```
gcloud docker -- push eu.gcr.io/wattx-infra/gke-snapshooter:{version}
```

Check the files in the `./kube` folder, it contains the docs and files for the k8s deployment.


## Vendoring dependencies

You need `govendor` to manage vendored dependencies: https://github.com/kardianos/govendor

- Update vendored dependencies: `govendor update +vendor`
- Add new dependencies: `govendor add -h`


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FWATTx%2Fgke-snapshooter.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FWATTx%2Fgke-snapshooter?ref=badge_large)