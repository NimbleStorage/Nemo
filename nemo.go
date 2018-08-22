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
    "os"
    "os/exec"
    "time"
    "strings"
    "strconv"
    "errors"
    "github.com/drajen/go-libzfs"
    log "github.com/Sirupsen/logrus"
)

type VolumeEntry map[string]interface{}

var OpenZFSOptions = map[string]zfs.Prop{
    OpenZFSReadonly:        zfs.DatasetPropReadonly,
    OpenZFSSnapdir:         zfs.DatasetPropSnapdir,
    OpenZFSRecordsize:      zfs.DatasetPropRecordsize,
    OpenZFSCompression:     zfs.DatasetPropCompression,
    OpenZFSPrimarycache:    zfs.DatasetPropPrimarycache,
    OpenZFSSecondarycache:  zfs.DatasetPropSecondarycache,
    OpenZFSLogbias:         zfs.DatasetPropLogbias,
    OpenZFSExec:            zfs.DatasetPropExec,
    OpenZFSDevices:         zfs.DatasetPropDevices,
    OpenZFSAtime:           zfs.DatasetPropAtime,
    OpenZFSXattr:           zfs.DatasetPropXattr,
}

var DefaultStatus = map[string]string{
    NVDApplicationCategory:     DefaultApplicationCategory,
    NVDCachePinned:             DefaultCachePinned,
    NVDCachingEnabled:          DefaultCachingEnabled,
    NVDConnections:             DefaultConnections,
    NVDDedupeEnabled:           DefaultDedupeEnabled,
    NVDEncryptionCipher:        DefaultEncryptionCipher,
    NVDLimitIOPS:               DefaultLimitIOPS,
    NVDLimitMBPS:               DefaultLimitMBPS,
    NVDLimitSnapPercentOfSize:  DefaultLimitSnapPercentOfSize,
    NVDLimitVolPercentOfSize:   DefaultLimitVolPercentOfSize,
    NVDPerfPolicy:              DefaultPerfPolicy,
    NVDPool:                    DefaultPool,
    NVDSnapUsageMiB:            DefaultSnapUsageMiB,
    NVDThinlyProvisioned:       DefaultThinlyProvisioned,
    NVDDestroyOnDetach:         DefaultDestroyOnDetach,
    NVDDestroyOnRm:             DefaultDestroyOnRm,
    NVDMountConflictDelay:      DefaultMountConflictDelay,
}

var UsrPropToVanity = map[string]string{
    UsrPropApplicationCategory:     NVDApplicationCategory,
    UsrPropCachePinned:             NVDCachePinned,
    UsrPropCachingEnabled:          NVDCachingEnabled,
    UsrPropConnections:             NVDConnections,
    UsrPropDedupeEnabled:           NVDDedupeEnabled,
    UsrPropEncryptionCipher:        NVDEncryptionCipher,
    UsrPropGroup:                   NVDGroup,
    UsrPropLimitIOPS:               NVDLimitIOPS,
    UsrPropLimitMBPS:               NVDLimitMBPS,
    UsrPropLimitSnapPercentOfSize:  NVDLimitSnapPercentOfSize,
    UsrPropLimitVolPercentOfSize:   NVDLimitVolPercentOfSize,
    UsrPropPerfPolicy:              NVDPerfPolicy,
    UsrPropPool:                    NVDPool,
    UsrPropFolder:                  NVDFolder,
    UsrPropSnapUsageMiB:            NVDSnapUsageMiB,
    UsrPropThinlyProvisioned:       NVDThinlyProvisioned,
    UsrPropDestroyOnRm:             NVDDestroyOnRm,
    UsrPropDestroyOnDetach:         NVDDestroyOnDetach,
    UsrPropMountConflictDelay:      NVDMountConflictDelay,
}

func Init(d *VolumeDriver) error {
    err := d.baseDatasetExist()
    if err != nil {
        return err
    }
    if datasetExist(d.baseds) != nil {
        usrp := make(map[string]string)
        usrp[UsrPropBaseDS] = OpenZFSTrue
        usrp[UsrPropDestroyOnRm] = d.destroy
        usrp[VDSize] = "0"
        err := d.CreateDataset(d.baseds, usrp)
        if err != nil {
            return errors.New(err.Error())
        }
    }
    if err := d.MountDataset(d.baseds); err != nil {
        return errors.New(err.Error())
    }
    log.Info("Ready to serve")
    return nil
}

func (d VolumeDriver) CreateDataset(ds string, usrp map[string]string) error {
    var err             error
    var name            string
    var size            string
    var fs              zfs.Dataset
    //Die
    if _, ok := usrp[VDHelp]; ok {
        return errors.New(getHelp(usrp[VDHelp]))
    }
    //Sanity checks
    if (usrp[VDSize] != "") && (usrp[VDSizeInGiB] != "") {
        return errors.New(VDSize + " and " + VDSizeInGiB + " are mutually exclusive")
    } else if _, ok := usrp[VDSize]; ok {
        size = usrp[VDSize]
    } else if _, ok := usrp[VDSizeInGiB]; ok {
        size = usrp[VDSizeInGiB]
    }
    props := make(map[zfs.Prop]zfs.Property)
    for k, v := range usrp {
        if OpenZFSOptions[k] != 0 {
            props[OpenZFSOptions[k]] = zfs.Property{Value: v}
            delete(usrp, k)
        }
    }
    if strings.HasPrefix(ds, d.baseds) == false {
        name = ds
        ds = d.prefixDs(ds)
        if strings.Contains(name, OpenZFSDivider) {
            return errors.New("Sub-datasets not allowed")
        }
    } else {
        props[zfs.DatasetPropMountpoint] = zfs.Property{Value: d.mountPoint}
    }
    if err := datasetExist(ds); err == nil {
        return errors.New("Dataset already exist. Try to import it to a new volume.")
    }
    //Create
    if _, ok := usrp[VDImportVol]; ok {
        err = d.importDataset(ds, usrp)
    } else if _, ok := usrp[VDImportVolAsClone]; ok {
        err = d.cloneDataset(ds, usrp, props)
    } else if _, ok := usrp[VDCloneOf]; ok {
        err = d.cloneDataset(ds, usrp, props)
    } else {
        fs, err = zfs.DatasetCreate(ds, zfs.DatasetTypeFilesystem, props)
        if err != nil {
            defer fs.Close()
            return err
        }
        defer fs.Close()
    }
    if err != nil {
        return err
    }
    fs, err = d.openDataset(ds)
    if err != nil {
        return err
    }
    for k, v := range usrp {
        if strings.Contains(k, OpenZFSUserprop) {
            err = d.setProperty(&fs, strings.ToLower(k), v)
            if err != nil {
                return err
            }
        }  else {
            err = d.setProperty(&fs, OpenZFSUserprop + strings.ToLower(k), v)
            if err != nil {
                return err
            }
        }
    }
    //Special props
    if err := d.setProperty(&fs, OpenZFSQuota, d.quotaCalc(size));err != nil {
        return err
    }
    //Internal props
    if err := d.setProperty(&fs, UsrPropCreated, timeString()); err != nil {
        return err
    }
    if name != "" {
        if err := d.setProperty(&fs, UsrPropName, name); err != nil {
            return err
        }
    }
    if _, ok := usrp[UsrPropBaseDS]; ok == false {
        if err := d.setProperty(&fs, UsrPropBaseDS, OpenZFSFalse); err != nil {
            return err
        }
    }
    if _, ok := usrp[VDDestroyOnRm]; ok {
        d.setProperty(&fs, UsrPropDestroyOnRm, usrp[VDDestroyOnRm])
    } else {
        d.setProperty(&fs, UsrPropDestroyOnRm, d.destroy)
    }
	defer fs.Close()
    // Post fixups
    if _, ok := usrp[VDFsMode]; ok == true {
        err = d.runFsCmd(name, UnixChmod, usrp[VDFsMode])
        if err != nil {
            return err
        }
    }
    if _, ok := usrp[VDFsOwner]; ok == true {
        err = d.runFsCmd(name, UnixChown, usrp[VDFsOwner])
        if err != nil {
            return err
        }
    }
    return err
}

func (d VolumeDriver) ListDatasets() ([]VolumeEntry, error) {
    var vols []VolumeEntry
    base, err := zfs.DatasetOpen(d.baseds)
    for i := range base.Children {
        ds := base.Children[i]
        if d.getProperty(&ds, UsrPropName) != "" {
            mpt,_ := d.GetDatasetMountpoint(d.getProperty(&ds, UsrPropName))
            e := VolumeEntry{   DVName: d.getProperty(&ds, UsrPropName),
                                DVMountpoint: mpt,
                                DVCreatedAt: d.getProperty(&ds, UsrPropCreated)}
            vols = append(vols, e)
        }
    }
    defer base.Close()
    return vols, err
}

func (d VolumeDriver) DestroyDataset(ds string) error {
    fs, err := d.openDataset(ds)
    if err != nil {
        return nil
    }
    if err := fs.Unmount(0); err != nil {
        return err
    }
    if d.getProperty(&fs, UsrPropDestroyOnRm) == OpenZFSTrue {
        if err := fs.DestroyRecursive(); err != nil {
            return errors.New("Recursive destroy failed: " + err.Error())
        }
    } else {
        if err := d.setProperty(&fs, UsrPropName, ""); err != nil {
            return err
        }
        if err := fs.Rename(d.prefixDs(ds) + VDDestroySuffix +
                timeString(), false, true); err != nil {
            return err
        }
    }
    defer fs.Close()
    return err
}

func (d VolumeDriver) GetDatasetMountpoint(ds string) (string, error) {
    var v string
    fs, err := d.openDataset(ds)
    if err != nil {
        return v, err
    }
    if d.getProperty(&fs, OpenZFSMounted) == OpenZFSYes {
        return d.getProperty(&fs, OpenZFSMountpoint), err
    }
    defer fs.Close()
    return v, err
}

func (d VolumeDriver) MountDataset(ds string) error {
    fs, err := d.openDataset(ds)
    if err != nil {
        return err
    }
    if mounted, _ := fs.IsMounted(); mounted == true {
        return nil
    }
    if err := fs.Mount("",0); err != nil {
        return err
    }
    defer fs.Close()
    return err
}

func (d VolumeDriver) UnmountDataset(ds string) error {
    fs, err := d.openDataset(ds)
    if err != nil {
        return err
    }
    if d.getProperty(&fs, UsrPropMountID) != "" {
        return nil
    }

    if mounted, _ := fs.IsMounted(); mounted == false {
        return nil
    }
    if err := fs.Unmount(0); err != nil {
        return err
    }
    if prop := d.getProperty(&fs, UsrPropDestroyOnDetach); prop == OpenZFSTrue {
        if err = fs.DestroyRecursive(); err != nil {
            return err
        }
    }
    defer fs.Close()
    return err
}

func (d VolumeDriver) AppendMountId(ds string, id string) error {
    var ids []string
    fs, err := d.openDataset(ds)
    mid := d.getProperty(&fs, UsrPropMountID)
    if mid != "" {
        ids = strings.Split(mid, ",")
        for i := range ids {
            if ids[i] == id {
                return nil
            }
        }
    }
    ids = append(ids, id)
    mid = strings.Join(ids, ",")
    err = d.setProperty(&fs, UsrPropMountID, mid)
    defer fs.Close()
    return err
}

func (d VolumeDriver) YankMountId(ds string, id string) error {
    var ids []string
    fs, err := d.openDataset(ds)
    mid := d.getProperty(&fs, UsrPropMountID)
    if mid != "" {
        ids = strings.Split(mid, ",")
        for i := range ids {
            if ids[i] == id {
                ids = ids[:i+copy(ids[i:], ids[i+1:])] // Stackoverflow :]
                break
            }
        }
    }
    mid = strings.Join(ids, ",")
    err = d.setProperty(&fs, UsrPropMountID, mid)
    defer fs.Close()
    return err
}

func (d VolumeDriver) GetDataset(ds string) (map[string]string, map[string]interface{}, error) {
    v := make(map[string]string)
    s := make(map[string]interface{})
    fs, err := d.openDataset(ds)
    if err != nil {
        return v, s, err
    }
    name := d.getProperty(&fs, UsrPropName)
    if name != "" {
        v[DVName] = d.getProperty(&fs, UsrPropName)
        v[DVMountpoint], _ = d.GetDatasetMountpoint(ds)
        v[DVCreatedAt] = d.getProperty(&fs, UsrPropCreated)
        s, err = d.getStatus(&fs, s)
    }
    defer fs.Close()
    return v, s, err
}

// Private functions
func (d VolumeDriver) getStatus(fs *zfs.Dataset, s map[string]interface{}) (map[string]interface{}, error) {
    var err error
    snaps := make([]map[string]string, 0)
    for i := range fs.Children {
        ds := fs.Children[i]
        if strings.Contains(ds.Properties[zfs.DatasetPropName].Value, OpenZFSSnapDivider) {
            snap := make(map[string]string)
            snap[NVDSnapName] = strings.TrimPrefix(ds.Properties[zfs.DatasetPropName].Value,
                            fs.Properties[zfs.DatasetPropName].Value + OpenZFSSnapDivider)
            snap[NVDSnapTime] = fromUnixDsTime(d.getProperty(&ds, OpenZFSCreation))
            snaps = append(snaps, snap)
        }
    }
    if prop := d.getProperty(fs, UsrPropDescription); prop != "" {
        s[NVDDescription] = prop
    } else {
        s[NVDDescription] = NVDDescriptionDefault + d.getProperty(fs, UsrPropName)
    }
    s[NVDVolUsage] = getMiBFromOpenZFS(d.getProperty(fs, OpenZFSUsed))
    s[NVDVolSize] = getMiBFromOpenZFS(d.getProperty(fs, OpenZFSQuota))
    s[NVDGroup], _ = os.Hostname()
    s[NVDSnapshots] = snaps
    for k, v := range UsrPropToVanity {
        if prop := d.getProperty(fs, k); prop != "" {
            s[UsrPropToVanity[k]] = prop
        } else {
            if _, ok := s[UsrPropToVanity[k]]; ok == false {
                s[UsrPropToVanity[k]] = DefaultStatus[v]
            }
        }
    }
    return s, err
}

func (d VolumeDriver) quotaCalc(size string) string {
    v := getOpenZFSFromGiB(d.defaultSize)
    if (size != "") {
        return getOpenZFSFromGiB(size)
    }
    return v
}

func getOpenZFSFromGiB(b string) string {
    v, err := strconv.Atoi(b)
    if err != nil {
        return "0"
    }
    v = v * 1024 * 1024 * 1024
    return strconv.Itoa(v)
}

func getMiBFromOpenZFS(b string) int {
    v, err := strconv.Atoi(b)
    if err != nil {
        v = 0
    } else {
        v = v / 1024 / 1024
    }
    return v
}

func (d VolumeDriver) prefixDs(ds string) string {
    if strings.HasPrefix(ds, d.baseds) == false {
        ds = d.baseds + OpenZFSDivider + ds
    }
    return ds
}

func (d VolumeDriver) getProperty(fs *zfs.Dataset, p string) string {
    var v string
    var prop zfs.Property
    if strings.Contains(p, OpenZFSUserprop) {
        prop, _ = fs.GetUserProperty(d.prefix + p)
        if prop.Source == OpenZFSSourceLocal {
            v = prop.Value
        }
    } else {
        for i := range fs.Properties {
            if zfs.DatasetPropertyToName(i) == p {
                v = fs.Properties[i].Value
            }
        }
    }
    return v
}

func (d VolumeDriver) setProperty(fs *zfs.Dataset, p string, v string) error {
    if strings.Contains(p, OpenZFSUserprop) {
        if err := d.setUserProperty(fs, p, v); err != nil {
            return err
        }
    } else {
        for i := range fs.Properties {
            if zfs.DatasetPropertyToName(i) == p {
                if err := fs.SetProperty(i, v); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func (d VolumeDriver) setUserProperty(fs *zfs.Dataset, p string, v string) error {
    err := fs.SetUserProperty(d.prefix + p, v)
    return err
}

func (d VolumeDriver) volumeExist(fs *zfs.Dataset) error {
    var err error
    if d.getProperty(fs, UsrPropName) == "" {
        return errors.New("Dataset exist without proper metadata")
    }
    return err
}

func datasetExist(ds string) error {
    fs, err := zfs.DatasetOpen(ds)
    defer fs.Close()
    return err
}

func (d VolumeDriver) baseDatasetExist() error {
    fs, err := d.openDataset(d.baseds)
    if err != nil {
        return nil
    }
    if p,_ := fs.GetUserProperty(d.prefix + UsrPropBaseDS); p.Value != OpenZFSTrue {
        defer fs.Close()
        return errors.New(d.prefix + " on " + d.baseds + " is not true or does not exist. Exiting")
    }
    defer fs.Close()
    return nil
}

func (d VolumeDriver) openDataset(ds string) (zfs.Dataset, error) {
    fs, err := zfs.DatasetOpen(d.prefixDs(ds))
    if err != nil {
        return fs, err
    }
    if ds == d.baseds {
        return fs, err
    }
    if d.volumeExist(&fs) != nil {
        return fs, err
    }
    return fs, err
}

func (d VolumeDriver) importDataset(ds string, usrp map[string]string) error {
    var snap     zfs.Dataset
    var snapshot string
    fs, err := zfs.DatasetOpen(d.prefixDs(usrp[VDImportVol]))
    if err == nil {
        if mounted, _ := fs.IsMounted(); mounted == true {
            return errors.New("Dataset is mounted")
        }
        if usrp[VDRestore] == OpenZFSTrue {
            if _, ok := usrp[VDSnapshot]; ok == true {
                snapshot = d.prefixDs(usrp[VDImportVol]) + OpenZFSSnapDivider + usrp[VDSnapshot]
            } else {
                snapshot, err = d.getLastSnap(usrp[VDImportVol])
                if err != nil {
                    fs.Close()
                    return err
                }
            }
            snap, err = zfs.DatasetOpen(snapshot)
            if err != nil {
                fs.Close()
                return err
            }
            if err = fs.Rollback(&snap, true); err != nil {
                fs.Close()
                return err
            }
        }
        if err = fs.Rename(d.prefixDs(ds), false, true); err != nil {
            defer fs.Close()
            return err
        }
    }
    defer fs.Close()
    return err
}

func (d VolumeDriver) getLastSnap(ds string) (string, error) {
    fs, err := zfs.DatasetOpen(d.prefixDs(ds))
    var time_ts int
    if err != nil {
        return ds, err
    }
    i := 0
    for k := range fs.Children {
        if strings.Contains(d.getProperty(&fs.Children[k], OpenZFSPropName), OpenZFSSnapDivider) {
            time_ts, err = strconv.Atoi(d.getProperty(&fs.Children[k], OpenZFSCreation))
            if time_ts >= i {
                i = time_ts
                ds = fs.Children[k].Properties[zfs.DatasetPropName].Value
            }
        }
    }
    if i == 0 {
        err = errors.New("No snapshots found on " + ds)
    }
    return ds, err
}

func (d VolumeDriver) cloneDataset(ds string, usrp map[string]string, props map[zfs.Prop]zfs.Property) error {
    var createOp    string
    var err         error
    var fs          zfs.Dataset
    var snap        zfs.Dataset
    var snap_s      string
    if _, ok := usrp[VDImportVolAsClone]; ok {
        createOp = VDImportVolAsClone
        fs, err = zfs.DatasetOpen(d.prefixDs(usrp[createOp]))
    } else if _, ok := usrp[VDCloneOf]; ok {
        createOp = VDCloneOf
        fs, err = d.openDataset(usrp[createOp])
    }
    if err == nil {
        if _,  ok := usrp[VDSnapshot]; ok {
            if val, _ := usrp[VDCreateSnapshot]; val == OpenZFSTrue {
                snap_s = d.prefixDs(usrp[createOp] + OpenZFSSnapDivider + usrp[VDSnapshot])
                snap, err = zfs.DatasetSnapshot(snap_s, false, props)
            } else {
                snap_s = d.prefixDs(usrp[createOp] + OpenZFSSnapDivider + usrp[VDSnapshot])
                snap, err = zfs.DatasetOpen(snap_s)
            }
        } else {
            snap_s = d.prefixDs(usrp[createOp] + OpenZFSSnapDivider + timeString())
            snap, err = zfs.DatasetSnapshot(snap_s, false, props)
        }
        if err != nil {
            return errors.New(err.Error() + ": " + snap_s)
        }
        clone, err := snap.Clone(d.prefixDs(ds), props)
        if err != nil {
            return err
        }
        defer snap.Close()
        defer clone.Close()
    } else {
        defer fs.Close()
        return err
    }
    defer fs.Close()
    return err
}

func timeString() string {
    return time.Now().UTC().Format(time.RFC3339Nano)
}

func fromUnixDsTime(sec string) string {
    ts, err := strconv.Atoi(sec)
    if err != nil {
        return time.Unix(int64(0), int64(0)).UTC().Format(time.RFC3339Nano)
    }
    return time.Unix(int64(ts), int64(0)).UTC().Format(time.RFC3339Nano)
}

func (d VolumeDriver) runFsCmd (ds string, unixcmd string, mode string) error {
    err := d.MountDataset(ds)
    if err != nil {
        return err
    }
    cmd := exec.Command(unixcmd, mode, d.mountPoint + OpenZFSDivider + ds)
    err = cmd.Run()
    uerr := d.UnmountDataset(ds)
    if err != nil {
        if uerr != nil {
            err = errors.New(err.Error() + uerr.Error())
        }
        return err
    }
    return err
}

func getHelp(help string) string {
    if help == OpenZFSHelp {
        return DVPOpenZFSHelp
    }
    return DVPHelp
}
