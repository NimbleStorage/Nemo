/*
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
*/

package main

const (
    DVName         = "Name"
    DVCreatedAt    = "CreatedAt"
    DVMountpoint   = "Mountpoint"
)

const (
    VDDestroySuffix             = "-removed-"
    VDDescription               = "description"
    VDFsMode                    = "fsMode"
    VDFsOwner                   = "fsOwner"
    VDHelp                      = "help"
    VDSize                      = "size"
    VDSizeInGiB                 = "sizeInGiB"
    VDCloneOf                   = "cloneOf"
    VDSnapshot                  = "snapshot"
    VDCreateSnapshot            = "createSnapshot"
    VDDestroyOnRm               = "destroyOnRm"
    VDDestroyOnDetach           = "destroyOnDetach"
    VDImportVol                 = "importVol"
    VDImportVolAsClone          = "importVolAsClone"
    VDMountConflictDelay        = "mountConflictDelay"
    VDRestore                   = "restore"
)

const (
    NVDVolSize                  = "VolSizeMiB"
    NVDVolUsage                 = "VolUsageMiB"
    NVDSnapshots                = "Snapshots"
    NVDDescriptionDefault       = "Docker knows this dataset as "
    NVDApplicationCategory      = "ApplicationCategory"
    NVDCachePinned              = "CachePinned"
    NVDCachingEnabled           = "CachingEnabled"
    NVDConnections              = "Connections"
    NVDDedupeEnabled            = "DedupeEnabled"
    NVDEncryptionCipher         = "EncryptionCipher"
    NVDGroup                    = "Group"
    NVDLimitIOPS                = "LimitIOPS"
    NVDLimitMBPS                = "LimitMBPS"
    NVDLimitSnapPercentOfSize   = "LimitSnapPercentOfSize"
    NVDLimitVolPercentOfSize    = "LimitVolPercentOfSize"
    NVDPerfPolicy               = "PerfPolicy"
    NVDPool                     = "Pool"
    NVDFolder                   = "Folder"
    NVDSnapUsageMiB             = "SnapUsageMiB"
    NVDThinlyProvisioned        = "ThinlyProvisioned"
    NVDDestroyOnRm              = "destroyOnRm"
    NVDDestroyOnDetach          = "destroyOnDetach"
    NVDMountConflictDelay       = "mountConflictDelay"
    NVDDescription              = "Description"
    NVDSnapName                 = "Name"
    NVDSnapTime                 = "Time"
)

const (
    OpenZFSQuota            = "quota"
    OpenZFSUsed             = "used"
    OpenZFSTrue             = "true"
    OpenZFSFalse            = "false"
    OpenZFSSourceLocal      = "local"
    OpenZFSMountpoint       = "mountpoint"
    OpenZFSMounted          = "mounted"
    OpenZFSYes              = "yes"
    OpenZFSNo               = "no"
    OpenZFSDivider          = "/"
    OpenZFSSnapDivider      = "@"
    OpenZFSUserprop         = ":"
    OpenZFSReadonly         = "readonly"
    OpenZFSSnapdir          = "snapdir"
    OpenZFSRecordsize       = "recordsize"
    OpenZFSCompression      = "compression"
    OpenZFSPrimarycache     = "primarycache"
    OpenZFSSecondarycache   = "secondarycache"
    OpenZFSLogbias          = "logbias"
    OpenZFSExec             = "exec"
    OpenZFSDevices          = "devices"
    OpenZFSAtime            = "atime"
    OpenZFSXattr            = "xattr"
    OpenZFSCreation         = "creation"
    OpenZFSPropName         = "name"
    OpenZFSHelp             = "OpenZFS"
)

const (
    UsrPropBaseDS                   = ":baseds"
    UsrPropCreated                  = ":created"
    UsrPropName                     = ":name"
    UsrPropMountID                  = ":mountid"
    UsrPropDescription              = ":description"
    UsrPropDestroyOnRm              = ":destroyonrm"
    UsrPropDestroyOnDetach          = ":destroyondetach"
    UsrPropApplicationCategory      = ":applicationcategory"
    UsrPropCachePinned              = ":cachepinned"
    UsrPropCachingEnabled           = ":cachingenabled"
    UsrPropConnections              = ":connections"
    UsrPropDedupeEnabled            = ":dedupeenabled"
    UsrPropEncryptionCipher         = ":encryptioncipher"
    UsrPropGroup                    = ":group"
    UsrPropLimitIOPS                = ":limitiops"
    UsrPropLimitMBPS                = ":limitmbps"
    UsrPropLimitSnapPercentOfSize   = ":limitsnappercentofsize"
    UsrPropLimitVolPercentOfSize    = ":limitvolpercentofsize"
    UsrPropPerfPolicy               = ":perfpolicy"
    UsrPropPool                     = ":pool"
    UsrPropFolder                   = ":folder"
    UsrPropSnapUsageMiB             = ":snapusagemib"
    UsrPropThinlyProvisioned        = ":thinlyprovisioned"
    UsrPropMountConflictDelay       = ":mountconflictdelay"
)
const (
    DefaultApplicationCategory      = "Virtual Server"
    DefaultCachePinned              = "false"
    DefaultCachingEnabled           = "true"
    DefaultConnections              = "0"
    DefaultDedupeEnabled            = "false"
    DefaultEncryptionCipher         = "none"
    DefaultLimitIOPS                = "-1"
    DefaultLimitMBPS                = "-1"
    DefaultLimitSnapPercentOfSize   = "-1"
    DefaultLimitVolPercentOfSize    = "100"
    DefaultPerfPolicy               = "DockerDefault"
    DefaultPool                     = "default"
    DefaultFolder                   = ""
    DefaultSnapUsageMiB             = "0"
    DefaultThinlyProvisioned        = "true"
    DefaultDestroyOnRm              = "false"
    DefaultDestroyOnDetach          = "false"
    DefaultMountConflictDelay       = "30"
)

const (
    UnixChown   = "chown"
    UnixChmod   = "chmod"
)

const (
    DVPHelp = ` -o help

Nemo: HPE Nimble Storage Advanced Data Services Emulator for Containers

Create, Clone or Import an existing OpenZFS dataset into a locally scoped 
Docker Volume. All options are optional. Every '-o key=value' will be stored 
on the OpenZFS Dataset.

********************************************************************************

                            D I S C L A I M E R

 HPE Nimble Storage is not using OpenZFS in any of its products. Nemo is a tool
 to educate users how to integrate Advanced Data Services into DevOps workflows 
 using common developer and IT operations tools without owning or using a 
 HPE Nimble Storage product. Nemo is not supported by HPE Nimble Storage.

********************************************************************************
 
Universal Options:
  -o help           This help
  -o size=X         X is the quota of volume specified in GiB, 0 means no quota

Nimble Compatible Global Options:
  -o description=X      X is a vanity description set on the volume
  -o sizeInGiB=X        X is the alternative to 'size'
  -o destroyOnRm=X      X is either 'true' or 'false' to actually destroy a 
                        dataset after it has been removed. Can be "imported" 
                        with 'importVol'. Global default runtime flag available
  -o fsMode=X           X X is 1 to 4 octal digits that represent the file mode
                        to be applied to the root directory of the filesystem
  -o fsOwner=X:Y        X:Y is the uid:gid that should own the root directory of
                        the filesystem, in the form of uid:gid (string or nums)

Nimble Compatible Clone Options:
  -o cloneOf=X          X is the name of Docker Volume to create a clone of
  -o snapshot=X         X is the name of the snapshot to base the clone on 
                        (optional, if missing, a new snapshot is created)
  -o createSnapshot=X   'true' or 'false', indicates that a new snapshot of the
                        volume should be taken and used for the clone
  -o destroyOnDetach=X  indicates that the dataseut (including snapshots)
                        backing this volume should be destroyed when this volume
                        is unmounted or detached

Nimble Compatible Import Dataset as Clone Options:
  -o importVolAsClone=X X is an exisiting dataset without ZDVP properties to
                        import into a ZDVP dataset as a clone
  -o snapshot=X         X is an optional dataset snapshot to import clone from
  -o createSnapshot=X   X is either 'true' or 'false' to create a above snapshot
  -o destroyOnDetach=X  X indicates that the dataset (including snapshots)
                        backing this volume should be destroyed when this volume
                        is unmounted or detached

Nimble Compatible Import Options:
  -o importVol=X    X is an exisiting unmounted dataset without ZDVP properties
                    to import into a ZDVP dataset
  -o snapshot=X     X is an optional dataset snapshot to import from
  -o restore=X      X is dataset snapshot name to roll back to on import
  -o forceImport=X  Vanity flag, accepts any value, all imports are forced

Nimble Compatible Vanity Options (will be displayed in the "Status" field):
  -o limitIOPS=X            Defaults to '-1'
  -o limitMBPS=X            Defaults to '-1'
  -o dedupe=X               Defaults to 'false'
  -o thick=X                Defaults to 'false'
  -o encryption=X           Defaults to 'none'
  -o folder=X               Defaults to 'none'
  -o pool=X                 Defaults to 'default'
  -o perfPolicy=X           Defaults to 'DockerDefault'
  -o protectionTemplate=X   Defaults to 'none'

Additional OpenZFS Options:
 -o help=OpenZFS
`
)

const (
    DVPOpenZFSHelp = ` -o help=OpenZFS

OpenZFS Dataset Properties (default):
  -o readonly=X         X is either ('off') or 'on'
  -o snapdir=X          X is either ('hidden') or 'visible'
  -o recordsize=X       X is '512' to '1048576' bytes, power of 2 ('131072')
  -o compression=X      X is either ('off'), 'on' or compression algorithm
  -o primarycache=X     X is either ('all'), 'none' or 'metadata'
  -o secondarycache=X   X is either ('all'), 'none' or 'metadata'
  -o logbias=X          X is either ('latency') or 'throughput'
  -o exec=X             X is either ('on') or 'off'
  -o devices=X          X is either ('on') or 'off'
  -o atime=X            X is either ('on') or 'off'
  -o xattr=X            X is either ('on') or 'off'

OpenZFS properties are inherited from its parent dataset unless explicitly set.
`
)
