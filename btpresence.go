package bt


import (
	"os/exec"
	"strings"
	"log"
)

func init(){
    
     log.SetFlags(log.LstdFlags | log.Lshortfile)
    
}

func Rssi(id string) string{

    cmd :="hcitool"
	args := []string{"-i","hci0","cc",id}
	_, err := exec.Command(cmd,args...).Output()
    if err!=nil {
	    
	    if (strings.Contains(err.Error(),"executable file not found")){
		   return "-1000" 
	    }
	    
        log.Println(err.Error())
        return "-1000"
    }
    
	args2 := []string{"-i","hci0","rssi",id}
	output2, err2 := exec.Command(cmd,args2...).Output()
    if err2!=nil {
        log.Println(err2.Error())
        return "-1000"
    }
    
    args3 := []string{"-i","hci0","dc",id,"19"} //19 is the reason - user disconnect
    	
	_, err3 := exec.Command(cmd,args3...).Output()
    if err3!=nil {
        log.Println(err3.Error())
        
    }
    
    //fmt.Println(string(output2))
		
	op := string(output2)
	


	stat := strings.Contains(op,"RSSI return value")
	
	
	if(stat==true){
		
		r1 := strings.Replace(op,"RSSI return value: ","",-1)
		r2 := strings.Trim(r1,"%0A")
		r3 := strings.Replace(r2,"\n","",-1)
		
		return r3
		
	}
	
    return "-1000"
    
 }   