package define

type RedisConfig struct {
	MaxIdleConn    int    `yaml:"maxIdleConn" json:"max_idle_conn,omitempty"`
	MaxActiveConn  int    `yaml:"maxActiveConn" json:"max_active_conn,omitempty"`
	MaxIdleTimeout int    `yaml:"max_idle_timeout" json:"max_idle_timeout,omitempty"`
	DialAddress    string `yaml:"dialAddress" json:"dial_address,omitempty"`
	DBIndex        int    `yaml:"dbIndex" json:"db_index,omitempty"`
	Password       string `yaml:"password" json:"password,omitempty"`
	ClusterMode    bool   `yaml:"cluster_mode" json:"cluster_mode,omitempty"`
	Timeout        int    `yaml:"timeout" json:"timeout,omitempty"`
	KeyPrefix      string `yaml:"key_prefix" json:"key_prefix,omitempty"`
}
