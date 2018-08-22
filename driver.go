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
import (
    "sync"
    "errors"
    "github.com/urfave/cli"
    "github.com/docker/go-plugins-helpers/volume"
    log "github.com/Sirupsen/logrus"
)

type VolumeDriver struct {
    mutex       *sync.Mutex
    baseds      string
    mountPoint  string
    prefix      string
    scope       string
    defaultSize string
    driverName  string
    socketOwner string
    destroy     string
}

func newVolumeDriver(args *cli.Context) VolumeDriver {
    driver := VolumeDriver{
                mutex:          &sync.Mutex{},
                scope:          "local",
                baseds:         args.String("baseds"),
		        mountPoint:     args.String("mountpath"),
		        prefix:         args.String("prefix"),
                defaultSize:    args.String("size"),
                driverName:     args.String("name"),
                socketOwner:    args.String("user"),
                destroy:        args.String("destroy"),
               }
    err := Init(&driver)
    if err != nil {
        log.Fatal("Unable to initialize base dataset: ", err)
    }
    return driver
}

func (d VolumeDriver) Unmount(r *volume.UnmountRequest) error {
    d.mutex.Lock()
    var err error
    if r.ID != "" {
        err = d.YankMountId(r.Name, r.ID)
    }
    err = d.UnmountDataset(r.Name)
    if err != nil {
        defer d.mutex.Unlock()
        return err
    }
    defer d.mutex.Unlock()
    return err
}

func (d VolumeDriver) Remove(r *volume.RemoveRequest) error {
    d.mutex.Lock()
    err := d.DestroyDataset(r.Name)
    defer d.mutex.Unlock()
    return err
}

func (d VolumeDriver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
    mpt, err := d.GetDatasetMountpoint(r.Name)
    return &volume.PathResponse{Mountpoint: mpt}, err
}

func (d VolumeDriver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
    d.mutex.Lock()
    err := d.MountDataset(r.Name)
    if err != nil {
        defer d.mutex.Unlock()
        return &volume.MountResponse{Mountpoint: ""}, err
    }
    if r.ID != "" {
        err = d.AppendMountId(r.Name, r.ID)
    }
    mpt, err := d.GetDatasetMountpoint(r.Name)
    defer d.mutex.Unlock()
    return &volume.MountResponse{Mountpoint: mpt}, err
}

func (d VolumeDriver) List() (*volume.ListResponse, error) {
    var list []*volume.Volume
    fs, err := d.ListDatasets()
    for i := range fs {
        name := fs[i][DVName].(string)
        mountpoint := fs[i][DVMountpoint].(string)
        created := fs[i][DVCreatedAt].(string)
        list = append(list, &volume.Volume{Name: name, Mountpoint: mountpoint, CreatedAt: created })
    }
    return &volume.ListResponse{Volumes: list }, err
}

func (d VolumeDriver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
    fs, s, err := d.GetDataset(r.Name);
    if fs[DVName] != "" {
        vol := &volume.Volume{  Name: fs[DVName],
                                Mountpoint: fs[DVMountpoint],
                                CreatedAt: fs[DVCreatedAt],
                                Status: s }
        return &volume.GetResponse{vol}, err
    }
    return &volume.GetResponse{}, errors.New("Can't find volume")
}

func (d VolumeDriver) Create(r *volume.CreateRequest) error {
    var err  error
    var fail error
    var res  error
    d.mutex.Lock()
    err = d.CreateDataset(r.Name, r.Options)
    if err != nil {
        fail = d.DestroyDataset(r.Name)
        if fail != nil {
            res = errors.New(err.Error() + " and while trying to clean up: " + fail.Error())
        } else {
            res = errors.New(err.Error())
        }
    }
    d.mutex.Unlock()
    return res
}

func (d VolumeDriver) Capabilities() *volume.CapabilitiesResponse {
    return &volume.CapabilitiesResponse{
	    Capabilities: volume.Capability{
		    Scope: d.scope,
	    },
    }
}
