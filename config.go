package carrot

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance Config
)

type Config struct {
	Server struct {
		BroadcastChannelSize int `yaml:"broadcast_channel_size"`
		Port                 int     `yaml:"port"`
		ServerSecret         string  `yaml:"server_secret"`
	}
	Session struct {
		NilSessionToken                     SessionToken  `yaml:"nil_session_token"`
		DefaultSessionClosedTimeoutDuration time.Duration `yaml:"default_sess_closed_timeout_duration_secs"`
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
	Client struct {
		SendMessageBufferSize int           `yaml:"send_msg_buffer_size"`
		SendTokenBufferSize   int           `yaml:"send_token_buffer_size"`
		MaxMessageSize        int64         `yaml:"max_message_size"`
		ClientSecretRequired  bool          `yaml:"client_secret_required"`
		WriteWaitSecs         time.Duration `yaml:"write_wait_secs"`
		PongWaitSecs          time.Duration `yaml:"pong_wait_secs"`
	}
	ClientPool struct {
		MaxClients               int `yaml:"max_clients"`
		MaxClientPoolQueueBackup int `yaml:"max_client_pool_queue_backup"`
		MaxOutboundMessages      int `yaml:"max_outbound_messages"`
	}
	Middleware struct {
		InputChannelSize int `yaml:"input_channel_size"`
	}
}

func AddConfig(pathToConfig string) Config {
	once.Do(func() {
		instance = Config{}
		yamlFile, err := ioutil.ReadFile(pathToConfig)
		if err != nil {
			log.Printf("yamlFile.Get err #%v ", err)
		}

		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
	})

	return instance
}
