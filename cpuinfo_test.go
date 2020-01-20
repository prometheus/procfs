// Copyright 2019 The Prometheus Authors
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

package procfs

import "testing"

func TestCPUInfoBasic(t *testing.T) {
	cpuinfo, err := getProcFixtures(t).CPUInfo()
	if err != nil {
		t.Fatal(err)
	}

	if cpuinfo == nil {
		t.Fatal("cpuinfo is nil")
	}

	if want, have := 8, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}

	if want, have := uint(7), cpuinfo[7].Processor; want != have {
		t.Errorf("want processor %v, have %v", want, have)
	}
	if want, have := "GenuineIntel", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "6", cpuinfo[1].CPUFamily; want != have {
		t.Errorf("want family %v, have %v", want, have)
	}
	if want, have := "142", cpuinfo[2].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz", cpuinfo[3].ModelName; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := uint(8), cpuinfo[4].Siblings; want != have {
		t.Errorf("want siblings %v, have %v", want, have)
	}
	if want, have := "1", cpuinfo[5].CoreID; want != have {
		t.Errorf("want core id %v, have %v", want, have)
	}
	if want, have := uint(4), cpuinfo[6].CPUCores; want != have {
		t.Errorf("want cpu cores %v, have %v", want, have)
	}
	if want, have := "vme", cpuinfo[7].Flags[1]; want != have {
		t.Errorf("want flag %v, have %v", want, have)
	}
}

func TestCPUInfoDetectFormat(t *testing.T) {
	cpuinfoX86Bytes := []byte("processor  : 0\nvendor_id : GenuineIntel")
	if want, have := platformX86, cpuinfoDetectFormat(cpuinfoX86Bytes); want != have {
		t.Errorf("want cpuinfo format %v, have %v", want, have)
	}
	cpuinfoArm7Bytes := []byte(cpuinfoArm7)
	if want, have := platformARM, cpuinfoDetectFormat(cpuinfoArm7Bytes); want != have {
		t.Errorf("want cpuinfo format %v, have %v", want, have)
	}
	cpuinfoS390xBytes := []byte(cpuinfoS390x)
	if want, have := platformS390X, cpuinfoDetectFormat(cpuinfoS390xBytes); want != have {
		t.Errorf("want cpuinfo format %v, have %v", want, have)
	}
	cpuinfoPpc64Bytes := []byte(cpuinfoPpc64)
	if want, have := platformPPC, cpuinfoDetectFormat(cpuinfoPpc64Bytes); want != have {
		t.Errorf("want cpuinfo format %v, have %v", want, have)
	}
}

func TestCPUInfoParseARM(t *testing.T) {
	cpuinfo, err := parseCPUInfo([]byte(cpuinfoArm7))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse arm cpu info: %v", err)
	}
	if want, have := 2, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "ARMv7 Processor rev 5 (v7l)", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "thumb", cpuinfo[1].Flags[2]; want != have {
		t.Errorf("want flag %v, have %v", want, have)
	}
}

func TestCPUInfoParseS390X(t *testing.T) {
	cpuinfo, err := parseCPUInfo([]byte(cpuinfoS390x))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse s390x cpu info: %v", err)
	}
	if want, have := 4, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "IBM/S390", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "ldisp", cpuinfo[1].Flags[4]; want != have {
		t.Errorf("want flag %v, have %v", want, have)
	}
	if want, have := 5000.0, cpuinfo[2].CPUMHz; want != have {
		t.Errorf("want cpu MHz %v, have %v", want, have)
	}
}

func TestCPUInfoParsePPC(t *testing.T) {
	cpuinfo, err := parseCPUInfo([]byte(cpuinfoPpc64))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse ppc cpu info: %v", err)
	}
	if want, have := 6, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := 3000.00, cpuinfo[2].CPUMHz; want != have {
		t.Errorf("want cpu mhz %v, have %v", want, have)
	}
}

const (
	cpuinfoArm7 = `
Processor : ARMv7 Processor rev 5 (v7l)
processor : 0
BogoMIPS : 2400.00

processor : 1
BogoMIPS : 2400.00

Features : swp half thumb fastmult vfp edsp thumbee neon vfpv3 tls vfpv4 idiva idivt
CPU implementer : 0x41
CPU architecture: 7
CPU variant : 0x0
CPU part : 0xc07
CPU revision : 5

Hardware : sun8i
Revision : 0000
Serial : 5400503583203c3c040e`

	cpuinfoS390x = `
vendor_id       : IBM/S390
# processors    : 4
bogomips per cpu: 3033.00
max thread id   : 0
features	: esan3 zarch stfle msa ldisp eimm dfp edat etf3eh highgprs te vx sie
facilities      : 0 1 2 3 4 6 7 8 9 10 12 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 30 31 32 33 34 35 36 37 40 41 42 43 44 45 46 47 48 49 50 51 52 53 55 57 73 74 75 76 77 80 81 82 128 129 131
cache0          : level=1 type=Data scope=Private size=128K line_size=256 associativity=8
cache1          : level=1 type=Instruction scope=Private size=96K line_size=256 associativity=6
cache2          : level=2 type=Data scope=Private size=2048K line_size=256 associativity=8
cache3          : level=2 type=Instruction scope=Private size=2048K line_size=256 associativity=8
cache4          : level=3 type=Unified scope=Shared size=65536K line_size=256 associativity=16
cache5          : level=4 type=Unified scope=Shared size=491520K line_size=256 associativity=30
processor 0: version = FF,  identification = 2733E8,  machine = 2964
processor 1: version = FF,  identification = 2733E8,  machine = 2964
processor 2: version = FF,  identification = 2733E8,  machine = 2964
processor 3: version = FF,  identification = 2733E8,  machine = 2964

cpu number      : 0
cpu MHz dynamic : 5000
cpu MHz static  : 5000

cpu number      : 1
cpu MHz dynamic : 5000
cpu MHz static  : 5000

cpu number      : 2
cpu MHz dynamic : 5000
cpu MHz static  : 5000

cpu number      : 3
cpu MHz dynamic : 5000
cpu MHz static  : 5000
`

	cpuinfoPpc64 = `
processor	: 0
cpu		: POWER7 (architected), altivec supported
clock		: 3000.000000MHz
revision	: 2.1 (pvr 003f 0201)

processor	: 1
cpu		: POWER7 (architected), altivec supported
clock		: 3000.000000MHz
revision	: 2.1 (pvr 003f 0201)

processor	: 2
cpu		: POWER7 (architected), altivec supported
clock		: 3000.000000MHz
revision	: 2.1 (pvr 003f 0201)

processor	: 3
cpu		: POWER7 (architected), altivec supported
clock		: 3000.000000MHz
revision	: 2.1 (pvr 003f 0201)

processor	: 4
cpu		: POWER7 (architected), altivec supported
clock		: 3000.000000MHz
revision	: 2.1 (pvr 003f 0201)

processor	: 5
cpu		: POWER7 (architected), altivec supported
clock		: 3000.000000MHz
revision	: 2.1 (pvr 003f 0201)

timebase	: 512000000
platform	: pSeries
model		: IBM,8233-E8B
machine		: CHRP IBM,8233-E8B
`
)
