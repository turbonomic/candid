package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/glog"
)

type TurboRestClient struct {
	client   *http.Client
	host     string
	username string
	password string
}

func NewTurboRestClient(host, uname, pass string) (*TurboRestClient, error) {

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

	return &TurboRestClient{
		client:   client,
		host:     host,
		username: uname,
		password: pass,
	}, nil
}

type Group struct {
	Links []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	// Group attributes
	UUID            string `json:"uuid"`
	DisplayName     string `json:"displayName"`
	ClassName       string `json:"className"`
	EntitiesCount   int    `json:"entitiesCount"`
	MembersCount    int    `json:"membersCount"`
	GroupType       string `json:"groupType"`
	Severity        string `json:"severity"`
	IsStatic        bool   `json:"isStatic"`
	LogicalOperator string `json:"logicalOperator"`
	EnvironmentType string `json:"environmentType"`
}

type Policy struct {
	Links []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	UUID          string `json:"uuid"`
	DisplayName   string `json:"displayName"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Enabled       bool   `json:"enabled"`
	Capacity      int    `json:"capacity"`
	CommodityType string `json:"commodityType"`
	ConsumerGroup Group
	ProviderGroup Group
}

func (c *TurboRestClient) getGroup(groupName string) ([]byte, error) {
	//1. a new http request
	urlStr := fmt.Sprintf("%s%s", c.host, TURBO_PATH_SEARCH+"?types=Group&q="+groupName)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return nil, err
	}

	//2. set queries
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	//3. set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

func parseGroupUUID(groupResponse []byte) (string, error) {
	var group Group

	if err := xml.Unmarshal(groupResponse, &group); err != nil {
		glog.Errorf("Failed ot unmrshal: %+v\n%v", groupResponse, err)
		return "", (err)
	}

	glog.V(2).Infof("response: %v", group)

	return group.UUID, nil
}

func (c *TurboRestClient) updateGroup(groupUUID string, memberUUIDs []string) ([]byte, error) {
	//0. create request data
	var members string
	for index, element := range memberUUIDs {
		//element: topology/pod-1/node-101/sys/phys-[eth1/5]
		dn := strings.Split(element,"/")
		leaf := strings.Replace(dn[2], "node", "leaf", -1)
		if index == 0 {
			members = leaf
		} else {
			members = members + "|" + leaf
		}
	}

	var groupDATA = "{\"isStatic\":false,\"memberUuidList\":[],\"displayName\":\"PMs_CandidLinkDown\",\"groupType\":\"PhysicalMachine\",\"criteriaList\":[{\"expType\":\"EQ\",\"expVal\":\"" + members + "\",\"filterType\":\"pmsBySwitch\",\"caseSensitive\":false}]}"

	//1. a new http request
	urlStr := fmt.Sprintf("%s%s", c.host, TURBO_PATH_GROUP+"/"+groupUUID)
	req, err := http.NewRequest("PUT", urlStr, bytes.NewBufferString(groupDATA))
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return nil, err
	}

	//2. set queries
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	//3. set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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
