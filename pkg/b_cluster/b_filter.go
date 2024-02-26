package cluster

type Filter struct {
	shardPos []map[uint32]int
	base     int
}

func (f *Filter) NumZones() int {
	return len(f.shardPos)
}
func (f *Filter) Get(zoneid int, shardid uint32) (nodeid int) {
	nodeid = f.shardPos[zoneid][shardid]
	return nodeid
}
func (f *Filter) set(zoneid int, shardid uint32, nodeid int) {
	f.shardPos[zoneid][shardid] = nodeid
}

func (f *Filter) inRange(nodeid int) bool {
	if f.base > 0 && nodeid >= f.base {
		return true
	}

	return false
}

// Select a shard to move to target node.
// Return the index in the shards.
func (f *Filter) selectShardForMove(shards []uint32, target Node) int {
	var ix = 0
	const scale = 1000
	var min = 1000 * scale

	if len(shards) == 1 {
		f.set(int(target.Zoneid), shards[0], int(target.Nodeid))
		return 0
	}

	numZones := f.NumZones()
	for i, shardid := range shards {

		var max = -1
		var count = 0

		for j := 0; j < numZones; j++ {

			if j == int(target.Zoneid) { // Skip my zone
				continue
			}

			nodeid := f.Get(j, shardid)

			if nodeid > max {
				max = nodeid
			}

			// Track count
			if f.inRange(nodeid) {
				if IsPrimary(shardid, uint32(j), uint32(numZones)) {
					count += 10
				} else {
					count += 9
				}
			}

		}

		score := count*scale + max
		if score >= 0 && score < min {
			min = score
			ix = i
		}
	}

	f.set(int(target.Zoneid), shards[ix], int(target.Nodeid))

	return ix
}
