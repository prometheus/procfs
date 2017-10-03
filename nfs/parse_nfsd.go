// /proc/net/rpc/nfsd parsing documented by https://www.svennd.be/nfsd-stats-explained-procnetrpcnfsd/
package nfs

import (
	"bufio"
	"fmt"
	"io"
)

// NewNFSdRPCStats returns stats read from /proc/net/rpc/nfsd
func (fs FS) NewNFSdRPCStats() (NFSdRPCStats, err) {
	var values []uint64

	f, err := os.Open(fs.Path("net/rpc/nfsd"))
	if err != nil {
		return Stat{}, err
	}
	defer f.Close()

	NFSdRPCStats := NFSdRPCStats{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(scanner.Text())
		// require at least <key> <value>
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid NFSd metric line %q", line)
		}
		label := parts[0]

		if label == "th" {
			values = procfs.ParseUint64s(parts[1:3])
		} else {
			values = procfs.ParseUint64s(parts[1:])
		}

		switch metricLine := parts[0]; metricLine {
		case "rc":
			replyCache, err := parseNFSdReplyCache(parts[1:])
		case "fh":
		case "io":
		case "th":
		case "ra":
		case "rpc":
		case "proc2":
		case "proc3":
		case "proc4":
		case "proc4ops":
		default:
			return nil, fmt.Errorf("invalid NFSd metric line %q", metricLine)
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing NFSdReplyCache: %s", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return Stat{}, fmt.Errorf("couldn't parse %s: %s", f.Name(), err)
	}

	return NFSdRPCStats, nil
}

func parseNFSdReplyCache(v []uint64) (NFSdReplyCache, err) {
	if len(v) != 3 {
		return nil, fmt.Errorf("invalid NFSdReplyCache line %q", v)
	}

	return NFSdReplyCache{
		Hits:    v[0],
		Misses:  v[1],
		NoCache: v[2],
	}, nil
}

func parseNFSdFileHandles(v []uint64) (NFSdFileHandles, err) {
	if len(v) != 5 {
		return nil, fmt.Errorf("invalid NFSdFileHandles, line %q", v)
	}

	return NFSdFileHandles{
		Stale:        v[0],
		TotalLookups: v[1],
		AnonLookups:  v[2],
		DirNoCache:   v[3],
		NoDirNoCache: v[4],
	}, nil
}

func parseNFSdInputOutput(v []uint64) (NFSdInputOutput, err) {
	if len(v) != 2 {
		return nil, fmt.Errorf("invalid NFSdInputOutput line %q", v)
	}

	return NFSdInputOutput{
		Read:  v[0],
		Write: v[1],
	}, nil
}

func parseNFSdThreads(v []uint64) (NFSdThreads, err) {
	if len(v) != 2 {
		return nil, fmt.Errorf("invalid NFSdThreads line %q", v)
	}

	return NFSdThreads{
		Threads: v[0],
		FullCnt: v[1],
	}, nil
}

func parseNFSdReadAheadCache(v []uint64) (NFSdReadAheadCache, err) {
	if len(v) != 12 {
		return nil, fmt.Errorf("invalid NFSdReadAheadCache line %q", v)
	}

	return NFSdReadAheadCache{
		CacheSize:      v[0],
		CacheHistogram: v[1:11],
		NotFound:       v[11],
	}, nil
}

func parseNFSdNetwork(v []uint64) (NFSdNetwork, err) {
	if len(v) != 4 {
		return nil, fmt.Errorf("invalid NFSdNetwork line %q", v)
	}

	return NFSdNetwork{
		NetCount:   v[0],
		UDPCount:   v[1],
		TCPCount:   v[2],
		TCPConnect: v[3],
	}, nil
}

func parseNFSdRPC(v []uint64) (NFSdRPC, err) {
	if len(v) != 5 {
		return nil, fmt.Errorf("invalid NFSdRPC line %q", v)
	}

	return NFSdRPC{
		RPCCount: v[0],
		BadCnt:   v[1],
		BadFmt:   v[2],
		BadAuth:  v[3],
		BadcInt:  v[4],
	}, nil
}

func parseNFSdv2Stats(v []uint64) (NFSdv2Stats, err) {
	values := v[0]
	if len(v) != values || values != 18 {
		return nil, fmt.Errorf("invalid NFSdv2Stats line %q", v)
	}

	return NFSdv2Stats{
		Null:     v[1],
		GetAttr:  v[2],
		SetAttr:  v[3],
		Root:     v[4],
		Lookup:   v[5],
		ReadLink: v[6],
		Read:     v[7],
		WrCache:  v[8],
		Write:    v[9],
		Create:   v[10],
		Remove:   v[11],
		Rename:   v[12],
		Link:     v[13],
		SymLink:  v[14],
		MkDir:    v[15],
		RmDir:    v[16],
		ReadDir:  v[17],
		FsStat:   v[18],
	}, nil
}

func parseNFSdv3Stats(v []uint64) (NFSdv3Stats, err) {
	values := v[0]
	if len(v) != values || values != 22 {
		return nil, fmt.Errorf("invalid NFSdv3Stats line %q", v)
	}

	return NFSdv3Stats{
		Null:        v[1],
		GetAttr:     v[2],
		SetAttr:     v[3],
		Lookup:      v[4],
		Access:      v[5],
		ReadLink:    v[6],
		Read:        v[7],
		Write:       v[8],
		Create:      v[9],
		MkDir:       v[10],
		SymLink:     v[11],
		MkNod:       v[12],
		Remove:      v[13],
		RmDir:       v[14],
		Rename:      v[15],
		Link:        v[16],
		ReadDir:     v[17],
		ReadDirPlus: v[18],
		FsStat:      v[19],
		FsInfo:      v[20],
		PathConf:    v[21],
		Commit:      v[22],
	}, nil
}

func parseNFSdv4Ops(v []uint64) (NFSdv4Ops, err) {
	values := v[0]
	if len(v) != values {
		return nil, fmt.Errorf("invalid NFSdv4Ops line %q", v)
	}

	stats := NFSdv4Ops{
		Op0Unused:    v[1],
		Op1Unused:    v[2],
		Op2Future:    v[3],
		Access:       v[4],
		Close:        v[5],
		Commit:       v[6],
		Create:       v[7],
		DelegPurge:   v[8],
		DelegReturn:  v[9],
		GetAttr:      v[10],
		GetFH:        v[11],
		Link:         v[12],
		Lock:         v[13],
		Lockt:        v[14],
		Locku:        v[15],
		Lookup:       v[16],
		LookupRoot:   v[17],
		Nverify:      v[18],
		Open:         v[19],
		OpenAttr:     v[20],
		OpenConfirm:  v[21],
		OpenDgrd:     v[22],
		PutFH:        v[23],
		PutPubFH:     v[24],
		PutRootFH:    v[25],
		Read:         v[26],
		ReadDir:      v[27],
		ReadLink:     v[28],
		Remove:       v[29],
		Rename:       v[31],
		Renew:        v[32],
		RestoreFH:    v[33],
		SaveFH:       v[34],
		SecInfo:      v[35],
		SetAttr:      v[36],
		Verify:       v[37],
		Write:        v[38],
		RelLockOwner: v[39],

	return stats, nil
}
