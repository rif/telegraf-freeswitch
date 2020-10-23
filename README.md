# Telegraf Freeswitch Plugin
[![Build Status](https://travis-ci.org/rif/telegraf-freeswitch.svg?branch=master)](https://travis-ci.org/rif/telegraf-freeswitch)

Collects status data from Freeswitch and makes it available for telegraf.

The collected data include:

- Sessions
  - active sessions
  - peak and peak5min
  - total
- Rate
  - current
  - max
  - peak and peak5min
- Running sessions per Sofia Sip Profiles


## Install
telegraf-freeswitch is a stand-alone binary with no dependencies. Just copy it on your system and run it.

Use [released binaries](https://github.com/rif/telegraf-freeswitch/releases) or ```go get -u github.com/rif/telegraf-freeswitch```

## Telegraf configuration

There are two operation modes for telegraf-freeswitch: one shot and server.

In one shot telegraf will start the plugin process which will connect to freeswitch via eventsocket get the status and profiles information and exit.

In server mode the plugin is started externally, connects to freeswitch and stays connected responding to http GET requestd from telegraf. This server mode is slightly more complicated to set up but it might be more efficient then the one shot mode.

Basically the server mode replaces the starting of the plugin process and freeswitch connection by an http GET request. However there are no measurements of how much of an optimization this is.

## One shot mode

```toml
[[inputs.exec]]
  ## Commands array
  commands = ["/usr/local/bin/telegraf-freeswitch -host 127.0.0.1 -port 8021 -pass ClueCon"]

  ## Timeout for each command to complete.
  timeout = "5s"

  ## Data format to consume.
  ## Each data format has it's own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
```

###Example Output
```
$ ./telegraf-freeswitch
freeswitch_sessions active=1636,peak=2380,peak_5min=1740,total=1377928,rate_current=10,rate_max=300,rate_peak=234,rate_peak_5min=25
freeswitch_profile_sessions,profile=dot51,ip=80.161.218.51 running=0
freeswitch_profile_sessions,profile=dot48,ip=80.161.218.48 running=147
freeswitch_profile_sessions,profile=dot52,ip=80.161.218.52 running=0
freeswitch_profile_sessions,profile=dot47,ip=80.161.218.47 running=177
freeswitch_profile_sessions,profile=dot49,ip=80.161.218.49 running=169
freeswitch_profile_sessions,profile=external,ip=80.161.218.17 running=988
freeswitch_profile_sessions,profile=dot50,ip=80.161.218.50 running=155
```

## Execd mode

``` toml
[[inputs.execd]]
#   ## Program to run as daemon
command = ["/usr/local/bin/telegraf-freeswitch", "-execd"]
#
#   ## Define how the process is signaled on each collection interval.
#   ## Valid values are:
#   ##   "none"   : Do not signal anything.
#   ##              The process must output metrics by itself.
#   ##   "STDIN"   : Send a newline on STDIN.
#   ##   "SIGHUP"  : Send a HUP signal. Not available on Windows.
#   ##   "SIGUSR1" : Send a USR1 signal. Not available on Windows.
#   ##   "SIGUSR2" : Send a USR2 signal. Not available on Windows.
signal = "STDIN"
#
#   ## Delay before the process is restarted after an unexpected termination
#   restart_delay = "10s"
#
#   ## Data format to consume.
#   ## Each data format has its own unique set of configuration options, read
#   ## more about them here:
#   ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
data_format = "influx"
```

###Example Output
```
$ ./telegraf-freeswitch -execd

freeswitch_sessions active=1636,peak=2380,peak_5min=1740,total=1377928,rate_current=10,rate_max=300,rate_peak=234,rate_peak_5min=25
freeswitch_profile_sessions,profile=dot51,ip=80.161.218.51 running=0
freeswitch_profile_sessions,profile=dot48,ip=80.161.218.48 running=147
freeswitch_profile_sessions,profile=dot52,ip=80.161.218.52 running=0
freeswitch_profile_sessions,profile=dot47,ip=80.161.218.47 running=177
freeswitch_profile_sessions,profile=dot49,ip=80.161.218.49 running=169
freeswitch_profile_sessions,profile=external,ip=80.161.218.17 running=988
freeswitch_profile_sessions,profile=dot50,ip=80.161.218.50 running=155
```

## Server mode

```toml
## Read flattened metrics from one or more JSON HTTP endpoints
[[inputs.httpjson]]
name_override = "freeswitch_sessions"
## URL of each server in the service's cluster
servers = [
  "http://localhost:9191/status/",
]


[[inputs.httpjson]]
name_override = "freeswitch_profiles_sessions"
## URL of each server in the service's cluster
servers = [
  "http://localhost:9191/profiles/",
]
## List of tag names to extract from top-level of JSON server response
tag_keys = [
  "profile",
  "ip"
]
```

Copy telegraf-freeswitch.service in /etc/systemd/system/ folder and run ```systemctl daemon-reload``` command to load the newly added file.

After that use the usual systemctl start/stop/restart telegraf-freeswitch.service commands to controll the telegraf-freeswitch server.

###Example Output
```
$ curl http://localhost:9191/status/
{
 "active": 53,
 "peak": 54,
 "peak_5min": 54,
 "total": 114,
 "rate_current": 3,
 "rate_max": 300,
 "rate_peak": 3,
 "rate_peak_5min": 3
}

$ curl http://localhost:9191/profiles/
[
 {
  "name": "dot3",
  "ip": "80.161.218.3",
  "running": "19"
 },
 {
  "name": "dot4",
  "ip": "80.161.218.4",
  "running": "10"
 },
 {
  "name": "external",
  "ip": "80.161.218.2",
  "running": "15"
 },
 {
  "name": "dot5",
  "ip": "80.161.218.5",
  "running": "14"
 },
 {
  "name": "dot6",
  "ip": "80.161.218.6",
  "running": "14"
 }
]
```

### Similar plugins
- [Telegraf plugin for FreeSWITCH ](https://github.com/areski/freeswitch-telegraf-plugin)
- [FreeSWITCH Metric Collection with Telegraf](https://github.com/moises-silva/freeswitch-telegraf)
