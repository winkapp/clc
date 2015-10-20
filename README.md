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
  
### config file

Clc config (default is clc.yaml). This file is used to set the discovery url for the production cluster, file includes for the generated cloud-config and user-data files (optional), the location of the file that defines our units, and a directory holding custom templates for any or all of the generated files. If no override templates are supplied, the [stock templates](https://github.com/winkapp/libclc/tree/master/templates) from [libclc](https://github.com/winkapp/libclc) will be used. 

```
discovery: { discovery url like https://discovery.etcd.io/3245sfgsdfgsdgsdfg34 }

unit_directory: { path to your units config file }

templates: { path to your template directory. optional }

files: {specification for files to be included in cloud-config and user-data, to be included on new cluster machines}
  - host_path: {path to the file on the current systen, relative to the directory specified as the clc root}
    path: {path the the desired destination of the file on cluster machines}
    owner: {desired owner of the file on cluster machines}
    permissions: {desired permissions of the file on cluster machines}
    
unit_config:
  units:
    - name: {name of the service}
      type: {multi|single - multi option names the file service-name@.service to enable running multiple instances}
      restart: {always|no - always will restart the service after any exit}
      image: {docker image to run}
      command: {command to pass when starting container. leave blank to use default container command}
      evironment:
        - SOME_ENV_VARIABLE_KEY {used to specify environmental variables to pass to the container. the actual values are set in etcd.}
```

#### note on environmental variables

Although the intention of this tool is to keep our clusters as close to stock as possible for simplicity, we make one exception to facilitate simple usage of etcd-backed environmental variables. Etcd is an incredible tool for sharing ephemeral information such as endpoints for services and api credentials across a cluster. We make use of it as a backing layer to provide env variables to our containers. If you would like to take advantage of this for your own units, first make sure the name/key of the env variable you want passed to your service is in the unit definition in your config file under the `environment` key. See above for example. After that, set the value of the env variable in etc like so:

    etcdctl set /services/{your-service-name}/env/SOME_ENV_VARIABLE_KEY value

Please note that services will not be restarted automatically if the etcd entry for an env variable changes; the env variables are passed on container startup.

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

## libclc

The actual code for creating unit files etc is actually kept in [libclc](https://github.com/winkapp/libclc). The clc project is just a command line tool that takes advantage of libclc. These two things are decoupled because it is sometimes useful to use the libclc logic in unrelated projects to do things like generate unit files on the fly and send them to a cluster programatically.
