# Copyright The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include Makefile.common

%/.unpacked: %.ttar
	@echo ">> extracting fixtures $*"
	./ttar -C $(dir $*) -x -f $*.ttar
	touch $@

fixture_list := testdata/fixtures/proc/.unpacked \
  testdata/fixtures/sys/block/.unpacked \
  testdata/fixtures/sys/bus/.unpacked \
  testdata/fixtures/sys/class/.unpacked \
  testdata/fixtures/sys/devices/.unpacked \
  testdata/fixtures/sys/fs/.unpacked \
  testdata/fixtures/sys/kernel/.unpacked


fixtures: $(fixture_list)

update_fixtures:
	rm -vf testdata/fixtures/proc/.unpacked
	./ttar -c -f testdata/fixtures/proc.ttar -C testdata/fixtures proc/
	rm -vf testdata/fixtures/sys/block/.unpacked
	./ttar -c -f testdata/fixtures/sys/block.ttar -C testdata/fixtures/sys block/
	rm -vf testdata/fixtures/sys/bus/.unpacked
	./ttar -c -f testdata/fixtures/sys/bus.ttar -C testdata/fixtures/sys bus/
	rm -vf testdata/fixtures/sys/class/.unpacked
	./ttar -c -f testdata/fixtures/sys/class.ttar -C testdata/fixtures/sys class/
	rm -vf testdata/fixtures/sys/devices/.unpacked
	./ttar -c -f testdata/fixtures/sys/devices.ttar -C testdata/fixtures/sys devices/
	rm -vf testdata/fixtures/sys/fs/.unpacked
	./ttar -c -f testdata/fixtures/sys/fs.ttar -C testdata/fixtures/sys fs/
	rm -vf testdata/fixtures/sys/kernel/.unpacked
	./ttar -c -f testdata/fixtures/sys/kernel.ttar -C testdata/fixtures/sys kernel/

.PHONY: build
build:

.PHONY: test
test: $(fixture_list) common-test
