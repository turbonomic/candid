package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/glog"
)

type CandidRestClient struct {
	client   *http.Client
	host     string
	username string
	password string
}

func NewCandidRestClient(host, uname, pass string) (*CandidRestClient, error) {

	//1. get http client
	client := &http.Client{
		Timeout: defaultTimeOut,
	}

	//2. check whether it is using ssl
	addr, err := url.Parse(host)
	if err != nil {
		glog.Errorf("Invalid url:%v, %v", host, err)
		return nil, err
	}
	if addr.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &CandidRestClient{
		client:   client,
		host:     host,
		username: uname,
		password: pass,
	}, nil
}

func (c *CandidRestClient) getRequest(URL string) ([]byte, error) {
	//1. a new http request
	urlStr := fmt.Sprintf("%s%s", c.host, URL)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return nil, err
	}

	//2. set queries
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	//3. set headers
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")

	//4. add user/password
	req.SetBasicAuth(c.username, c.password)

	//5. send request and receive result
	resp, err := c.client.Do(req)
	if err != nil {
		glog.Errorf("Failed to send request: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read response: %v", err)
		return nil, err
	}

	glog.V(4).Infof("resp: %++v", resp)
	return result, nil
}

type FabricResponse struct {
	ResponseBean struct {
		Value struct {
			Data []struct {
				UUID                  string `xml:"uuid"`
				SystemObject          string `xml:"systemObject"`
				Interval              int    `xml:"interval"`
				AnalysisTimeoutInSecs string `xml:"analysisTimeoutInSecs"`
				Status                string `xml:"status"`
				Links                 struct {
					Links struct {
						Rel  string `xml:"rel"`
						Href string `xml:"href"`
					} `xml:"links"`
				} `xml:"links"`
				UniqueName         string `xml:"unique_name"`
				AssuredNetworkType string `xml:"assured_network_type"`
				ApicHostnames      struct {
					ApicHostnames []string `xml:"apic_hostnames"`
				} `xml:"apic_hostnames"`
				Username        string `xml:"username"`
				ApplicationID   string `xml:"application_id"`
				Active          string `xml:"active"`
				OperationalMode string `xml:"operational_mode"`
			} `xml:"data"`
			DataSummary struct {
				Links struct {
					Links struct {
						Rel  string `xml:"rel"`
						Href string `xml:"href"`
					} `xml:"links"`
				} `xml:"links"`
				TotalCount  int  `xml:"total_count"`
				HasMoreData bool `xml:"has_more_data"`
			} `xml:"data_summary"`
		} `xml:"value"`
	} `xml:"ResponseBean"`
}

type SmartEventsResponse struct {
	ResponseBean struct {
		Value struct {
			Data []struct {
				UUID    string `xml:"uuid"`
				Summary struct {
					Leafs struct {
						Leafs struct {
							Name string `xml:"name"`
						} `xml:"leafs"`
					} `xml:"leafs"`
					Tenants struct {
						Tenants struct {
							Dn   string `xml:"dn"`
							Name string `xml:"name"`
						} `xml:"tenants"`
					} `xml:"tenants"`
					AppProfiles struct {
						AppProfiles struct {
							Dn   string `xml:"dn"`
							Name string `xml:"name"`
						} `xml:"appProfiles"`
					} `xml:"appProfiles"`
					Interfaces struct {
						Interfaces struct {
							Dn   string `xml:"dn"`
							Name string `xml:"name"`
						} `xml:"interfaces"`
					} `xml:"interfaces"`
					Epgs struct {
						Epgs struct {
							Dn   string `xml:"dn"`
							Name string `xml:"name"`
						} `xml:"epgs"`
					} `xml:"epgs"`
				} `xml:"summary"`
				Category      string `xml:"category"`
				Mnemonic      string `xml:"mnemonic"`
				EventCode     string `xml:"eventCode"`
				Severity      string `xml:"severity"`
				ConditionCode string `xml:"conditionCode"`
				SubCategory   string `xml:"sub_category"`
				OriginTime    int    `xml:"origin_time"`
				DetectionTime int    `xml:"detection_time"`
			} `xml:"data"`
			DataSummary struct {
				Links struct {
					Links struct {
						Rel  string `xml:"rel"`
						Href string `xml:"href"`
					} `xml:"links"`
				} `xml:"links"`
				TotalCount        int  `xml:"total_count"`
				HasMoreData       bool `xml:"has_more_data"`
				PageSize          int  `xml:"page_size"`
				CurrentPageNumber int  `xml:"current_page_number"`
				TotalPageCount    int  `xml:"total_page_count"`
			} `xml:"data_summary"`
		} `xml:"value"`
	} `xml:"ResponseBean"`
}

type EventDetailResponse struct {
	ResponseBean struct {
		Value struct {
			Data struct {
				Category         string `xml:"category"`
				Mnemonic         string `xml:"mnemonic"`
				Severity         string `xml:"severity"`
				Pod              string `xml:"pod"`
				NodeName         string `xml:"nodeName"`
				InterfaceSummary struct {
					Dn          string `xml:"dn"`
					Name        string `xml:"name"`
					DisplayName string `xml:"displayName"`
				} `xml:"interfaceSummary"`
				Details struct {
					Details struct {
						InterfaceDetails struct {
							InterfaceType    string `xml:"interfaceType"`
							InterfaceUseType string `xml:"interfaceUseType"`
							InterfaceProfile struct {
								Dn   string `xml:"dn"`
								Name string `xml:"name"`
							} `xml:"interfaceProfile"`
							InterfaceSelector struct {
								Dn        string `xml:"dn"`
								Name      string `xml:"name"`
								ApicClass string `xml:"apicClass"`
							} `xml:"interfaceSelector"`
							InterfacePolicyGroup struct {
								Dn   string `xml:"dn"`
								Name string `xml:"name"`
							} `xml:"interfacePolicyGroup"`
							Aep struct {
								Dn                        string `xml:"dn"`
								Name                      string `xml:"name"`
								InfrastructureVlanEnabled bool   `xml:"infrastructureVlanEnabled"`
							} `xml:"aep"`
						} `xml:"interfaceDetails"`
					} `xml:"details"`
				} `xml:"details"`
				SubCategory string `xml:"sub_category"`
			} `xml:"data"`
			DataSummary struct {
				TotalCount  int `xml:"total_count"`
				HasMoreData int `xml:"has_more_data"`
			} `xml:"data_summary"`
		} `xml:"value"`
	} `xml:"ResponseBean"`
}

func (c *CandidRestClient) getFabric() ([]byte, error) {
	return c.getRequest(CANDID_PATH_ASSUREDNETWORKS + CANDID_PATH_ACIFABRIC)
}

func parseFabric(fabricXML []byte) ([]string, error) {
	var fabricInfo FabricResponse
	var fabricIDs []string

	if err := xml.Unmarshal(fabricXML, &fabricInfo); err != nil {
		glog.Errorf("Failed ot unmarshal: %+v\n%v", fabricInfo, err)
		return nil, (nil)
	}

	glog.V(2).Infof("response: %v", fabricInfo)

	for _, element := range fabricInfo.ResponseBean.Value.Data {
		fabricIDs = append(fabricIDs, CANDID_PATH_ASSUREDNETWORKS+"/"+element.UUID)
	}
	return fabricIDs, nil
}

func (c *CandidRestClient) getSmartEvents(fabricURL string) ([]byte, error) {
	return c.getRequest(fabricURL + CANDID_PATH_SMARTEVENTS)
}

func parseEvents(eventsXML []byte, fabricURL string) ([]string, error) {
	var eventsInfo SmartEventsResponse
	var eventIDs []string

	if err := xml.Unmarshal(eventsXML, &eventsInfo); err != nil {
		glog.Errorf("Failed ot unmrshal: %+v\n%v", eventsXML, err)
		return nil, (nil)
	}

	glog.V(2).Infof("response: %v", eventsInfo)

	for _, element := range eventsInfo.ResponseBean.Value.Data {
		eventIDs = append(eventIDs, fabricURL+CANDID_PATH_SMARTEVENTS+"/"+element.UUID)
	}
	return eventIDs, nil
}

func (c *CandidRestClient) getEventDetail(eventURL string) ([]byte, error) {
	return c.getRequest(eventURL)
}

func parseEventDetail(eventDetailXML []byte) (string, error) {
	var eventDetail EventDetailResponse

	if err := xml.Unmarshal(eventDetailXML, &eventDetail); err != nil {
		glog.Errorf("Failed ot unmrshal: %+v\n%v", eventDetailXML, err)
		return "", (err)
	}

	glog.V(2).Infof("response: %v", eventDetail)

	var interfaceSummary = eventDetail.ResponseBean.Value.Data.InterfaceSummary
	return interfaceSummary.Dn, nil
}
