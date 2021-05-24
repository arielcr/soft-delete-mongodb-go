package entities

// MongoDBConfig Config struct to store the database connection string
type MongoDBConfig struct {
	URI        string `env:"MONGODB_URI"`
	Database   string `env:"MONGODB_DATABASE"`
	Collection string `env:"MONGODB_COLLECTION"`
}
