# qb_tracker_updater

Update your qBittorrent.conf trackers list automatically!

## Prerequisites

You must enable the **Automatically add these trackers to new downloads** config  
![prerequisite](https://i.ibb.co/jfjtzDc/image.png)

## Usage
```shell
tracker_updater <qBittorent.conf path> [profile]
```
or  
```shell
node ./main.js <qBittorent.conf path> [profile]
```

## Installing binary version
```shell
# mv qb_tracker_updater-[arch] /usr/local/bin/qb-tracker-updater
$ qb-tracker-updater --help
```

## Profiles
- best 
- best_ip [default]
- all 
- all_http 
- all_https
- all_ws

## Systemd integration

This program can run seamlessly before qBittorrent stars, just add **ExecStartPre** options to your systemd service file

```
...
User= user
ExecStartPre=/usr/local/bin/node /home/user/qb_tracker_updater/main.js /home/user/.config/qBittorrent/qBittorrent.conf
ExecStart=/usr/bin/qbittorrent-nox
...
```

or 

```
...
User= user
ExecStartPre=/usr/local/bin/node /usr/local/bin/qb-tracker-updater main.js /home/user/.config/qBittorrent/qBittorrent.conf
ExecStart=/usr/bin/qbittorrent-nox
...
```