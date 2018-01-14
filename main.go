package main

import (
	"flag"

	"github.com/golang/glog"
)

var (
	candidhost     string
	candiduser     string
	candidpassword string
	turbohost      string
	turbouser      string
	turbopassword  string
)

func parseFlags() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&candidhost, "candidhost", "https://192.168.130.21", "the address of turbo.server")
	flag.StringVar(&candiduser, "candiduser", "admin", "username to login to turbo.server")
	flag.StringVar(&candidpassword, "candidpass", "Cisco!123", "password to login to turbo.server")
	flag.StringVar(&turbohost, "turbohost", "https://192.168.130.31", "the address of turbo.server")
	flag.StringVar(&turbouser, "turbouser", "administrator", "username to login to turbo.server")
	flag.StringVar(&turbopassword, "turbopass", "Cisco!123", "password to login to turbo.server")
	flag.Parse()
}

func getCandidInterfaces() ([]string, error) {
	var interfaceDNs []string
	client, _ := NewCandidRestClient(candidhost, candiduser, candidpassword)
	fabricXML, _ := client.getFabric()
	fabricURLs, _ := parseFabric(fabricXML)
	for _, fabricElement := range fabricURLs {
		var smartEventsIDs []string
		smartEventXML, _ := client.getSmartEvents(fabricElement)
		smartEventURLs, _ := parseEvents(smartEventXML, fabricElement)
		for _, eventElement := range smartEventURLs {
			smartEventsIDs = append(smartEventsIDs, eventElement)
		}
		for _, eventElement := range smartEventsIDs {
			eventDetailXML, _ := client.getEventDetail(eventElement)
			interfaceDn, _ := parseEventDetail(eventDetailXML)
			interfaceDNs = append(interfaceDNs, interfaceDn)
		}
	}

	glog.V(2).Infof("result=%+v", interfaceDNs)
	return interfaceDNs, (nil)
}

func updateTurboGroup(interfaceDNs []string) {
	client, _ := NewTurboRestClient(turbohost, turbouser, turbopassword)
	providerXML, _ := client.getGroup("PMs_CandidLinkDown")
	providerGroupUUID, _ := parseGroupUUID(providerXML)
	providerGroup, _ := client.updateGroup(providerGroupUUID, interfaceDNs)
	glog.V(2).Infof("result=%+v", string(providerGroup))
}

func main() {
	parseFlags()
	interfaceDNs, _ := getCandidInterfaces()
	updateTurboGroup(interfaceDNs)
}
