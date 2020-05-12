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
	sha1ver   string // sha1 revision used to build the program
	buildTime string // Time and date for executable was built
	release string
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
		fmt.Printf("Release %s build on %s from rev %s\n",release, buildTime, sha1ver)
		os.Exit(2)
	}

	if *zabbixProxy == "" && *zabbixHost == "" {
		fmt.Println("Proxy and host must be set")
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
		disres(&data)
	}

	if *poll {
		resourceData(&data)
		nodeData(&data)
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

func disres(mon *CrmMon) bool {

	var resMap []interface{}
	for _, grp := range mon.Resources.Group {
		for _, res := range grp.Resource {
			g := make(map[string]string)
			g["{#GROUP}"] = grp.ID
			g["{#RES}"] = res.ID
			resMap = append(resMap, g)
		}
	}

	SendLLD(cfg.ZabbixTargetHost, cfg.ZabbixProxy, "pacemaker.discover.resources", resMap, "Resource discovery", cfg.Debug)
	return true
}

func resourceData(mon *CrmMon) bool {
	var mData []*Metric

	for _, grp := range mon.Resources.Group {
		for _, res := range grp.Resource {
			//
			if res.Active {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.active["+grp.ID+"."+res.ID+"]", "1", time.Now().Unix()))
			} else {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.active["+grp.ID+"."+res.ID+"]", "0", time.Now().Unix()))
			}

			if res.Failed {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.failed["+grp.ID+"."+res.ID+"]", "1", time.Now().Unix()))
			} else {
				mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.resource.failed["+grp.ID+"."+res.ID+"]", "0", time.Now().Unix()))
			}
		}
	}

	SendMetrics(cfg.ZabbixProxy, mData, "Resources", cfg.Debug)
	return true
}

func nodeData(mon *CrmMon) bool {
	var mData []*Metric
	online := 0
	for _, node := range mon.Nodes.Node {

		if node.Online {
			online++
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.online["+strings.Split(node.Name, ".")[0]+"]", "1", time.Now().Unix()))
		} else {
			mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.online["+strings.Split(node.Name, ".")[0]+"]", "0", time.Now().Unix()))
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
	mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.online", strconv.Itoa(online), time.Now().Unix()))
	mData = append(mData, NewMetric(cfg.ZabbixTargetHost, "pacemaker.node.total", strconv.Itoa(len(mon.Nodes.Node)), time.Now().Unix()))

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
