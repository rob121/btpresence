package main 


import(
	"net/http"
    "log"
    "fmt"
    "github.com/GeertJohan/go.rice"
    "text/template"
    "encoding/json"
    "bytes"
    "os"
    "io/ioutil"
    "encoding/gob"
    "path/filepath"
    "strings"
    "time"
   // "strconv"
)



func start_server(){
	
	
	log.Print("Listening on port ", port)
	
	

	box := rice.MustFindBox("www")
	assetsFileServer := http.StripPrefix("/assets/", http.FileServer(box.HTTPBox()))
	http.Handle("/assets/", assetsFileServer)
	
	if(mode=="master"){
	//node reports its findings for a device
	http.HandleFunc("/report/", reportHandler)
    //register a node 
    http.HandleFunc("/register/", registerHandler)
    //list devices
	http.HandleFunc("/devices/", deviceHandler)
	
	//save a device
	http.HandleFunc("/devicesave/", deviceSaveHandler)
	
	//remove a device
	http.HandleFunc("/deviceremove/", deviceRemoveHandler)
    }
    
    http.HandleFunc("/status/", statusHandler)
	
	http.HandleFunc("/", defaultHandler)
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	
}

func deviceRemoveHandler(w http.ResponseWriter, r *http.Request) { 
	
    w.Header().Set("Content-type", "application/json")

	r.ParseForm()
	
	path := r.URL.Path

	if r.Method != "POST" {
		respond(w, 500, "Invalid Request - must POST", "")
		return

	}
	
	
	path = strings.Replace(path, " ", "_", -1)
	path = strings.ToLower(path)
	
	parts := strings.Split(path, "/")
	
	if parts[2] == "" {
		
		respond(w,500,"No Device Provided","")
		
		
	}
	
	mac := parts[2]

   
	
	removeDevice(mac)
	
	dev := getDeviceSaved();
	
	respond(w,200,"Device Removed",dev)
	
	
}

func deviceSaveHandler(w http.ResponseWriter, r *http.Request) { 
	
    w.Header().Set("Content-type", "application/json")

	r.ParseForm()

	if r.Method != "POST" {
		respond(w, 500, "Invalid Request - must POST", "")
		return

	}

	mac := r.FormValue("mac")
	
	label := r.FormValue("label")
	
	action := r.FormValue("on_detect_http")
	
	saveDevices(mac,label,action)
	
	dev := getDeviceSaved();
	
	respond(w,200,"Good Report",dev)
	
	
}

func registerHandler(w http.ResponseWriter, r *http.Request) { 
	
    w.Header().Set("Content-type", "application/json")

	r.ParseForm()
	
	path := r.URL.Path

	
	
	path = strings.Replace(path, " ", "_", -1)
	path = strings.ToLower(path)
	
	parts := strings.Split(path, "/")
	
	if len(parts) < 5 {
		
		respond(w,500,"Missing Values","")
		return
		
	}
	
	
	id := parts[2]
	
	addr := parts[3]
	
	label := parts[4]


     nodes[id]=Node{Id:id,Label: label,Addr: addr,Status:"slave"} //by defintion it has to be a slave,since only the master gets this update 
	
	
	
	respond(w,200,"Node Added",nodes)
	
}

func reportHandler(w http.ResponseWriter, r *http.Request) { 
	
	
	//report is in the format
	
	//report/{node}/{device}/{rssi}
	
	r.ParseForm()
	
	path := r.URL.Path


	path = strings.Replace(path, " ", "_", -1)
	path = strings.ToLower(path)
	
	parts := strings.Split(path, "/")
	
	if len(parts)<5 {
		
		respond(w,500,"No Node ID Provided","")
		
		
	}

	
	now := time.Now()
    tmnow := now.Unix()
	
	node := parts[2]
	device := parts[3]
	rssi := parts[4]
	
	reports[node+device]=DeviceState{Rssi: rssi,Reporter: node,Id: device,Reported: tmnow}

	
	respond(w,200,"Good Report",reports)
	
}

func statusHandler(w http.ResponseWriter, r *http.Request) { 
	
		
	dev := getDeviceSaved();
	
	output := make(map[string]interface{})
	
	output["devices"] = dev
	
	output["nodes"] = nodes
	
	output["reports"] = reports
	
	respond(w,200,"Good Status",output)
	
}

func defaultHandler(w http.ResponseWriter, r *http.Request) { 
	
	var repl = make(map[string]string)
	
	repl["Mode"]="manual"
	
	respond_tmpl(w,"tmpl/index.html",repl)
	
}

func deviceHandler(w http.ResponseWriter, r *http.Request) { 
	
	var repl = make(map[string]string)
	
	repl["Mode"]="manual"
	
	respond_tmpl(w,"tmpl/device.html",repl)
	
}

func removeDevice(mac string){
	
	dev := getDeviceSaved();
	
	delete(dev,mac);

	
	fp := filepath.FromSlash(settings+"devices.gob")
	
	file, err := os.Create(fp)
   
    if err == nil { 
       
       
    }
        
    encoder := gob.NewEncoder(file)
     
    if err := encoder.Encode(dev); err != nil {
		
	}
	
	file.Close()
	
}

func getDeviceSaved() map[string]Device{
	
		// Create a file for IO
		
	fp := filepath.FromSlash(settings+"devices.gob")	
	
	byt, err := ioutil.ReadFile(fp)
	
	encodeFile := bytes.NewReader(byt)
	
	if err != nil {
	
	}


	
	decoder := gob.NewDecoder(encodeFile)
	
		// Place to decode into
	out := make(map[string]Device)

	// Decode -- We need to pass a pointer otherwise accounts2 isn't modified
	decoder.Decode(&out)
	

	
	return out
	
}

func saveDevices(mac string,label string,action string){
	
	
	dev := getDeviceSaved();

	
	dev[mac]=Device{Id: mac,Label: label,Action: action}

	fp := filepath.FromSlash(settings+"devices.gob")	 
	 
	file, err := os.Create(fp)
   
    if err == nil { 
       
       
    }
        
    encoder := gob.NewEncoder(file)
     
    if err := encoder.Encode(dev); err != nil {
		
	}
	
	file.Close()
	
	 
}



func tmpl_setup(templ string) (t *template.Template,err error) { 
	
		templateBox, err := rice.FindBox("www")
		
		if err != nil {
		 log.Println(err)
		 return nil,err
	    }
	    
		templateString, err := templateBox.String(templ)
		
		if err != nil {
			log.Println(err)
		    return nil,err
		}
		// parse and execute the template
		tmplMessage, err := template.New("message").Parse(templateString)
		
		if err != nil {
		log.Println(err)
	    return nil,err
	    }
	    
	    return tmplMessage,nil

	
}

func respond_tmpl(w http.ResponseWriter,tmpl string,replace map[string]string){

   t,err := tmpl_setup(tmpl)
   
   if(err!=nil){ 
	   
	   
   }
   
   t.Execute(w, replace)
	 
}


func respond(w http.ResponseWriter, code int, message string, payload interface{}) {

	resp := JsonResp{
		Code:    code,
		Payload: payload,
		Message: message,
	}

	var jsonData []byte
	jsonData, err := json.Marshal(resp)

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	fmt.Fprintln(w, string(jsonData))

}