package main

import "encoding/xml"

type CrmMon struct {
	XMLName xml.Name `xml:"crm_mon"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Summary struct {
		Text  string `xml:",chardata"`
		Stack struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"stack"`
		CurrentDc struct {
			Text       string `xml:",chardata"`
			Present    string `xml:"present,attr"`
			Version    string `xml:"version,attr"`
			Name       string `xml:"name,attr"`
			ID         string `xml:"id,attr"`
			WithQuorum string `xml:"with_quorum,attr"`
		} `xml:"current_dc"`
		LastUpdate struct {
			Text string `xml:",chardata"`
			Time string `xml:"time,attr"`
		} `xml:"last_update"`
		LastChange struct {
			Text   string `xml:",chardata"`
			Time   string `xml:"time,attr"`
			User   string `xml:"user,attr"`
			Client string `xml:"client,attr"`
			Origin string `xml:"origin,attr"`
		} `xml:"last_change"`
		NodesConfigured struct {
			Text          string `xml:",chardata"`
			Number        string `xml:"number,attr"`
			ExpectedVotes string `xml:"expected_votes,attr"`
		} `xml:"nodes_configured"`
		ResourcesConfigured struct {
			Text     string `xml:",chardata"`
			Number   string `xml:"number,attr"`
			Disabled string `xml:"disabled,attr"`
			Blocked  string `xml:"blocked,attr"`
		} `xml:"resources_configured"`
		ClusterOptions struct {
			Text             string `xml:",chardata"`
			StonithEnabled   bool   `xml:"stonith-enabled,attr"`
			SymmetricCluster bool   `xml:"symmetric-cluster,attr"`
			NoQuorumPolicy   string `xml:"no-quorum-policy,attr"`
			MaintenanceMode  bool   `xml:"maintenance-mode,attr"`
		} `xml:"cluster_options"`
	} `xml:"summary"`
	Nodes struct {
		Text string `xml:",chardata"`
		Node []struct {
			Text             string `xml:",chardata"`
			Name             string `xml:"name,attr"`
			ID               string `xml:"id,attr"`
			Online           bool   `xml:"online,attr"`
			Standby          bool   `xml:"standby,attr"`
			StandbyOnfail    bool   `xml:"standby_onfail,attr"`
			Maintenance      bool   `xml:"maintenance,attr"`
			Pending          string `xml:"pending,attr"`
			Unclean          string `xml:"unclean,attr"`
			Shutdown         bool   `xml:"shutdown,attr"`
			ExpectedUp       string `xml:"expected_up,attr"`
			IsDc             string `xml:"is_dc,attr"`
			ResourcesRunning string `xml:"resources_running,attr"`
			Type             string `xml:"type,attr"`
		} `xml:"node"`
	} `xml:"nodes"`
	Resources struct {
		Text  string `xml:",chardata"`
		Resource []struct {
			Text           string `xml:",chardata"`
			ID             string `xml:"id,attr"`
			ResourceAgent  string `xml:"resource_agent,attr"`
			Role           string `xml:"role,attr"`
			Active         bool   `xml:"active,attr"`
			Orphaned       bool   `xml:"orphaned,attr"`
			Blocked        bool   `xml:"blocked,attr"`
			Managed        bool   `xml:"managed,attr"`
			Failed         bool   `xml:"failed,attr"`
			FailureIgnored bool   `xml:"failure_ignored,attr"`
			NodesRunningOn string `xml:"nodes_running_on,attr"`
			Node           struct {
				Text   string `xml:",chardata"`
				Name   string `xml:"name,attr"`
				ID     string `xml:"id,attr"`
				Cached string `xml:"cached,attr"`
			} `xml:"node"`
		} `xml:"resource"`
		Group []struct {
			Text            string `xml:",chardata"`
			ID              string `xml:"id,attr"`
			NumberResources string `xml:"number_resources,attr"`
			Resource        []struct {
				Text           string `xml:",chardata"`
				ID             string `xml:"id,attr"`
				ResourceAgent  string `xml:"resource_agent,attr"`
				Role           string `xml:"role,attr"`
				Active         bool   `xml:"active,attr"`
				Orphaned       bool   `xml:"orphaned,attr"`
				Blocked        bool   `xml:"blocked,attr"`
				Managed        bool   `xml:"managed,attr"`
				Failed         bool   `xml:"failed,attr"`
				FailureIgnored bool   `xml:"failure_ignored,attr"`
				NodesRunningOn string `xml:"nodes_running_on,attr"`
				Node           struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					ID     string `xml:"id,attr"`
					Cached string `xml:"cached,attr"`
				} `xml:"node"`
			} `xml:"resource"`
		} `xml:"group"`
	} `xml:"resources"`
	NodeAttributes struct {
		Text string `xml:",chardata"`
		Node []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"node"`
	} `xml:"node_attributes"`
	NodeHistory struct {
		Text string `xml:",chardata"`
		Node []struct {
			Text            string `xml:",chardata"`
			Name            string `xml:"name,attr"`
			ResourceHistory []struct {
				Text               string `xml:",chardata"`
				ID                 string `xml:"id,attr"`
				Orphan             string `xml:"orphan,attr"`
				MigrationThreshold string `xml:"migration-threshold,attr"`
				FailCount          string `xml:"fail-count,attr"`
				LastFailure        string `xml:"last-failure,attr"`
				OperationHistory   []struct {
					Text         string `xml:",chardata"`
					Call         string `xml:"call,attr"`
					Task         string `xml:"task,attr"`
					LastRcChange string `xml:"last-rc-change,attr"`
					LastRun      string `xml:"last-run,attr"`
					ExecTime     string `xml:"exec-time,attr"`
					QueueTime    string `xml:"queue-time,attr"`
					Rc           string `xml:"rc,attr"`
					RcText       string `xml:"rc_text,attr"`
					Interval     string `xml:"interval,attr"`
				} `xml:"operation_history"`
			} `xml:"resource_history"`
		} `xml:"node"`
	} `xml:"node_history"`
	Failures struct {
		Text    string `xml:",chardata"`
		Failure []struct {
			Text         string `xml:",chardata"`
			OpKey        string `xml:"op_key,attr"`
			Node         string `xml:"node,attr"`
			Exitstatus   string `xml:"exitstatus,attr"`
			Exitreason   string `xml:"exitreason,attr"`
			Exitcode     string `xml:"exitcode,attr"`
			Call         string `xml:"call,attr"`
			Status       string `xml:"status,attr"`
			LastRcChange string `xml:"last-rc-change,attr"`
			Queued       string `xml:"queued,attr"`
			Exec         string `xml:"exec,attr"`
			Interval     string `xml:"interval,attr"`
			Task         string `xml:"task,attr"`
		} `xml:"failure"`
	} `xml:"failures"`
	Tickets string `xml:"tickets"`
	Bans    struct {
		Text string `xml:",chardata"`
		Ban  struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id,attr"`
			Resource   string `xml:"resource,attr"`
			Node       string `xml:"node,attr"`
			Weight     string `xml:"weight,attr"`
			MasterOnly string `xml:"master_only,attr"`
		} `xml:"ban"`
	} `xml:"bans"`
}
