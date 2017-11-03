# Telegraf Freeswitch Plugin

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

Use [releases binaries](https://github.com/rif/telegraf-freeswitch/releases) or ```go get -u github.com/rif/telegraf-freeswitch```

## Telegraf configuration

```
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

### Similar plugins
- [Telegraf plugin for FreeSWITCH ](https://github.com/areski/freeswitch-telegraf-plugin)
- [FreeSWITCH Metric Collection with Telegraf](https://github.com/moises-silva/freeswitch-telegraf)
