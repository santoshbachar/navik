package constants

//const ResourceDir = "/Users/santoshbachar/dev/go/src/github.com/santoshbachar/navik/resources/"
// const ResourceDir = "/home/santosh/dev/go/navik/resources/"

// my demo VM
const ResourceDir = "./resources/"

var common_args string

func GetCommonArgs() string {
	return common_args
}

func SetCommonArgs(args string) {
	common_args = args
}
