package resource

type Config struct {
	//dsn := "root@tcp(localhost:3306)/coursera?"
	// указываем кодировку
	//dsn += "&charset=utf8"
	// отказываемся от prapared statements
	// параметры подставляются сразу
	//dsn += "&interpolateParams=true"
	//RESTAPIPort int    `envconfig:"PORT" default:"8080" required:"true"`
	DBURL  string `envconfig:"DB_URL" default:"postgres://db-user:db-password@localhost:5429/tripdb?sslmode=disable" required:"true"`
	ApiKey string `envconfig:"Api_Key" default:"WOp3bI0eN2_gEG1ob-orRSViXwd-53mYAa_Vn8dyuMM" required:"true"`
}

/*type HereConfig struct {
	AppID string // ibKEG20RngEygrvyo1I8
	ApiKey string // WOp3bI0eN2_gEG1ob-orRSViXwd-53mYAa_Vn8dyuMM
}*/
