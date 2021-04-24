package lutra

type LutraConfig struct {
	SmtpServer           string
	SmtpLogin            string
	SmtpPassword         string
	SmtpPort             int
	MongoDbConnectionStr string
	MongoDbName          string
	PrivateKeyPw         string
}
