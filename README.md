# clc: a stupid simple cluster configuration tool

## tl;dr

This simple tool automates the tedious task of creating configurations for running your containers on vanilla CoreOS clusters in development and production.

## usage

The clc utility has three commands and two flags. All of the rest of the options are set through a configuration file. The format for using clc is:

    clc [options] command

### commands

`cc`

  Create a new cloud config file for usage with an IaaS with stock support for file includes 
and etcd based env variables.

`vagrant`

  Creates a user-data file, Vagrantfile and config.rb for use on your local system 
with options mirroring those of your cloud config.

`units`

  Creates unit files for services defined in your configuration.
  
`new`

  Does all of the above. Intended for spinning up a new cluster.
  
### flags

`-config`

  Points to config file. (default "clc.yaml")
  
`-root`

  Points to root directory for configs and output. (default "./")

## rationale

There are plenty of good ways to orchestrate a cluster, like [deis](https://github.com/deis/deis) and [kubernetes](https://github.com/kubernetes/kubernetes). These orchestration layers remove the need for manually creating unit files and worrying about things like environmental variables. For 99% of workloads, those tools are the way to go. That said, there are occasionally times where your workload doesn't quite fit the mold, or you just want to run something on a vanilla cluster, without installing anything. We use this tool for prototyping data pipelines, and for running containers that expose multiple non-http services.

The idea is that you can create a development environment for your cluster in basically no time:

```
clc vagrant
vagrant up
export FLEETCTL_TUNNEL={vagrant machine ip}
```

Then easily iterate on the unit files, re-using much of your existing docker-compose configuration:

```
[edit units config]
clc units
fleetctl start your.service
```

And finally deploy to production on an IaaS using the generated cloud-config file. For instance, if you are using AWS, you can create an AutoScale config and group with the generated could-config.

```
clc cc
```
