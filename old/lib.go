package main 

import (

    "net"
    "errors"
    "github.com/tidwall/gjson"
    "fmt"
    "net/http"
    "time"
    "io/ioutil"
    "strconv"
)


func find_master(){
	//stub
}

//range over reports

func closest_device() string{
	
	    max := 0
	    
	    max_id := ""
	    	      
		for _, rep := range reports {

			rssi,_ := strconv.Atoi(rep.Rssi)
	        			 
			 if(rssi > max && rssi>0){
				 
				 max = rssi
				 max_id = rep.Id
				 
			 }
			 
	   }
	   
	   return max_id 
	
	
	
}

func search_for_devices(){
	
	ticker := time.NewTicker(5 * time.Second)
	    
    go func() {
	    
	for range ticker.C {    
	
	dev := getDeviceSaved()
	
	
	node, err := getMacAddr()
		
	if err != nil {}
	
	for device ,_ := range dev{
		
		
		rssi := device_strength(device)
		
		_, err2 := http.Get("http://"+master+"/report/"+node+"/"+device+"/"+rssi)
		
	    if err2 != nil {
	        fmt.Printf("%s", err)
	        return
	        
	    }
		
	}
	
	}
	
	
	}()
	
	
	
}

func node_register(){
	
	
	ticker := time.NewTicker(30 * time.Second)
	    
    go func() {

	for range ticker.C {

		addr, err := externalIP()
	    
		if err != nil {}
	    
	    id, err := getMacAddr()
		
		if err != nil {}
		
		_, err2 := http.Get("http://"+master+"/register/"+id+"/"+addr+"/"+label)
		
	    if err2 != nil {
	        fmt.Printf("%s", err)
	        return
	        
	    }
    
    }
    
    }() 

	
}

func get_devices() { 
	
		if(mode=="master"){
			   
			   return //we already have devices local
			   
		}
	
	    ticker := time.NewTicker(10 * time.Second)

	    
	    go func() {

		for range ticker.C {
	

		
	
	
	    response, err := http.Get("http://"+master+"/status")
	
	    if err != nil {
	        fmt.Printf("%s", err)
	        return
	        
	    } 
    
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        
        if err != nil {
            fmt.Printf("%s", err)
            return 
        }
        
        fmt.Printf("%s\n", string(contents))
        
        //cleanup existing devices
		temp := getDeviceSaved()
		
		for k, _ := range temp{
			
			removeDevice(k)
			
		}
        
        result := gjson.Get(string(contents), "payload.devices")
		result.ForEach(func(key, value gjson.Result) bool {
			//fmt.Println(value.String()) 
			
		
			
			mac := gjson.Parse(value.String()).Get("id")
			label := gjson.Parse(value.String()).Get("label")
			action := gjson.Parse(value.String()).Get("action")
			
			if(mac.String()!=""){
			
			saveDevices(mac.String(),label.String(),action.String())
			
			}
			
			return true // keep iterating
		})
		
	 }

	}()
		
				    
	
	
}

func getMacAddr() (string, error) {
    ifas, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    
    for _, ifa := range ifas {
        a := ifa.HardwareAddr.String()
        if a != "" {
            return a, nil
        }
    }
    
    e2 := errors.New("Interface Not Found")
    
    return "",e2

}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}