# Nemo as a managed plugin for Docker 1.13+
Using a standalone Docker Engine for dev and test is a common use case. Nemo is available as a Docker Volume plugin on Docker Hub and may be installed on a host OS that has a recent `zfs.ko` kernel module. This is how one would use it on a stock Ubuntu 18.04:

```
sudo apt-get install -y docker.io
sudo docker plugin install --alias nemo --grant-all-permissions nimblestorage/nemo:1.0.0
sudo docker volume create -d nemo myvol1
sudo docker run --rm -it -v myvol1:/data bash
```
A container prompt should present itself:
```
bash-4.4# df /data
Filesystem           1K-blocks      Used Available Use% Mounted on
tank/v2/myvol1        10485760       128  10485632   0% /data
```
**Note:** When used for emulating HPE Nimble Storage, use `--alias nimble` to able to run `docker compose` and such unmodified from production environments.

# Customizing Nemo as a managed plugin
Managed plugins are controlled with "settable" environment variables. Inspecting the plugin once it's installed on the host reveals a few self explanatory options:

```
"Env": [
    {
        "Description": "Log level",
        "Name": "NEMO_DEBUG",
        "Settable": [
            "value"
        ],
        "Value": "debug"
    },
    {
        "Description": "Default volume size (in GiB)",
        "Name": "NEMO_DEFAULTSIZE",
        "Settable": [
            "value"
        ],
        "Value": "10"
    },
    {
        "Description": "Nemo OpenZFS files",
        "Name": "NEMO_POOLDIR",
        "Settable": [
            null
        ],
        "Value": "/var/lib/nemo"
    },
    {
        "Description": "Nemo OpenZFS pool name",
        "Name": "NEMO_POOLNAME",
        "Settable": [
            "value"
        ],
        "Value": "tank"
    },
    {
        "Description": "Nemo OpenZFS pool size (in GiB)",
        "Name": "NEMO_POOLSIZE",
        "Settable": [
            "value"
        ],
        "Value": "64"
    }
]
```

The managed plugin default behavior is to run off a loopback device mapped to a file that persists with the host, no OpenZFS experience needed. If there's already a pool on the host that matches `NEMO_POOLNAME`, that pool will be used instead. The plugin will try its best to maintain the `/dev/loop*` devices across reboots, restarts and upgrades. 

As an example, to use a custom pool name (plugin needs to be disabled):
```
sudo docker plugin set nemo NEMO_POOLNAME=zwimming
sudo docker plugin enable nemo
```

**Pro tip:** The plugin can be disabled on install if changes to the defaults are desired. Just add `--disable` to `docker plugin install`. Additionally, parameters may be added to the install as well, by adding `KEY=VALUE` to the very end of the `docker plugin install` string.

Running multiple instances of the plugin against multiple pools should be safe.

# Building a managed plugin
The required `Dockerfile` and `Makefile` is available in this directory:
```
sudo make VERSION=myplugin
sudo docker plugin enable nimblestorage/nemo:myplugin
sudo docker volume create -d nimblestorage/nemo:myplugin myvol1
```
