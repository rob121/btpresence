# Bluetooth Presence 

This code allows you to setup bluetooth nodes on your network that detect which device you are closer to

## Building 

This is used to build for a raspberry pi, optionally use the util/build file which contains this command
```
env GOOS=linux GOARCH=arm GOARM=5 go build -o devicepresence
```

## Install 

There is an install script in util that will copy everything where it needs to go, for a linux/raspi 

## Web Interface

By default there is a web interface exposed on port 15784, that lets you see nodes, etc.

Access it by going to

http://ip:15784/manage


## Operating Rationale

Each node looks for others on the network, through a simple conseus process they elect a master and then send all the location data to that node, that node then connects to whatever automation system is configured in drivers.

The system uses hcitool and at least on the raspberry pi had to be run as root in order to access the hardware.

## Drivers

Drivers can be added to the drivers folder and you must create two methods 

Arrived(mac string,room string)
Departed(mac string,room string)

## Running Client

On the device use the bin/devicepresence binary (or build your own) and run as follows

./devicepresence --room=name --devices=./devices

### Flags

* --room (string): This is a identifier for the room name, doesn't matter so long as it's meaningful to you
* --devices (string, file path): Full file path and file name to a file that has a list of bluetooth mac's to scan for, one per line

## Client Configuration

There are possibly two types of configuration necessary 

1) Device list, in the root name a file "devices" (or configure with flag) with a list of mac addresses to monitor for, one per line
2) Whatever specific driver configuration is necessary
