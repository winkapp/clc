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
  
### flags

`-config`

  Points to config file. (default "clc.yaml")
  
`-root`

  Points to root directory for configs and output. (default "./")

