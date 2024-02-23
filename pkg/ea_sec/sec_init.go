package sec

import initmgr "junodb_lite/pkg/e_initmgr"

var (
	Initializer initmgr.IInitializer = initmgr.NewInitializer(Initialize, Finalize)
)

func Initialize(args ...interface{}) (err error) {
	return nil
}

func Finalize() {
}
