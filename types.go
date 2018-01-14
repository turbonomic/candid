package main

import (
	"time"
)

const (
	CANDID_PATH_ASSUREDNETWORKS = "/api/v1/assured-networks"
	CANDID_PATH_ACIFABRIC       = "/aci-fabric"
	CANDID_PATH_SMARTEVENTS     = "/smart-events"
	TURBO_PATH_SEARCH           = "/vmturbo/rest/search"
	TURBO_PATH_GROUP            = "/vmturbo/rest/group"
	TURBO_PATH_POLICY           = "/vmturbo/rest/policy"

	defaultTimeOut = time.Duration(60 * time.Second)
)
