// Copyright The Prometheus Authors
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

//go:build linux

package procfs

import "testing"

const (
	cpuinfoArm7Legacy = `
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

	cpuinfoArm7LegacyV1 = `
Processor       : ARMv6-compatible processor rev 5 (v6l)
BogoMIPS        : 791.34
Features        : swp half thumb fastmult vfp edsp java 
CPU implementer : 0x41
CPU architecture: 6TEJ
CPU variant     : 0x1
CPU part        : 0xb36
CPU revision    : 5

Hardware        : IMAPX200
Revision        : 0000
Serial          : 0000000000000000`

	cpuinfoArm7 = `
processor : 0
model name : ARMv7 Processor rev 3 (v7l)
BogoMIPS : 108.00
Features : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant : 0x0
CPU part : 0xd08
CPU revision : 3

processor : 1
model name : ARMv7 Processor rev 3 (v7l)
BogoMIPS : 108.00
Features : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant : 0x0
CPU part : 0xd08
CPU revision : 3

processor : 2
model name : ARMv7 Processor rev 3 (v7l)
BogoMIPS : 108.00
Features : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant : 0x0
CPU part : 0xd08
CPU revision : 3

processor : 3
model name : ARMv7 Processor rev 3 (v7l)
BogoMIPS : 108.00
Features : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32
CPU implementer : 0x41
CPU architecture: 7
CPU variant : 0x0
CPU part : 0xd08
CPU revision : 3

Hardware : BCM2835
Revision : c03111
`

	cpuinfoArm64 = `
processor	: 0
BogoMIPS	: 2000.00
Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm jscvt fcma lrcpc dcpop sha3 sm3 sm4 asimddp sha512 sve asimdfhm dit uscat ilrcpc flagm sb paca pacg dcpodp sve2 sveaes svepmull svebitperm svesha3 svesm4 flagm2 frint svei8mm svebf16 i8mm bf16 dgh rng bti
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x0
CPU part	: 0xd4f
CPU revision	: 1

processor	: 1
BogoMIPS	: 2000.00
Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm jscvt fcma lrcpc dcpop sha3 sm3 sm4 asimddp sha512 sve asimdfhm dit uscat ilrcpc flagm sb paca pacg dcpodp sve2 sveaes svepmull svebitperm svesha3 svesm4 flagm2 frint svei8mm svebf16 i8mm bf16 dgh rng bti
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x0
CPU part	: 0xd4f
CPU revision	: 1
`

	cpuinfoArm64Heterogeneous = `
processor	: 0
BogoMIPS	: 2000.00
Features	: fp asimd evtstrm aes
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x0
CPU part	: 0xd03
CPU revision	: 4

processor	: 1
BogoMIPS	: 2000.00
Features	: fp asimd evtstrm aes
CPU implementer	: 0x48
CPU architecture: 8
CPU variant	: 0x1
CPU part	: 0xd40
CPU revision	: 2
`

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
physical id     : 2
core id         : 0
siblings        : 8
cpu cores       : 4
cpu MHz dynamic : 5000
cpu MHz static  : 5000

cpu number      : 1
physical id     : 2
core id         : 0
siblings        : 8
cpu cores       : 4
cpu MHz dynamic : 5000
cpu MHz static  : 5000

cpu number      : 2
physical id     : 2
core id         : 1
siblings        : 8
cpu cores       : 4
cpu MHz dynamic : 5000
cpu MHz static  : 5000

cpu number      : 3
physical id     : 2
core id         : 1
siblings        : 8
cpu cores       : 4
cpu MHz dynamic : 5000
cpu MHz static  : 5000
`

	cpuinfoMips = `
system type		: UBNT_E100
machine			: Unknown
processor		: 0
cpu model		: Cavium Octeon+ V0.1
BogoMIPS		: 1000.00
wait instruction	: yes
microsecond timers	: yes
tlb_entries		: 64
extra interrupt vector	: yes
hardware watchpoint	: yes, count: 2, address/irw mask: [0x0ffc, 0x0ffb]
isa			: mips1 mips2 mips3 mips4 mips5 mips64r2
ASEs implemented	:
shadow register sets	: 1
kscratch registers	: 0
core			: 0
VCED exceptions		: not available
VCEI exceptions		: not available

processor		: 1
cpu model		: Cavium Octeon+ V0.1
BogoMIPS		: 1000.00
wait instruction	: yes
microsecond timers	: yes
tlb_entries		: 64
extra interrupt vector	: yes
hardware watchpoint	: yes, count: 2, address/irw mask: [0x0ffc, 0x0ffb]
isa			: mips1 mips2 mips3 mips4 mips5 mips64r2
ASEs implemented	:
shadow register sets	: 1
kscratch registers	: 0
core			: 1
VCED exceptions		: not available
VCEI exceptions		: not available

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

	cpuinfoRiscv64 = `
processor	: 0
hart		: 0
isa		: rv64imafdcsu
mmu		: sv48

processor	: 1
hart		: 1
isa		: rv64imafdcsu
mmu		: sv48
`

	cpuinfoLoong64 = `
system type		: generic-loongson-machine

processor		: 0
package			: 0
core			: 0
CPU Family		: Loongson-64bit
Model Name		: Loongson-3A5000
CPU Revision		: 0x10
FPU Revision		: 0x00
CPU MHz			: 2500.00
BogoMIPS		: 5000.00
TLB Entries		: 2112
Address Sizes		: 48 bits physical, 48 bits virtual
ISA			: loongarch32 loongarch64
Features		: cpucfg lam ual fpu complex crypto lvz
Hardware Watchpoint	: yes, iwatch count: 0, dwatch count: 0

processor		: 1
package			: 0
core			: 1
CPU Family		: Loongson-64bit
Model Name		: Loongson-3A5000
CPU Revision		: 0x10
FPU Revision		: 0x00
CPU MHz			: 2500.00
BogoMIPS		: 5000.00
TLB Entries		: 2112
Address Sizes		: 48 bits physical, 48 bits virtual
ISA			: loongarch32 loongarch64
Features		: cpucfg lam ual fpu complex crypto lvz
Hardware Watchpoint	: yes, iwatch count: 0, dwatch count: 0

processor		: 2
package			: 0
core			: 2
CPU Family		: Loongson-64bit
Model Name		: Loongson-3A5000
CPU Revision		: 0x10
FPU Revision		: 0x00
CPU MHz			: 2500.00
BogoMIPS		: 5000.00
TLB Entries		: 2112
Address Sizes		: 48 bits physical, 48 bits virtual
ISA			: loongarch32 loongarch64
Features		: cpucfg lam ual fpu complex crypto lvz
Hardware Watchpoint	: yes, iwatch count: 0, dwatch count: 0

processor		: 3
package			: 0
core			: 3
CPU Family		: Loongson-64bit
Model Name		: Loongson-3A5000
CPU Revision		: 0x10
FPU Revision		: 0x00
CPU MHz			: 2500.00
BogoMIPS		: 5000.00
TLB Entries		: 2112
Address Sizes		: 48 bits physical, 48 bits virtual
ISA			: loongarch32 loongarch64
Features		: cpucfg lam ual fpu complex crypto lvz
Hardware Watchpoint	: yes, iwatch count: 0, dwatch count: 0
`
)

func TestCPUInfoX86(t *testing.T) {
	parseCPUInfo = parseCPUInfoX86
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

func TestCPUInfoParseARMLegacy(t *testing.T) {
	cpuinfo, err := parseCPUInfoARM([]byte(cpuinfoArm7Legacy))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse arm cpu info: %v", err)
	}
	if want, have := 2, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "ARMv7 Processor rev 5 (v7l)", cpuinfo[0].ModelName; want != have {
		t.Errorf("want model name %v, have %v", want, have)
	}
	if want, have := "0x41", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "7", cpuinfo[0].CPUFamily; want != have {
		t.Errorf("want family %v, have %v", want, have)
	}
	if want, have := "0xc07", cpuinfo[0].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "r0p5", cpuinfo[0].Stepping; want != have {
		t.Errorf("want stepping %v, have %v", want, have)
	}
	if want, have := "thumb", cpuinfo[1].Flags[2]; want != have {
		t.Errorf("want flag %v, have %v", want, have)
	}
	if want, have := "0x41", cpuinfo[1].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "0xc07", cpuinfo[1].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "r0p5", cpuinfo[1].Stepping; want != have {
		t.Errorf("want stepping %v, have %v", want, have)
	}
}

func TestCPUInfoParseARMLegacyV1(t *testing.T) {
	cpuinfo, err := parseCPUInfoARM([]byte(cpuinfoArm7LegacyV1))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse arm cpu info: %v", err)
	}
	if want, have := 1, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "ARMv6-compatible processor rev 5 (v6l)", cpuinfo[0].ModelName; want != have {
		t.Errorf("want model name %v, have %v", want, have)
	}
	if want, have := "0x41", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "6TEJ", cpuinfo[0].CPUFamily; want != have {
		t.Errorf("want family %v, have %v", want, have)
	}
	if want, have := "0xb36", cpuinfo[0].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "r1p5", cpuinfo[0].Stepping; want != have {
		t.Errorf("want stepping %v, have %v", want, have)
	}
	if want, have := "thumb", cpuinfo[0].Flags[2]; want != have {
		t.Errorf("want flag %v, have %v", want, have)
	}
}

func TestCPUInfoParseARM(t *testing.T) {
	cpuinfo, err := parseCPUInfoARM([]byte(cpuinfoArm7))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse arm cpu info: %v", err)
	}
	if want, have := 4, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "ARMv7 Processor rev 3 (v7l)", cpuinfo[0].ModelName; want != have {
		t.Errorf("want model name %v, have %v", want, have)
	}
	if want, have := "0x41", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "7", cpuinfo[0].CPUFamily; want != have {
		t.Errorf("want family %v, have %v", want, have)
	}
	if want, have := "0xd08", cpuinfo[0].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "r0p3", cpuinfo[0].Stepping; want != have {
		t.Errorf("want stepping %v, have %v", want, have)
	}
	if want, have := "thumb", cpuinfo[1].Flags[1]; want != have {
		t.Errorf("want flag %v, have %v", want, have)
	}
}

func TestCPUInfoParseARM64(t *testing.T) {
	cpuinfo, err := parseCPUInfoARM([]byte(cpuinfoArm64))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse arm64 cpu info: %v", err)
	}
	if want, have := 2, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	for i := range cpuinfo {
		if want, have := "0x41", cpuinfo[i].VendorID; want != have {
			t.Errorf("cpu %d: want vendor %v, have %v", i, want, have)
		}
		if want, have := "8", cpuinfo[i].CPUFamily; want != have {
			t.Errorf("cpu %d: want family %v, have %v", i, want, have)
		}
		if want, have := "0xd4f", cpuinfo[i].Model; want != have {
			t.Errorf("cpu %d: want model %v, have %v", i, want, have)
		}
		if want, have := "r0p1", cpuinfo[i].Stepping; want != have {
			t.Errorf("cpu %d: want stepping %v, have %v", i, want, have)
		}
		if want, have := 2000.00, cpuinfo[i].BogoMips; want != have {
			t.Errorf("cpu %d: want bogomips %v, have %v", i, want, have)
		}
		if want, have := "fp", cpuinfo[i].Flags[0]; want != have {
			t.Errorf("cpu %d: want flag %v, have %v", i, want, have)
		}
	}
}

func TestCPUInfoParseARM64Heterogeneous(t *testing.T) {
	cpuinfo, err := parseCPUInfoARM([]byte(cpuinfoArm64Heterogeneous))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse heterogeneous arm64 cpu info: %v", err)
	}
	if want, have := 2, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "0x41", cpuinfo[0].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "8", cpuinfo[0].CPUFamily; want != have {
		t.Errorf("want family %v, have %v", want, have)
	}
	if want, have := "0xd03", cpuinfo[0].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "r0p4", cpuinfo[0].Stepping; want != have {
		t.Errorf("want stepping %v, have %v", want, have)
	}

	if want, have := "0x48", cpuinfo[1].VendorID; want != have {
		t.Errorf("want vendor %v, have %v", want, have)
	}
	if want, have := "8", cpuinfo[1].CPUFamily; want != have {
		t.Errorf("want family %v, have %v", want, have)
	}
	if want, have := "0xd40", cpuinfo[1].Model; want != have {
		t.Errorf("want model %v, have %v", want, have)
	}
	if want, have := "r1p2", cpuinfo[1].Stepping; want != have {
		t.Errorf("want stepping %v, have %v", want, have)
	}
}

func TestCPUInfoParseS390X(t *testing.T) {
	cpuinfo, err := parseCPUInfoS390X([]byte(cpuinfoS390x))
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
	if want, have := uint(8), cpuinfo[3].Siblings; want != have {
		t.Errorf("want siblings %v, have %v", want, have)
	}
	if want, have := "1", cpuinfo[3].CoreID; want != have {
		t.Errorf("want core id %v, have %v", want, have)
	}
	if want, have := uint(4), cpuinfo[2].CPUCores; want != have {
		t.Errorf("want cpu cores %v, have %v", want, have)
	}
	if want, have := "2", cpuinfo[2].PhysicalID; want != have {
		t.Errorf("want physical id %v, have %v", want, have)
	}
}

func TestCPUInfoParseMips(t *testing.T) {
	cpuinfo, err := parseCPUInfoMips([]byte(cpuinfoMips))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse mips cpu info: %v", err)
	}
	if want, have := 2, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := 1000.00, cpuinfo[0].BogoMips; want != have {
		t.Errorf("want BogoMIPS %v, have %v", want, have)
	}
	if want, have := "Cavium Octeon+ V0.1", cpuinfo[1].ModelName; want != have {
		t.Errorf("want ModelName '%v', have '%v'", want, have)
	}
}

func TestCPUInfoParsePPC(t *testing.T) {
	cpuinfo, err := parseCPUInfoPPC([]byte(cpuinfoPpc64))
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

func TestCPUInfoParseRISCV64(t *testing.T) {
	cpuinfo, err := parseCPUInfoRISCV([]byte(cpuinfoRiscv64))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse ppc cpu info: %v", err)
	}
	if want, have := 2, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "1", cpuinfo[1].CoreID; want != have {
		t.Errorf("want CoreId %v, have %v", want, have)
	}
	if want, have := "rv64imafdcsu", cpuinfo[1].ModelName; want != have {
		t.Errorf("want ModelName %v, have %v", want, have)
	}
}

func TestCPUInfoParseLoong64(t *testing.T) {
	cpuinfo, err := parseCPUInfoLoong([]byte(cpuinfoLoong64))
	if err != nil || cpuinfo == nil {
		t.Fatalf("unable to parse loong cpu info: %v", err)
	}
	if want, have := 4, len(cpuinfo); want != have {
		t.Errorf("want number of processors %v, have %v", want, have)
	}
	if want, have := "Loongson-64bit", cpuinfo[1].CPUFamily; want != have {
		t.Errorf("want CPUFamily '%v', have '%v'", want, have)
	}
}

func TestFormatSteppingARM(t *testing.T) {
	tests := []struct {
		variant  string
		revision string
		expected string
	}{
		{"", "", ""},
		{"0x0", "1", "r0p1"},
		{"0x1", "5", "r1p5"},
		{"0x0", "3", "r0p3"},
		{"0x2", "0", "r2p0"},
		{"0", "1", "r0p1"},
		{"1", "2", "r1p2"},
		{"0x0", "", "r0"},
		{"", "1", "p1"},
		{"abc", "def", "rabcpdef"},
	}

	for _, tt := range tests {
		got := formatSteppingARM(tt.variant, tt.revision)
		if got != tt.expected {
			t.Errorf("formatSteppingARM(%q, %q) = %q, want %q", tt.variant, tt.revision, got, tt.expected)
		}
	}
}

func TestCPUInfoParseARMEdgeCases(t *testing.T) {
	input := `
processor	: 0
cpu implementer	: 0x41
cpu architecture: 8
cpu variant	: 0x0
cpu part	: 0xd03

processor	: 1
cpu implementer	: 0x41
cpu architecture: 8
cpu revision	: 2
cpu part	: 0xd03
`
	cpuinfo, err := parseCPUInfoARM([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cpuinfo) != 2 {
		t.Fatalf("expected 2 cpus, got %d", len(cpuinfo))
	}
	if cpuinfo[0].VendorID != "0x41" {
		t.Errorf("expected vendor 0x41, got %q", cpuinfo[0].VendorID)
	}
	if cpuinfo[0].CPUFamily != "8" {
		t.Errorf("expected family 8, got %q", cpuinfo[0].CPUFamily)
	}
	if cpuinfo[0].Model != "0xd03" {
		t.Errorf("expected model 0xd03, got %q", cpuinfo[0].Model)
	}
	if cpuinfo[0].Stepping != "r0" {
		t.Errorf("expected stepping r0, got %q", cpuinfo[0].Stepping)
	}

	if cpuinfo[1].Stepping != "p2" {
		t.Errorf("expected stepping p2, got %q", cpuinfo[1].Stepping)
	}

	_, err = parseCPUInfoARM([]byte("invalid line\n"))
	if err == nil {
		t.Errorf("expected error for malformed line, got nil")
	}
}

