package global

type Notice struct {
	Enable   bool `yaml:"enable"`
	DingTalk struct {
		Enable  bool   `yaml:"enable"`
		WebHook string `yaml:"web_hook"`
	} `yaml:"DingTalk"`
	QyWeixin struct {
		Enable     bool   `yaml:"enable"`
		CorpID     string `yaml:"corp_id"`
		CorpSecret string `yaml:"corp_secret"`
	} `yaml:"QyWeixin"`
}

type KVConfig struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
}

// Config 服务器设置
type Config struct {
	MQTT          KVConfig `mqtt`
	Addr          string   `yaml:"addr"`          // 监听地址
	Data          string   `yaml:"data"`          // 数据目录
	Negative_data string   `yaml:"negative_data"` // 负面词数据目录
	Debug         bool     `yaml:"debug"`         // 调试模式
	INFO          bool     `yaml:"info"`          // 输出信息
	AllowDrop     bool     `yaml:"allow_drop"`    // 允许删除数据库
	Dictionary    string   `yaml:"dictionary"`    // 字典路径
	EnableAdmin   bool     `yaml:"enableAdmin"`   //启用admin
	Gomaxprocs    int      `yaml:"gomaxprocs"`    //GOMAXPROCS
	Shard         int      `yaml:"shard"`         //分片数
	Auth          string   `yaml:"auth"`          //认证
	EnableGzip    bool     `yaml:"enableGzip"`    //是否开启gzip压缩
	Timeout       int64    `yaml:"timeout"`       //超时时间
	BufferNum     int      `yaml:"bufferNum"`     //分片缓冲数
	Notice        Notice   `yaml:"notice"`        //通知配置
}
