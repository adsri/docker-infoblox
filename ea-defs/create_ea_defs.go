package main

import (
	"github.com/infobloxopen/docker-infoblox/common"
	ibclient "github.com/infobloxopen/infoblox-go-client"
	"log"
)

func GetRequiredEADefs() []ibclient.EADefinition {
	ea_defs := common.RequiredEADefs
	res := make([]ibclient.EADefinition, len(ea_defs))
	for i, d := range ea_defs {
		res[i] = *ibclient.NewEADefinition(d)
	}

	return res
}

func main() {
	config, err := common.LoadConfig()
	if config == nil || err != nil {
		log.Fatal(err)
	}

	hostConfig := ibclient.HostConfig{
		Host:     config.GridHost,
		Version:  config.WapiVer,
		Port:     config.WapiPort,
		Username: config.WapiUsername,
		Password: config.WapiPassword,
	}

	transportConfig := ibclient.NewTransportConfig(
		config.SslVerify,
		int(config.HttpRequestTimeout),
		int(config.HttpPoolConnections),
	)

	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}

	conn, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)

	if err != nil {
		log.Fatal(err)
	}

	objMgr := ibclient.NewObjectManager(conn, "Docker", "")

	reqEaDefs := GetRequiredEADefs()
	for _, e := range reqEaDefs {
		eadef, err := objMgr.GetEADefinition(e.Name)

		if err != nil {
			log.Printf("GetEADefinition(%s) error '%s'", e.Name, err)
			continue
		}

		if eadef != nil {
			log.Printf("EA Definition '%s' already exists", eadef.Name)

		} else {
			log.Printf("EA Definition '%s' not found.", e.Name)
			newEadef, err := objMgr.CreateEADefinition(e)
			if err == nil {
				log.Printf("EA Definition '%s' created", newEadef.Name)
			}
		}
	}
}
