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
)

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

func KeyNodeIpport(zone int, node int) string {
	return Key(TagNodeIpport, zone, node)
}

func KeyNodeShards(zone int, node int) string {
	return Key(TagNodeShards, zone, node)
}

// Keys for redistribution
var (
	TagRedistEnablePrefix       = "redist_enable"
	TagRedistFromNode           = "redist_from"
	TagRedistPrefix             = "redist"
	TagRedistNodePrefix         = "redist_node"
	TagRedistNodeIpport         = "redist_node_ipport"
	TagRedistNodeShards         = "redist_node_shards"
	TagRedistStatePrefix        = "redist_state"
	TagRedistStateSummary       = "redist_state_summary"
	TagRedistTgtStatePrefix     = "redist_tgtstate"
	TagRedistShardMoveSeparator = "|"

	TagRedistStateBegin          = "begin"
	TagRedistStateInprogress     = "in-progress"
	TagRedistStateFinishSnapshot = "finish_snapshot"
	TagRedistStateFinish         = "finish"

	// vaules
	TagRedistEnabledReady    = "ready"
	TagRedistEnabledSource   = "yes_source"
	TagRedistEnabledSourceRL = "yes_source_rl"
	TagRedistEnabledTarget   = "yes_target"
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

func KeyRedistEnable(zone int) string {
	return Key(TagRedistEnablePrefix, zone)
}

func KeyRedistFromNodeByZone(zone int) string {
	return Key(TagRedistFromNode, zone)
}

func KeyRedistNodeState(zone int, node int, shardid int) string {
	return Key(TagRedistStatePrefix, zone, node, shardid)
}
func KeyRedistFromNode(zone int, node int) string {
	return Key(TagRedistFromNode, zone, node)
}

func KeyRedistTgtNodeState(zone int, node int) string {
	return Key(TagRedistTgtStatePrefix, zone, node)
}
