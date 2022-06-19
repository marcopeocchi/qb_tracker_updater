# qb_tracker_updater

Update your qBittorrent.conf trackers list automatically!

## Prerequisites

You must enable the **Automatically add these trackers to new downloads** config  
![prerequisite](https://i.ibb.co/jfjtzDc/image.png)  

or add manually, if not already specified
```
[Preferences]
...
Bittorrent\TrackersList=     <- this line under preferences section
...
```

## Usage
```shell
tracker_updater -conf <qBittorent.conf path> -profile <tracker list type>
```

## Installing binary version
```shell
$ mv qb_tracker_updater-[arch] /usr/local/bin/qb-tracker-updater
$ qb-tracker-updater --help
```

## Profiles
1. best 
2. best_ip [default]
3. all 
4. all_http 
5. all_https
6. all_ws
7. all_ip

## Systemd integration

This program can run seamlessly before qBittorrent stars, just add **ExecStartPre** options to your systemd service file

```
...
User= user
ExecStartPre=/usr/local/bin/qb-tracker-updater
# /usr/local/bin/qb-tracker-updater -conf [to manually specify] -profile [profile number]
ExecStart=/usr/bin/qbittorrent-nox
...
```