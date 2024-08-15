package flag

import "flag"

var Port *string

func init() {
	Port = flag.String("port", "10002", "server port")
	flag.Parse()
}
