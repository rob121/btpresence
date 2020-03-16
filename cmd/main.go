package main 

import(
"flag"
"log"
"fmt"
"time"
"github.com/spf13/viper"
"github.com/spf13/pflag"
"strconv"
"github.com/rob121/btpresence"

)


var room string
var startTime time.Time
var peers chan Host
var peer_list map[string]Host
var master_host string
var master bool
var slave bool
var minrssi int
var checker map[string]time.Time
 
func main(){
     
    flag.String("room","default", "Room Name")	
    flag.Int("minrssi",-5, "Minimum Valid RSSI")	
    
    pflag.CommandLine.AddGoFlagSet(flag.CommandLine) 
    pflag.Parse()
    viper.BindPFlags(pflag.CommandLine)
    
    config()
    
    room = viper.GetString("room")
    minrssi = viper.GetInt("minrssi")
 
    log.Println("Room:",room)
 
    startTime = time.Now() 
     
     // scan for devices
     
    peers  = make(chan Host) 
    peer_list = make(map[string]Host)   
    checker = make(map[string]time.Time)
    
    log.Println(viper.Get("devices"))  
 
    go discover_peers()
    
    // start a web server here
    
    go search_devices()
    
    select{} 
     
     
     
 }
 
 func config(){
     
     
    viper.SetConfigName("config") // name of config file (without extension)
    viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name
    viper.AddConfigPath("/etc/btpresence/")   // path to look for the config file in
    viper.AddConfigPath("$HOME/.btpresence")  // call multiple times to add many search paths
    viper.AddConfigPath(".")               // optionally look for config in the working directory
    err := viper.ReadInConfig() // Find and read the config file
    if err != nil { // Handle errors reading the config file
    	log.Print("Config file error: %e", err)
    }
         
     
 }
 
 func search_devices(){
     
     for {
         
         for k,v := range viper.GetStringMap("devices") {
             
             log.Println("Searching for",k,v)
             rssi := bt.Rssi(k)
             
             log.Println("Rssi",rssi)
             
             
             sig,_ := strconv.Atoi(rssi)
             
             
             if(sig>minrssi){
                 
               dev := viper.GetString("devices."+k+".hubitat_device")  
               
               str := viper.GetString("hubitat_update")
               
               url := fmt.Sprintf(str,dev,room)
               
               log.Println(url)
               
               if tm,ok := checker[k]; ok { 
                   
                   if(tm.Unix() > time.Now().Add(-30 * time.Second).Unix()){
                      log.Println("Update made for "+k+" in the last 30 seconds, skipping")
                      continue 
                       
                   }
                   
               }
                 
               out:=httpGet(string(url))
               log.Println(out) 
               
               checker[k] = time.Now()
             
             }
             
             
         }
         
         time.Sleep(1 * time.Second)
         
     }
     
     
 }