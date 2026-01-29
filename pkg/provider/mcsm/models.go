package mcsm

type InstanceDetail struct {
	Status int `json:"status"`
	Data   struct {
		InstanceUUID  string `json:"instanceUuid"`
		Started       int    `json:"started"`
		AutoRestarted int    `json:"autoRestarted"`
		Status        int    `json:"status"`
		Config        struct {
			Nickname          string        `json:"nickname"`
			StartCommand      string        `json:"startCommand"`
			StopCommand       string        `json:"stopCommand"`
			Cwd               string        `json:"cwd"`
			Ie                string        `json:"ie"`
			Oe                string        `json:"oe"`
			CreateDatetime    int64         `json:"createDatetime"`
			LastDatetime      int64         `json:"lastDatetime"`
			Type              string        `json:"type"`
			Tag               []string      `json:"tag"`
			EndTime           int           `json:"endTime"`
			FileCode          string        `json:"fileCode"`
			ProcessType       string        `json:"processType"`
			UpdateCommand     string        `json:"updateCommand"`
			RunAs             string        `json:"runAs"`
			Crlf              int           `json:"crlf"`
			Category          int           `json:"category"`
			BasePort          int           `json:"basePort"`
			EnableRcon        bool          `json:"enableRcon"`
			RconPassword      string        `json:"rconPassword"`
			RconPort          int           `json:"rconPort"`
			RconIP            string        `json:"rconIp"`
			ActionCommandList []interface{} `json:"actionCommandList"`
			TerminalOption    struct {
				HaveColor    bool `json:"haveColor"`
				Pty          bool `json:"pty"`
				PtyWindowCol int  `json:"ptyWindowCol"`
				PtyWindowRow int  `json:"ptyWindowRow"`
			} `json:"terminalOption"`
			EventTask struct {
				AutoStart           bool `json:"autoStart"`
				AutoRestart         bool `json:"autoRestart"`
				AutoRestartMaxTimes int  `json:"autoRestartMaxTimes"`
				Ignore              bool `json:"ignore"`
			} `json:"eventTask"`
			Java struct {
				ID string `json:"id"`
			} `json:"java"`
			Docker struct {
				ContainerName    string        `json:"containerName"`
				Image            string        `json:"image"`
				Ports            []string      `json:"ports"`
				ExtraVolumes     []interface{} `json:"extraVolumes"`
				Memory           int           `json:"memory"`
				MemorySwap       interface{}   `json:"memorySwap"`
				MemorySwappiness interface{}   `json:"memorySwappiness"`
				NetworkMode      string        `json:"networkMode"`
				NetworkAliases   []string      `json:"networkAliases"`
				CpusetCpus       string        `json:"cpusetCpus"`
				CPUUsage         int           `json:"cpuUsage"`
				MaxSpace         int           `json:"maxSpace"`
				Io               int           `json:"io"`
				Network          int           `json:"network"`
				WorkingDir       string        `json:"workingDir"`
				Env              []string      `json:"env"`
				ChangeWorkdir    bool          `json:"changeWorkdir"`
				Labels           []string      `json:"labels"`
			} `json:"docker"`
			PingConfig struct {
				IP   string `json:"ip"`
				Port int    `json:"port"`
				Type int    `json:"type"`
			} `json:"pingConfig"`
			ExtraServiceConfig struct {
				OpenFrpTunnelID string `json:"openFrpTunnelId"`
				OpenFrpToken    string `json:"openFrpToken"`
				IsOpenFrp       bool   `json:"isOpenFrp"`
			} `json:"extraServiceConfig"`
		} `json:"config"`
		Info struct {
			McPingOnline   bool          `json:"mcPingOnline"`
			CurrentPlayers int           `json:"currentPlayers"`
			MaxPlayers     int           `json:"maxPlayers"`
			Version        string        `json:"version"`
			FileLock       int           `json:"fileLock"`
			PlayersChart   []interface{} `json:"playersChart"`
			OpenFrpStatus  bool          `json:"openFrpStatus"`
			Latency        int           `json:"latency"`
			AllocatedPorts []struct {
				Host      int    `json:"host"`
				Container int    `json:"container"`
				Protocol  string `json:"protocol"`
			} `json:"allocatedPorts"`
			CPUUsage           int   `json:"cpuUsage"`
			RxBytes            int   `json:"rxBytes"`
			TxBytes            int   `json:"txBytes"`
			MemoryUsagePercent int   `json:"memoryUsagePercent"`
			StorageUsage       int64 `json:"storageUsage"`
			StorageLimit       int64 `json:"storageLimit"`
			MemoryUsage        int64 `json:"memoryUsage"`
			MemoryLimit        int64 `json:"memoryLimit"`
		} `json:"info"`
		Space       int         `json:"space"`
		ProcessInfo interface{} `json:"processInfo"`
	} `json:"data"`
	Time int64 `json:"time"`
}

type InstanceSearch struct {
	Status int `json:"status"`
	Data   struct {
		Page     int      `json:"page"`
		PageSize int      `json:"pageSize"`
		MaxPage  int      `json:"maxPage"`
		AllTags  []string `json:"allTags"`
		Data     []struct {
			InstanceUUID  string `json:"instanceUuid"`
			Started       int    `json:"started"`
			AutoRestarted int    `json:"autoRestarted"`
			Status        int    `json:"status"`
			Config        struct {
				Nickname          string        `json:"nickname"`
				StartCommand      string        `json:"startCommand"`
				StopCommand       string        `json:"stopCommand"`
				Cwd               string        `json:"cwd"`
				Ie                string        `json:"ie"`
				Oe                string        `json:"oe"`
				CreateDatetime    int64         `json:"createDatetime"`
				LastDatetime      int64         `json:"lastDatetime"`
				Type              string        `json:"type"`
				Tag               []string      `json:"tag"`
				EndTime           int           `json:"endTime"`
				FileCode          string        `json:"fileCode"`
				ProcessType       string        `json:"processType"`
				UpdateCommand     string        `json:"updateCommand"`
				RunAs             string        `json:"runAs"`
				Crlf              int           `json:"crlf"`
				Category          int           `json:"category"`
				BasePort          int           `json:"basePort"`
				EnableRcon        bool          `json:"enableRcon"`
				RconPassword      string        `json:"rconPassword"`
				RconPort          int           `json:"rconPort"`
				RconIP            string        `json:"rconIp"`
				ActionCommandList []interface{} `json:"actionCommandList"`
				TerminalOption    struct {
					HaveColor    bool `json:"haveColor"`
					Pty          bool `json:"pty"`
					PtyWindowCol int  `json:"ptyWindowCol"`
					PtyWindowRow int  `json:"ptyWindowRow"`
				} `json:"terminalOption"`
				EventTask struct {
					AutoStart           bool `json:"autoStart"`
					AutoRestart         bool `json:"autoRestart"`
					AutoRestartMaxTimes int  `json:"autoRestartMaxTimes"`
					Ignore              bool `json:"ignore"`
				} `json:"eventTask"`
				Java struct {
					ID string `json:"id"`
				} `json:"java"`
				Docker struct {
					ContainerName    string        `json:"containerName"`
					Image            string        `json:"image"`
					Ports            []string      `json:"ports"`
					ExtraVolumes     []string      `json:"extraVolumes"`
					Memory           int           `json:"memory"`
					MemorySwap       interface{}   `json:"memorySwap"`
					MemorySwappiness interface{}   `json:"memorySwappiness"`
					NetworkMode      string        `json:"networkMode"`
					NetworkAliases   []interface{} `json:"networkAliases"`
					CpusetCpus       string        `json:"cpusetCpus"`
					CPUUsage         int           `json:"cpuUsage"`
					MaxSpace         int           `json:"maxSpace"`
					Io               int           `json:"io"`
					Network          int           `json:"network"`
					WorkingDir       string        `json:"workingDir"`
					Env              []string      `json:"env"`
					ChangeWorkdir    bool          `json:"changeWorkdir"`
					Labels           []interface{} `json:"labels"`
				} `json:"docker"`
				PingConfig struct {
					IP   string `json:"ip"`
					Port int    `json:"port"`
					Type int    `json:"type"`
				} `json:"pingConfig"`
				ExtraServiceConfig struct {
					OpenFrpTunnelID string `json:"openFrpTunnelId"`
					OpenFrpToken    string `json:"openFrpToken"`
				} `json:"extraServiceConfig"`
			} `json:"config"`
			Info struct {
				McPingOnline   bool          `json:"mcPingOnline"`
				CurrentPlayers int           `json:"currentPlayers"`
				MaxPlayers     int           `json:"maxPlayers"`
				Version        string        `json:"version"`
				FileLock       int           `json:"fileLock"`
				PlayersChart   []interface{} `json:"playersChart"`
				OpenFrpStatus  bool          `json:"openFrpStatus"`
				Latency        int           `json:"latency"`
				AllocatedPorts []struct {
					Host      int    `json:"host"`
					Container int    `json:"container"`
					Protocol  string `json:"protocol"`
				} `json:"allocatedPorts"`
				CPUUsage           int     `json:"cpuUsage"`
				RxBytes            float64 `json:"rxBytes"`
				TxBytes            float64 `json:"txBytes"`
				MemoryUsagePercent int     `json:"memoryUsagePercent"`
				MemoryUsage        int     `json:"memoryUsage"`
				MemoryLimit        int64   `json:"memoryLimit"`
				StorageUsage       int     `json:"storageUsage"`
				StorageLimit       int64   `json:"storageLimit"`
			} `json:"info"`
		} `json:"data"`
	} `json:"data"`
	Time int64 `json:"time"`
}
