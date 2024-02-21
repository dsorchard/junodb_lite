package etcd

import "fmt"

const (
	TagCompDelimiter          = "_"
	TagVersion                = "version"
	TagAlgVer                 = "algver"
	TagNumZones               = "numzones"
	TagNumShards              = "numshards"
	TagNodePrefix             = "node"
	TagNodeIpport             = "node_ipport"
	TagNodeShards             = "node_shards"
	TagShardDelimiter         = ","
	TagPrimSecondaryDelimiter = "|"
	TagZoneMarkDown           = "zonemarkdown"
	TagLimitsConfig           = "config_limits"
	TagRedistNodeIpport       = "redist_node_ipport"
	TagRedistNodeShards       = "redist_node_shards"
	TagRedistPrefix           = "redist"
	TagRedistEnablePrefix     = "redist_enable"
	TagRedistEnabledTarget    = "yes_target"
	TagRedistEnabledSource    = "yes_source"

	// vaules
	TagRedistEnabledReady    = "ready"
	TagRedistEnabledSourceRL = "yes_source_rl"
	TagRedistDisabled        = "no"
	TagRedistAbortAll        = "abort_all"
	TagRedistAbortZone       = "abort_zone"
	TagRedistResume          = "source_resume"
	TagRedistResumeRL        = "source_resume_rl"
	TagRedistTgtStateInit    = "init"
	TagRedistTgtStateReady   = "ready"
	TagRedistRateLimit       = "ratelimit"
	TagFieldSeparator        = "|"
	TagKeyValueSeparator     = "="
)

func KeyNodeIpport(zone int, node int) string {
	return Key(TagNodeIpport, zone, node)
}

func KeyNodeShards(zone int, node int) string {
	return Key(TagNodeShards, zone, node)
}

func Key(Prefix string, list ...int) string {
	var key string = Prefix
	if len(list) == 0 {
		return key
	}

	if len(list) >= 1 {
		key = fmt.Sprintf("%s%s%02d", key, TagCompDelimiter, list[0])
	}
	if len(list) >= 2 {
		key = fmt.Sprintf("%s%s%03d", key, TagCompDelimiter, list[1])
	}
	for i := 2; i < len(list); i++ {
		key = fmt.Sprintf("%s%s%05d", key, TagCompDelimiter, list[i])
	}

	return key
}

func KeyRedistEnable(zone int) string {
	return Key(TagRedistEnablePrefix, zone)
}
