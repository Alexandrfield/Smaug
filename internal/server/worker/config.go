package worker

type Config struct {
	MasterKey  []byte
	DatabasDSN string
	Port       string
}
