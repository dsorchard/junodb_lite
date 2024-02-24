package patch

import (
	config "junodb_lite/cmd/bb_dbscanserv/b_config"
	app "junodb_lite/cmd/bb_dbscanserv/c_app"
)

// Called by storageserv.
func Init(cfg *config.DbScan) {
	app.InitPatch(cfg)
}
