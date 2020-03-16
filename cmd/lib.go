package main 

import (
 
    "net"
    "log"
    "time"
    "net/http"
    "io/ioutil"
)

func resolve_master() string{
	
    min := 0
    resp := ""
    

	for ip, val := range peer_list {
		
		    //fmt.Println("IP:",ip," INT:",val)
		
	        if (val.Uptime > min) {
	            min = val.Uptime
	            resp = ip
	        }
	}

    master_host = resp
    
    if(master_host==getlocalip()){
	    
	    master = true;
	    slave = false;
	    
	    log.Println("I AM THE MASTER: ",master_host)
	    
    }else{
	    
	     log.Println("I AM A SLAVE")
	    //retrieve the devices?
	    master = false;
	    slave = true;
    }
	
	return resp
	
	
}


func getoutboundip() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
            
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP.String()
}

func getlocalip() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}


func uptime() int {
    return int(time.Since(startTime).Seconds())
}

func httpGet(url string) string{
	
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		
	}

	req.Header.Set("User-Agent", "device-client")

	res, getErr := client.Do(req)
	if getErr != nil {
		
	}

	body, _ := ioutil.ReadAll(res.Body)
	
	return string(body)
	
	
}