package models

type EnvVariables struct {
	Port         string `default:"8080"`
	DatabaseURL  string `split_words:"true"`
	DatabaseName string `split_words:"true" required:"true"`
	Collection   string `required:"true"`
	RoboCpuUrl   string `default:"https://robotstakeover20210903110417.azurewebsites.net/robotcpu"`
	LogLevel     string `split_words:"true"`
	LogFormat    string `split_words:"true" default:"console"`
}
type Survivor struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Gender   string   `json:"gender"`
	Location Location `json:"location"`
	Resource []string `json:"resource"`
	Status   int      `json:"status"`
}
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}
