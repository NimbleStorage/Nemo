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
    "os/user"
    "strconv"
    "github.com/urfave/cli"
    "github.com/docker/go-plugins-helpers/volume"
    log "github.com/Sirupsen/logrus"
)

func main () {
    app := cli.NewApp()
    app.Name = "Nemo"
    app.Usage = "Emulates HPE Nimble Storage Docker Volume Plugin"
    app.Author = "Michael Mattsson <michael.mattsson@hpe.com>"
    app.Version = "1.0.0"

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   "baseds, b",
            Usage:  "Base OpenZFS dataset (will be created)",
            Value:  "tank/default",
            EnvVar: "NEMO_BASEDS",
        },
        cli.StringFlag{
            Name:   "mountpath, m",
            Usage:  "Base mountpoint (on base dataset creation)",
            Value:  "/opt/nemo",
            EnvVar: "NEMO_MOUNTPATH",
        },
        cli.StringFlag{
            Name:   "prefix, p",
            Usage:  "OpenZFS properties prefix",
            Value:  "com.hpe.dev.nemo",
            EnvVar: "NEMO_PREFIX",
        },
        cli.StringFlag{
            Name:   "size, s",
            Usage:  "Default volume size in GiB (0=unlimited)",
            Value:  "0",
            EnvVar: "NEMO_DEFAULTSIZE",
        },
        cli.StringFlag{
            Name:   "destroy",
            Usage:  "Recursively destroy datasets on rm",
            Value:  "false",
            EnvVar: "NEMO_DESTROY",
        },
        cli.StringFlag{
            Name:   "name, n",
            Usage:  "Driver name",
            Value:  "nemo",
            EnvVar: "NEMO_DRIVERNAME",
        },
        cli.StringFlag{
            Name:   "user, u",
            Usage:  "Socket owner",
            Value:  "root",
            EnvVar: "NEMO_SOCKETOWNER",
        },
        cli.StringFlag{
            Name:   "debug, d",
            Usage:  "Sets log level (debug, info, warn, error)",
            Value:  "error",
            EnvVar: "NEMO_DEBUG",
        },
    }

    app.Action = func(c *cli.Context) error {
        switch c.String("debug") {
            case "debug":
                log.SetLevel(log.DebugLevel)
            case "info":
                log.SetLevel(log.InfoLevel)
            case "warn":
                log.SetLevel(log.WarnLevel)
            case "error":
                log.SetLevel(log.ErrorLevel)
            default:
                log.SetLevel(log.ErrorLevel)
        }

        d := newVolumeDriver(c)
        h := volume.NewHandler(d)
        u, _ := user.Lookup(d.socketOwner)
        gid, _ := strconv.Atoi(u.Gid)
        log.Debug("About to serve: ", d.driverName)
        h.ServeUnix(d.driverName, gid)
		return nil
	}

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
