package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	. "github.com/marstid/go-zabbix"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"time"
)

var (
	sha1ver   string
	buildTime string
	release   string
)

type Config struct {
	ZabbixProxy      string
	ZabbixTargetHost string
	Debug            bool
}

var cfg Config

func main() {

	discover := flag.Bool("discover", false, "Run discovery and send LLD to Zabbix")
	poll := flag.Bool("poll", false, "Poll data and send to Zabbix")
	version := flag.Bool("version", false, "Print version and build date")
	zabbixProxy := flag.String("P", "", "Zabbix Proxy server to recieve data")
	zabbixHost := flag.String("H", "", "Target Zabbix host")
	flag.Parse()

	if *version {
		fmt.Printf("Release %s build on %s from rev %s\n", release, buildTime, sha1ver)
		os.Exit(2)
	}

	if *zabbixProxy == "" && *zabbixHost == "" {
		fmt.Println("Proxy and host must be set")
		os.Exit(2)
	}

	if *zabbixProxy == "" {
		fmt.Println("Proxy must be set")
		os.Exit(2)
	}

	if *zabbixHost == "" {
		fmt.Println("Host must be set")
		os.Exit(2)
	}

	cfg = Config{
		ZabbixProxy:      *zabbixProxy,
		ZabbixTargetHost: *zabbixHost,
		Debug:            false,
	}

	var data CrmMon
	data = getData()

	if *discover {
		disNode(&data)
		disRes(&data)
		disResGroup(&data)
	}

	if *poll {
		resourceData(&data)
		resourceDataGroup(&data)
		nodeData(&data)
		clusterData(&data)
	}
}

func disNode(mon *CrmMon) bool {

	var nodeMap []interface{}
	for _, node := range mon.Nodes.Node {
		ne := make(map[string]string)
		ne["{#NAME}"] = strings.Split(node.Name, ".")[0]
		nodeMap = append(nodeMap, ne)
	}

	SendLLD(cfg.ZabbixTargetHost, cfg.ZabbixProxy, "pacemaker.discover.node", nodeMap, "Node discovery", cfg.Debug)
	return true
}

func disRes(mon *CrmMon) bool {

	var resMap []interface{}
	for _, res := range mon.Resources.Resource {
		g := make(map[string]string)
		g["{#RES}"] = res.ID
		resMap = append(resMap, g)
	}

	SendLLD(cfg.ZabbixTargetHost, cfg.ZabbixProxy, "pacemaker.discover.resources", resMap, "Resource discovery", cfg.Debug)
	return true
}

func disResGroup(mon *CrmMon) bool {

	var resMap []interface{}
	for _, grp := range mon.Resources.Group {
		for _, res := range grp.Resource {
			g := make(map[string]string)
			g["{#GROUP}"] = grp.ID
			g["{#RES}"] = res.ID
			resMap = append(resMap, g)
		}
	}

	SendLLD(cfg.ZabbixTargetHost, cfg.ZabbixProxy, "pacemaker.discover.resourcesgroup", resMap, "Resource group discovery", cfg.Debug)
	return true
}

func resourceData(mon *CrmMon) bool {
	var mData []*Metric

	for _, res := range mon.Resources.Resource {
		if res.Active {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.active["+res.ID+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.active["+res.ID+"]", "0", time.Now().Unix()))
		}

		if res.Failed {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.failed["+res.ID+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.failed["+res.ID+"]", "0", time.Now().Unix()))
		}
	}

	SendMetrics(cfg.ZabbixProxy, mData, "Resources", cfg.Debug)
	return true
}

func resourceDataGroup(mon *CrmMon) bool {
	var mData []*Metric

	for _, grp := range mon.Resources.Group {
		for _, res := range grp.Resource {
			if res.Active {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resourcegroup.active["+grp.ID+"."+res.ID+"]", "1", time.Now().Unix()))
			} else {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resourcegroup.active["+grp.ID+"."+res.ID+"]", "0", time.Now().Unix()))
			}

			if res.Failed {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resourcegroup.failed["+grp.ID+"."+res.ID+"]", "1", time.Now().Unix()))
			} else {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resourcegroup.failed["+grp.ID+"."+res.ID+"]", "0", time.Now().Unix()))
			}
		}
	}

	SendMetrics(cfg.ZabbixProxy, mData, "Resources", cfg.Debug)
	return true
}

func nodeData(mon *CrmMon) bool {
	var mData []*Metric
	online := 0
	standby := 0
	standby_onfail := 0
	for _, node := range mon.Nodes.Node {

		if node.Online {
			online++
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.online["+strings.Split(node.Name, ".")[0]+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.online["+strings.Split(node.Name, ".")[0]+"]", "0", time.Now().Unix()))
		}

		if node.Standby {
			standby++
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.standby["+strings.Split(node.Name, ".")[0]+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.standby["+strings.Split(node.Name, ".")[0]+"]", "0", time.Now().Unix()))
		}

		if node.StandbyOnfail {
			standby_onfail++
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.standby_onfail["+strings.Split(node.Name, ".")[0]+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.standby_onfail["+strings.Split(node.Name, ".")[0]+"]", "0", time.Now().Unix()))
		}

		if node.Shutdown {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.shutdown["+strings.Split(node.Name, ".")[0]+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.shutdown["+strings.Split(node.Name, ".")[0]+"]", "0", time.Now().Unix()))
		}

		if node.Maintenance {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.maint["+strings.Split(node.Name, ".")[0]+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.maint["+strings.Split(node.Name, ".")[0]+"]", "0", time.Now().Unix()))
		}
	}
	mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.total", strconv.Itoa(len(mon.Nodes.Node)), time.Now().Unix()))
	mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.online", strconv.Itoa(online), time.Now().Unix()))
	mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.standby", strconv.Itoa(standby), time.Now().Unix()))
	mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.standby_onfail", strconv.Itoa(standby_onfail), time.Now().Unix()))

	SendMetrics(cfg.ZabbixProxy, mData, "Node", cfg.Debug)
	return true
}

func clusterData(mon *CrmMon) bool {
	var mData []*Metric

	if mon.Summary.ClusterOptions.MaintenanceMode {
		mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.cluster.maint", "1", time.Now().Unix()))
	} else {
		mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.cluster.maint", "0", time.Now().Unix()))
	}

	SendMetrics(cfg.ZabbixProxy, mData, "Node", cfg.Debug)
	return true
}

func getData() (data CrmMon) {

	cmd := exec.Command("crm_mon", "-X")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("crm_mon -X failed with %s\n", err)
	}

	xml.Unmarshal(out, &data)
	return

}
