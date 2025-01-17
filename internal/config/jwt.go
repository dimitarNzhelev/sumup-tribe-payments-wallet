package config

type JWTConfig struct {
	// Secret is the secret key used to sign the JWT token.
	Secret string `default:"your-very-secret-key" envconfig:"JWT_SECRET"`
}
