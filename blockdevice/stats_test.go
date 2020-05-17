// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blockdevice

import (
	"testing"
)

const (
	failMsgFormat  = "%v, expected %v, actual %v"
	procfsFixtures = "../fixtures/proc"
	sysfsFixtures  = "../fixtures/sys"
)

func TestDiskstats(t *testing.T) {
	blockdevice, err := NewFS(procfsFixtures, sysfsFixtures)
	if err != nil {
		t.Fatalf("failed to access blockdevice fs: %v", err)
	}
	diskstats, err := blockdevice.ProcDiskstats()
	if err != nil {
		t.Fatal(err)
	}
	expectedNumOfDevices := 52
	if len(diskstats) != expectedNumOfDevices {
		t.Errorf(failMsgFormat, "Incorrect number of devices", expectedNumOfDevices, len(diskstats))
	}
	if diskstats[0].DeviceName != "ram0" {
		t.Errorf(failMsgFormat, "Incorrect device name", "ram0", diskstats[0].DeviceName)
	}
	if diskstats[1].IoStatsCount != 14 {
		t.Errorf(failMsgFormat, "Incorrect number of stats read", 14, diskstats[0].IoStatsCount)
	}
	if diskstats[24].WriteIOs != 28444756 {
		t.Errorf(failMsgFormat, "Incorrect writes completed", 28444756, diskstats[24].WriteIOs)
	}
	if diskstats[48].DiscardTicks != 11130 {
		t.Errorf(failMsgFormat, "Incorrect discard time", 11130, diskstats[48].DiscardTicks)
	}
	if diskstats[48].IoStatsCount != 18 {
		t.Errorf(failMsgFormat, "Incorrect number of stats read", 18, diskstats[48].IoStatsCount)
	}
	if diskstats[49].IoStatsCount != 20 {
		t.Errorf(failMsgFormat, "Incorrect number of stats read", 20, diskstats[50].IoStatsCount)
	}
	if diskstats[49].FlushRequestsCompleted != 127 {
		t.Errorf(failMsgFormat, "Incorrect number of flash requests completed", 127, diskstats[50].FlushRequestsCompleted)
	}
	if diskstats[49].TimeSpentFlushing != 182 {
		t.Errorf(failMsgFormat, "Incorrect time spend flushing", 182, diskstats[50].TimeSpentFlushing)
	}
}

func TestBlockDevice(t *testing.T) {
	blockdevice, err := NewFS("../fixtures/proc", "../fixtures/sys")
	if err != nil {
		t.Fatalf("failed to access blockdevice fs: %v", err)
	}
	devices, err := blockdevice.SysBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	expectedNumOfDevices := 2
	if len(devices) != expectedNumOfDevices {
		t.Fatalf(failMsgFormat, "Incorrect number of devices", expectedNumOfDevices, len(devices))
	}
	if devices[0] != "dm-0" {
		t.Errorf(failMsgFormat, "Incorrect device name", "dm-0", devices[0])
	}
	device0stats, count, err := blockdevice.SysBlockDeviceStat(devices[0])
	if err != nil {
		t.Fatal(err)
	}
	if count != 11 {
		t.Errorf(failMsgFormat, "Incorrect number of stats read", 11, count)
	}
	if device0stats.ReadIOs != 6447303 {
		t.Errorf(failMsgFormat, "Incorrect read I/Os", 6447303, device0stats.ReadIOs)
	}
	if device0stats.WeightedIOTicks != 6088971 {
		t.Errorf(failMsgFormat, "Incorrect time in queue", 6088971, device0stats.WeightedIOTicks)
	}
	device1stats, count, err := blockdevice.SysBlockDeviceStat(devices[1])
	if count != 15 {
		t.Errorf(failMsgFormat, "Incorrect number of stats read", 15, count)
	}
	if err != nil {
		t.Fatal(err)
	}
	if device1stats.WriteSectors != 286915323 {
		t.Errorf(failMsgFormat, "Incorrect write merges", 286915323, device1stats.WriteSectors)
	}
	if device1stats.DiscardTicks != 12 {
		t.Errorf(failMsgFormat, "Incorrect discard ticks", 12, device1stats.DiscardTicks)
	}
	blockQueueStat, err := blockdevice.SysBlockDeviceQueueStats(devices[1])
	if err != nil {
		t.Fatal(err)
	}
	if blockQueueStat.AddRandom != 1 {
		t.Errorf(failMsgFormat, "Incorrect add_random", 1, blockQueueStat.AddRandom)
	}
	if blockQueueStat.ChunkSectors != 0 {
		t.Errorf(failMsgFormat, "Incorrect chunk_sectors", 0, blockQueueStat.ChunkSectors)
	}
	if blockQueueStat.DAX != 0 {
		t.Errorf(failMsgFormat, "Incorrect dax", 0, blockQueueStat.DAX)
	}
	if blockQueueStat.DiscardGranularity != 0 {
		t.Errorf(failMsgFormat, "Incorrect discard_granularity", 0, blockQueueStat.DiscardGranularity)
	}
	if blockQueueStat.DiscardMaxHWBytes != 0 {
		t.Errorf(failMsgFormat, "Incorrect discard_max_hw_bytes", 0, blockQueueStat.DiscardMaxHWBytes)
	}
	if blockQueueStat.DiscardMaxBytes != 0 {
		t.Errorf(failMsgFormat, "Incorrect discard_max_bytes", 0, blockQueueStat.DiscardMaxBytes)
	}
	if blockQueueStat.FUA != 0 {
		t.Errorf(failMsgFormat, "Incorrect fua", 0, blockQueueStat.FUA)
	}
	if blockQueueStat.HWSectorSize != 512 {
		t.Errorf(failMsgFormat, "Incorrect hw_sector_size", 512, blockQueueStat.HWSectorSize)
	}
	if blockQueueStat.IOPoll != 0 {
		t.Errorf(failMsgFormat, "Incorrect io_poll", 0, blockQueueStat.IOPoll)
	}
	if blockQueueStat.IOPollDelay != -1 {
		t.Errorf(failMsgFormat, "Incorrect io_poll_delay", -1, blockQueueStat.IOPollDelay)
	}
	if blockQueueStat.IOTimeout != 30000 {
		t.Errorf(failMsgFormat, "Incorrect io_timeout", 30000, blockQueueStat.IOTimeout)
	}
	if blockQueueStat.IOStats != 1 {
		t.Errorf(failMsgFormat, "Incorrect iostats", 1, blockQueueStat.IOStats)
	}
	if blockQueueStat.LogicalBlockSize != 512 {
		t.Errorf(failMsgFormat, "Incorrect logical_block_size", 512, blockQueueStat.LogicalBlockSize)
	}
	if blockQueueStat.MaxDiscardSegments != 1 {
		t.Errorf(failMsgFormat, "Incorrect max_discard_segments", 1, blockQueueStat.MaxDiscardSegments)
	}
	if blockQueueStat.MaxHWSectorsKB != 32767 {
		t.Errorf(failMsgFormat, "Incorrect max_hw_sectors_kb", 32767, blockQueueStat.MaxHWSectorsKB)
	}
	if blockQueueStat.MaxIntegritySegments != 0 {
		t.Errorf(failMsgFormat, "Incorrect max_integrity_segments", 0, blockQueueStat.MaxIntegritySegments)
	}
	if blockQueueStat.MaxSectorsKB != 1280 {
		t.Errorf(failMsgFormat, "Incorrect max_sectors_kb", 1280, blockQueueStat.MaxSectorsKB)
	}
	if blockQueueStat.MaxSegments != 168 {
		t.Errorf(failMsgFormat, "Incorrect max_segments", 168, blockQueueStat.MaxSegments)
	}
	if blockQueueStat.MaxSegmentSize != 65536 {
		t.Errorf(failMsgFormat, "Incorrect max_segment_size", 65536, blockQueueStat.MaxSegmentSize)
	}
	if blockQueueStat.MinimumIOSize != 512 {
		t.Errorf(failMsgFormat, "Incorrect minimum_io_size", 512, blockQueueStat.MinimumIOSize)
	}
	if blockQueueStat.NoMerges != 0 {
		t.Errorf(failMsgFormat, "Incorrect nomerges", 0, blockQueueStat.NoMerges)
	}
	if blockQueueStat.NRRequests != 64 {
		t.Errorf(failMsgFormat, "Incorrect nr_requests", 64, blockQueueStat.NRRequests)
	}
	if blockQueueStat.NRZones != 0 {
		t.Errorf(failMsgFormat, "Incorrect nr_zones", 0, blockQueueStat.NRZones)
	}
	if blockQueueStat.OptimalIOSize != 0 {
		t.Errorf(failMsgFormat, "Incorrect optimal_io_size", 0, blockQueueStat.OptimalIOSize)
	}
	if blockQueueStat.PhysicalBlockSize != 512 {
		t.Errorf(failMsgFormat, "Incorrect physical_block_size", 512, blockQueueStat.PhysicalBlockSize)
	}
	if blockQueueStat.ReadAHeadKB != 128 {
		t.Errorf(failMsgFormat, "Incorrect read_ahead_kb", 128, blockQueueStat.ReadAHeadKB)
	}
	if blockQueueStat.Rotational != 1 {
		t.Errorf(failMsgFormat, "Incorrect rotational", 1, blockQueueStat.Rotational)
	}
	if blockQueueStat.RQAffinity != 1 {
		t.Errorf(failMsgFormat, "Incorrect rq_affinity", 1, blockQueueStat.RQAffinity)
	}
	if blockQueueStat.SchedulerCurrent != "bfq" {
		t.Errorf(failMsgFormat, "Incorrect current scheduler", "bfq", blockQueueStat.SchedulerCurrent)
	}
	if len(blockQueueStat.SchedulerList) != 4 {
		t.Errorf(failMsgFormat, "Incorrect scheduler list", 4, len(blockQueueStat.SchedulerList))
	}
	if blockQueueStat.WriteCache != "write back" {
		t.Errorf(failMsgFormat, "Incorrect write_cache", "write back", blockQueueStat.WriteCache)
	}
	if blockQueueStat.WriteSameMaxBytes != 0 {
		t.Errorf(failMsgFormat, "Incorrect write_same_max_bytes", 0, blockQueueStat.WriteSameMaxBytes)
	}
	if blockQueueStat.WBTLatUSec != 75000 {
		t.Errorf(failMsgFormat, "Incorrect wbt_lat_usec", 75000, blockQueueStat.WBTLatUSec)
	}
	if blockQueueStat.ThrottleSampleTime != nil {
		t.Errorf(failMsgFormat, "Incorrect throttle_sample_time", nil, blockQueueStat.ThrottleSampleTime)
	}
	if blockQueueStat.WriteZeroesMaxBytes != 0 {
		t.Errorf(failMsgFormat, "Incorrect write_zeroes_max_bytes", 0, blockQueueStat.WriteZeroesMaxBytes)
	}
	if blockQueueStat.Zoned != "none" {
		t.Errorf(failMsgFormat, "Incorrect zoned", 0, blockQueueStat.Zoned)
	}
}
