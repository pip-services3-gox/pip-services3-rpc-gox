package services

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
)

type AboutOperations struct {
	RestOperations
	contextInfo *cinfo.ContextInfo
}

func (c *AboutOperations) SetReferences(references crefer.IReferences) {
	c.RestOperations.SetReferences(references)

	depResult := references.GetOneOptional(crefer.NewDescriptor("pip-services", "context-info", "*", "*", "*"))
	if depResult != nil {
		c.contextInfo = depResult.(*cinfo.ContextInfo)
	}

}

func (c *AboutOperations) GetAboutOperation() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		c.About(res, req)
	}
}

func (c *AboutOperations) GetNetworkAddresses() []string {
	interfacesAddrs, _ := net.InterfaceAddrs()
	var addresses []string
	for _, address := range interfacesAddrs {
		ipAddr, ipOk := address.(*net.IPAddr)
		if !ipOk {
			continue
		}
		if ipAddr.IP.IsGlobalUnicast() && ipAddr.Network() == "ip4" {
			addresses = append(addresses, address.String())
		}

	}
	return addresses
}

func (c *AboutOperations) About(res http.ResponseWriter, req *http.Request) {
	about := make(map[string]interface{}, 0)
	server := make(map[string]interface{})

	server["name"] = "unknown"
	server["description"] = ""
	server["properties"] = ""
	server["uptime"] = ""
	server["start_time"] = ""
	if c.contextInfo != nil {
		server["name"] = c.contextInfo.Name
		server["description"] = c.contextInfo.Description
		server["properties"] = c.contextInfo.Properties
		server["uptime"] = c.contextInfo.Uptime()
		server["start_time"] = c.contextInfo.StartTime
	}

	server["current_time"] = time.Now().Format(time.RFC3339)
	server["protocol"] = req.URL.Scheme
	server["host"] = HttpRequestDetector.DetectServerHost(req)
	server["addresses"] = c.GetNetworkAddresses()
	server["port"] = HttpRequestDetector.DetectServerPort(req)
	server["url"] = req.URL.String()

	about["server"] = server

	client := make(map[string]interface{})
	client["address"] = HttpRequestDetector.DetectAddress(req)
	client["client"] = HttpRequestDetector.DetectBrowser(req)
	client["platform"] = HttpRequestDetector.DetectPlatform(req)
	client["user"] = req.URL.User

	about["client"] = client

	jsonObj, jsonErr := json.Marshal(about)
	if jsonErr == nil {
		io.WriteString(res, (string)(jsonObj))
	}
	//res.json(about)
}
