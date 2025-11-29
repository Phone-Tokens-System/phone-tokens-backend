package env

type CredsEnv struct {
	Email  string `env:"EMAIL"`
	ApiKey string `env:"API_KEY"`
}
