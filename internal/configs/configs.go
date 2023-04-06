package configs

//const Port = ":8080"
//const Host = "http://localhost" + Port + "/"

type Config struct {
	//Port          string `env:"PORT" envDefault:":8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`

	//Host string `env:"HOST" envDefault:"http://localhost:8080/"`
}
