// /proc/net/rpc/nfsd parsing documented by https://www.svennd.be/nfsd-stats-explained-procnetrpcnfsd/
package nfs

// rc line: Reply Cache
type NFSdReplyCache struct {
	Hits    uint64
	Misses  uint64
	NoCache uint64
}

// fh line: File Handles
type NFSdFileHandles struct {
	Stale        uint64
	TotalLookups uint64
	AnonLookups  uint64
	DirNoCache   uint64
	NoDirNoCache uint64
}

// io line: Input Output
type NFSdInputOutput struct {
	Read  uint64
	Write uint64
}

// th line: Threads
type NFSdThreads struct {
	Threads uint64
	FullCnt uint64
}

// ra line: Read Ahead Cache
type NFSdReadAheadCache struct {
	CacheSize      uint64
	CacheHistogram [10]uint64
	NotFound       uint64
}

// net line: Network
type NFSdNetwork struct {
	NetCount   uint64
	UDPCount   uint64
	TCPCount   uint64
	TCPConnect uint64
}

// rpc line:
type NFSdRPC struct {
	RPCCount uint64
	BadCnt   uint64
	BadFmt   uint64
	BadAuth  uint64
	BadcInt  uint64
}

// proc2 line: NFSv2 Stats
type NFSdv2Stats struct {
	Values   uint64 // Should be 18.
	Null     uint64
	GetAttr  uint64
	SetAttr  uint64
	Root     uint64
	Lookup   uint64
	ReadLink uint64
	Read     uint64
	WrCache  uint64
	Write    uint64
	Create   uint64
	Remove   uint64
	Rename   uint64
	Link     uint64
	SymLink  uint64
	MkDir    uint64
	RmDir    uint64
	ReadDir  uint64
	FsStat   uint64
}

// proc3 line: NFSv3 Stats
type NFSdv3Stats struct {
	Values      uint64 // Should be 22.
	Null        uint64
	GetAttr     uint64
	SetAttr     uint64
	Lookup      uint64
	Access      uint64
	ReadLink    uint64
	Read        uint64
	Write       uint64
	Create      uint64
	MkDir       uint64
	SymLink     uint64
	MkNod       uint64
	Remove      uint64
	RmDir       uint64
	Rename      uint64
	Link        uint64
	ReadDir     uint64
	ReadDirPlus uint64
	FsStat      uint64
	FsInfo      uint64
	PathConf    uint64
	Commit      uint64
}

// proc4 line: NFSv4 Stats
type NFSdv4Stats struct {
	Values   uint64 // Should be 2.
	Null     uint64
	Compound uint64
}

// proc4ops line: NFSv4 operations
// Variable list, see:
// v4.0 https://tools.ietf.org/html/rfc3010 (38 operations)
// v4.1 https://tools.ietf.org/html/rfc5661 (58 operations)
// v4.2 https://tools.ietf.org/html/draft-ietf-nfsv4-minorversion2-41 (71 operations)
type NFSdv4Ops struct {
	Values       uint64 // Variable depending on v4.x sub-version.
	Op0Unused    uint64
	Op1Unused    uint64
	Op2Future    uint64
	Access       uint64
	Close        uint64
	Commit       uint64
	Create       uint64
	DelegPurge   uint64
	DelegReturn  uint64
	GetAttr      uint64
	GetFH        uint64
	Link         uint64
	Lock         uint64
	Lockt        uint64
	Locku        uint64
	Lookup       uint64
	LookupRoot   uint64
	Nverify      uint64
	Open         uint64
	OpenAttr     uint64
	OpenConfirm  uint64
	OpenDgrd     uint64
	PutFH        uint64
	PutPubFH     uint64
	PutRootFH    uint64
	Read         uint64
	ReadDir      uint64
	ReadLink     uint64
	Remove       uint64
	Rename       uint64
	Renew        uint64
	RestoreFH    uint64
	SaveFH       uint64
	SecInfo      uint64
	SetAttr      uint64
	Verify       uint64
	Write        uint64
	RelLockOwner uint64
}

// All stats from /proc/net/rpc/nfsd
type NFSdRPCStats struct {
	NFSdReplyCache     NFSdReplyCache
	NFSdFileHandles    NFSdFileHandles
	NFSdInputOutput    NFSdInputOutput
	NFSdThreads        NFSdThreads
	NFSdReadAheadCache NFSdReadAheadCache
	NFSdNetwork        NFSdNetwork
	NFSdRPC            NFSdRPC
	NFSdv2Stats        NFSdv2Stats
	NFSdv3Stats        NFSdv3Stats
	NFSdv4Stats        NFSdv4Stats
	NFSdv4Ops          NFSdv4Ops
	NFSdRPCStats       NFSdRPCStats
}
