package config

type AppEnv string

const (
	AppEnvDevelopment AppEnv = "development"
	AppEnvProduction  AppEnv = "production"
)

func (e AppEnv) IsProduction() bool {
	return e == AppEnvProduction
}

func (e AppEnv) IsDevelopment() bool {
	return e == AppEnvDevelopment
}
