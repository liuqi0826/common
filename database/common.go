package database

func EditMySQLPath(address, user, passwrod, database, charset string) (path string) {
	path += user
	path += ":"
	path += passwrod
	path += "@tcp("
	path += address
	path += ")/"
	path += database
	path += "?"
	if charset == "" {
		path += "charset=utf8mb4"
	} else {
		path += "charset=" + charset
	}
	return
}

func CreateMySQLConnect(address string) (proxy *SQLProxy, err error) {
	proxy = &SQLProxy{}
	err = proxy.Connect(DATABASE_MYSQL, address)
	return
}
func CreateSQLite3Connect(path string) (proxy *SQLProxy, err error) {
	proxy = &SQLProxy{}
	err = proxy.Connect(DATABASE_SQLITE3, path)
	return
}
func CreateRedisConnect(address, passwrod string, database int) (proxy *CacheProxy, err error) {
	proxy = &CacheProxy{}
	err = proxy.ConnectRedis(address, passwrod, database)
	return
}
func CreateSSDBConnect(address string, port int) (proxy *CacheProxy, err error) {
	proxy = &CacheProxy{}
	err = proxy.ConnectSSDB(address, port)
	return
}
func CreateLevelDBConnect(path string) (proxy *CacheProxy, err error) {
	proxy = &CacheProxy{}
	err = proxy.ConnectLevelDB(path)
	return
}
