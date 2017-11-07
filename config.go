package carrot

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	Server struct {
		BroadcastChannelSize int64  `yaml:"broadcast_channel_size"`
		Port                 int    `yaml:"port"`
		ServerSecret         string `yaml:"server_secret"`
	}
	Session struct {
		NilSessionToken                     SessionToken  `yaml:"nil_session_token"`
		DefaultSessionClosedTimeoutDuration time.Duration `yaml:"default_sess_closed_timeout_duration_secs"`
	}
	Client struct {
		SendMessageBufferSize int           `yaml:"send_msg_buffer_size"`   // size of client send channel
		SendTokenBufferSize   int           `yaml:"send_token_buffer_size"` // size of sendToken channel
		MaxMessageSize        int64         `yaml:"max_message_size"`       // maximum message size allowed from the websocket
		ClientSecretRequired  bool          `yaml:"client_secret_required"` // toggle to require a client secret token on WS upgrade request
		WriteWaitSecs         time.Duration `yaml:"write_wait_secs"`        // time allowed to write a message to the websocket
		PongWaitSecs          time.Duration `yaml:"pong_wait_secs"`         // time allowed to read the next pong message from the websocket
		writeWait             time.Duration
		pongWait              time.Duration
		pingPeriod            time.Duration // send pings to the websocket with this period, must be less than pongWait
	}
	Router struct {
		RouteDelimiter       string `yaml:"route_delimiter"`
		StreamIdentifier     string `yaml:"stream_identifier"`
		ControllerIdentifier string `yaml:"controller_identifier"`
	}
	Dispatcher struct {
		DoCacheControllers               bool `yaml:"do_cache_controllers"`
		MaxNumCachedControllers          int  `yaml:"max_num_cached_controllers"`
		MaxNumDispatcherIncomingRequests int  `yaml:"max_num_dispatcher_incoming_requests"`
	}
	ClientPool struct {
		MaxClients               int `yaml:"max_clients"`
		MaxClientPoolQueueBackup int `yaml:"max_client_pool_queue_backup"`
		MaxOutboundMessages      int `yaml:"max_outbound_messages"`
	}
	Middleware struct {
		InputChannelSize int     `yaml:"input_channel_size"`
		Rate             float64 `yaml:"rate"`
	}
}

func SetConfig(pathToConfig string) {
	once.Do(func() {
		config = Config{}
		yamlFile, err := ioutil.ReadFile(pathToConfig)
		if err != nil {
			log.Printf("yamlFile.Get err #%v ", err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		setUpClientConfig(&config)
	})
}

func setUpClientConfig(config *Config) {
	clientConfig := &config.Client
	clientConfig.writeWait = clientConfig.WriteWaitSecs * time.Second
	clientConfig.pongWait = clientConfig.PongWaitSecs * time.Second
	clientConfig.pingPeriod = (clientConfig.pongWait * 9) / 10
}
