package main

import(
"flag"
//"github.com/rob121/btpresence"
)


var room string
var port int
var mode string
var settings string
var nodes map[string]Node
var devices map[string]Device
var master string
var label string
var reports map[string]DeviceState

type JsonResp struct {
	Code    int         `json:"code"`
	Payload interface{} `json:"payload"`
	Message string      `json:"message"`
}


type Node struct{
	Id string `json:"id"`	
	Addr string `json:"addr"`
	Label string `json:"label"`
	Status string `json:"status"`
}

type Device struct{
	Id string `json:"id"`
	Label string `json:"label"`
	Action string `json:"action"`	
}


type DeviceState struct {
	Rssi string `json:"rssi"`
	Reporter string `json:"reporter"` //mac of node
	Id string `json:"id"` //mac of device 
	Reported int64 `json:"reported"`
}


func main() {
	
	nodes = make(map[string]Node)
	reports = make(map[string]DeviceState)
	
	flag.StringVar(&label, "label","default", "Room Name")		
	flag.IntVar(&port,"port",15784,"Port to listen on")
	flag.StringVar(&mode,"mode","master","Master or Slave?") 
	flag.StringVar(&master,"master","","Master ip address/port") 
	flag.StringVar(&settings,"settings","./","File path to save settings") 
    flag.Parse()
    
    //add yourself as a node!
    
    addr, err := externalIP()
    
	if err != nil {}
    
    id, err := getMacAddr()
	
	if err != nil {}
    
   
	//add the node info to this device if master, otherwise find and send to master
	
	if(mode=="master"){
		
	 nodes[id] = Node{Id: id,Label: label,Addr: addr,Status: mode}	
	 go start_server()
	
	}
	
      //register myself as a node
      
    go node_register() //register node every 30 seconds in case master falls off etc.
	
	go find_master() //find the master on the network, for now just use the string for master location
	
	go get_devices() //get the list from the master, if master then you have the list in "getDevicesSaved"
	
	go search_for_devices()
	
	
	
	
	
	select{}
}
