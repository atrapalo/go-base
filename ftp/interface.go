package ftp

type Client interface {
	Connect(host string, port int) error
	Login(user string, pass string) error
	Quit() error
	RetrieveFile(remoteFilepath, localFilepath string) error
	StoreFile(remoteFilepath string, data []byte) error
}
