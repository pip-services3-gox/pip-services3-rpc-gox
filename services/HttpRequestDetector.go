package services

import (
	"net/http"
	"regexp"
)

/*
HttpRequestDetector рelper class that retrieves parameters from HTTP requests.
*/
var HttpRequestDetector = THttpRequestDetector{}

type THttpRequestDetector struct {
}

// DetectPlatform method are detects the platform (using "user-agent") from which the given HTTP request was made.
//   -  req  *http.Request an HTTP request to process.
//   Returns the detected platform and version. Detectable platforms: "mobile", "iphone",
//   "ipad",  "macosx", "android",  "webos", "mac", "windows". Otherwise - "unknown" will
// be returned.
func (c *THttpRequestDetector) DetectPlatform(req *http.Request) string {
	ua := req.Header.Get("user-agent")
	var version string
	var pattern string

	pattern = "/mobile/i"
	match, _ := regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "mobile"
	}

	pattern = "/like Mac OS X/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		re := regexp.MustCompile(`/CPU( iPhone)? OS ([0-9\._]+) like Mac OS X/`)
		result := re.FindAllStringSubmatch(ua, -1)
		re = regexp.MustCompile("/_/g")
		version = re.ReplaceAllString(result[0][2], ".")

		pattern = "/iPhone/"
		match, _ = regexp.Match(pattern, ([]byte)(ua))
		if match {
			return "iphone " + version
		}
		pattern = "/iPad/"
		match, _ = regexp.Match(pattern, ([]byte)(ua))

		if match {
			return "ipad " + version
		}
		return "macosx " + version
	}

	pattern = "/Android/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		re := regexp.MustCompile(`/Android ([0-9\.]+)[\);]/`)
		version = re.FindAllStringSubmatch(ua, -1)[0][1]
		return "android " + version
	}

	pattern = `/webOS\//`
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		re := regexp.MustCompile(`/webOS\/([0-9\.]+)[\);]\`)
		version = re.FindAllStringSubmatch(ua, -1)[0][1]
		return "webos " + version
	}

	pattern = "/(Intel|PPC) Mac OS X/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		re := regexp.MustCompile(`/(Intel|PPC) Mac OS X ?([0-9\._]*)[\)\;]/`)
		result := re.FindAllStringSubmatch(ua, -1)
		re = regexp.MustCompile("/_/g")
		version = re.ReplaceAllString(result[0][2], ".")
		return "mac " + version
	}

	pattern = "/Windows NT/"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		//try {
		re := regexp.MustCompile(`/Windows NT ([0-9\._]+)[\);]/`)
		version = re.FindAllStringSubmatch(ua, -1)[0][1]
		return "windows " + version
		// }
		//catch (ex) {
		return "unknown"
		//}
	}
	return "unknown"
}

// Detects the browser (using "user-agent") from which the given HTTP request was made.
//   -  req  *http.Reques an HTTP request to process.
//   Returns the detected browser. Detectable browsers: "chrome", "msie", "firefox",
// "safari". Otherwise - "unknown" will be returned.
func (c *THttpRequestDetector) DetectBrowser(req *http.Request) string {

	ua := req.Header.Get("user-agent")

	var pattern string
	pattern = "/chrome/i"
	match, _ := regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "chrome"
	}

	pattern = "/msie/i"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "msie"
	}

	pattern = "/firefox/i"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "firefox"
	}

	pattern = "/safari/i"
	match, _ = regexp.Match(pattern, ([]byte)(ua))
	if match {
		return "safari"
	}

	if ua == "" {
		return "unknown"
	}
	return ua
}

// DetectAddress method are detects the IP address from which the given HTTP request was received.
//   -  req *http.Reques an HTTP request to process.
//   Returns the detected IP address (without a port). If no IP is detected -
// nil will be returned.
func (c *THttpRequestDetector) DetectAddress(req *http.Request) string {
	var ip string
	// TODO: need to write!!

	// if req.headers["x-forwarded-for"] {
	// 	ip = req.headers["x-forwarded-for"].split(",")[0]
	// }

	// if ip == nil && req.ip {
	// 	ip = req.ip
	// }

	// if ip == nil && req.connection {
	// 	ip = req.connection.remoteAddress
	// 	if !ip && req.connection.socket {
	// 		ip = req.connection.socket.remoteAddress
	// 	}
	// }

	// if ip == nil && req.socket {
	// 	ip = req.socket.remoteAddress
	// }

	// // Remove port
	// if ip != nil {
	// 	ip = ip.toString()
	// 	var index = ip.indexOf(":")
	// 	if index > 0 {
	// 		ip = ip.substring(0, index)
	// 	}
	// }

	return ip
}

// DetectServerHost method are detects the host name of the request"s destination server.
//   -  req *http.Request  an HTTP request to process.
//   Returns the destination server"s host name.
func (c *THttpRequestDetector) DetectServerHost(req *http.Request) string {
	//TODO: Need fix this
	return "" + req.URL.Hostname() // socket.localAddress
}

// DetectServerPort method are detects the request"s destination port number.
//   -  req  *http.Request an HTTP request to process.
//   Returns the detected port number or <code>80</code> (if none are detected).
func (c *THttpRequestDetector) DetectServerPort(req *http.Request) string {
	return req.URL.Port() //socket.localPort
}
