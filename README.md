# Nemo: HPE Nimble Storage Advanced Data Services Emulator for Containers
Nemo is a Docker Volume plugin that strives to be compatible with the HPE Nimble Storage Docker Volume plugin. Nemo uses OpenZFS to emulate the snapshot and cloning capabilities of the HPE Nimble Storage plugin. OpenZFS is also capable of storing volume metadata similar how NimbleOS stores metadata which makes it simple to emulate the behavior.

# Disclaimer
HPE Nimble Storage is not using OpenZFS in any of its products. Nemo is a tool to educate users how to integrate Advanced Data Services into DevOps workflows using common developer and IT operations tools without owning or using a HPE Nimble Storage product. Nemo is not supported by HPE Nimble Storage.

# Limitations
While aspiring to be as feature complete as possible, certain things are not possible to emulate or to achieve with OpenZFS.
* Not supported by HPE Nimble Storage
* Locally scoped Docker Volume driver - meaning volumes are not shared between nodes in a cluster
* Certain volume options are "vanity options", meaning they get set, but doesn't do anything. Run `docker volume create -d nemo -o help` for list of vanity options
* The output of `docker volume inspect` is not fully accurate with several missing representations and all field values are strings, not the actual JSON representation of its data type.

# Requirements
Nemo has only been built and tested with OpenZFS on Ubuntu 18.04 and is the preferred runtime environment. Nemo will eventually be packaged in binary form to multiple outlets to ease integration into popular tools and environments. Nemo is compatible with Docker 1.13+ and is 100% Docker Volume API compatible. 

# Building and Running Nemo
Building on a vanilla installed Ubuntu 18.04:
```
sudo apt-get update && sudo apt-get install zfsutils-linux libzfslinux-dev golang git docker.io 
git clone https://github.com/NimbleStorage/Nemo && cd Nemo
make
```

Nemo assumes a pre-created OpenZFS pool by the name of 'tank', run `./nemo -h` for optional placements. If there is no extra devices on your host OS to create a pool, here's a recipe to create a pool on a sparse file through a loop device:
```
dd if=/dev/zero of=/tmp/tank bs=512 seek=20M count=1
sudo losetup /tmp/tank
sudo zpool create tank /dev/loop0
```

Nemo will now be able to run with default settings:
```
sudo ./nemo
```

Switch to another terminal, you should now be able to run:
```
sudo docker volume create -d nemo -o help
```

# Help
Please file issues on this GitHub repo for bug reports, inquiries and help tickets. If you like to chat, join the HPE DEV slack at https://www.labs.hpe.com/slack, #NimbleStorage

# Contributing
Contributions to Nemo from outside HPE require a contributor license agreement (CLA) with HPE. It's a legality to ensure contributors understands their contributions are subject to the Open Source license this software is made available.

# License
```
(C) Copyright 2018 Hewlett Packard Enterprise Development LP

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
    License for the specific language governing permissions and limitations
    under the License.
```
Full terms available in [LICENSE](LICENSE)
