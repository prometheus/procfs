// Copyright 2017 Alex Lau (AvengerMoJo) <alau@suse.com>
//
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


package iscsi

import (
    "fmt"
    "os"
    "path/filepath"
)

// infomation about iSCSI info from kernel/config/target/iscsi

const TARGET_PATH = "/sys/kernel/config/target/iscsi"

type TPGT struct {
    tpgt_path   string
    is_enable   bool
    luns_path   []string
}

type LUN struct {
    lun_path    string
    backstore   string
    object_name, type_number,
}

type Stats struct {
    Name    string
    tpgt    []TPGT
}

func GetStats(iqn_path string) (*Stats, error) {
    var iscsi Stats

    iscsi.Name := filepath.Base(iqn_path)
    matches, err := filepath.Glob(filepath.Join(iqn_path, iscsi.TARGET_PATH, "tpgt*"))
    if err != nil {
        fmt.Errorf("getTPGT error %v\n", err )
        return nil, fmt;
    }
    iscsi.tpgt := make([]string, len(matches))
    // iscsi.is_enable := make([]bool, len(matches))
    for ii, tpgt_path := range matches {
        iscsi.tpgt[ii].tpgt_path := filepath.Base(tpgt_path)
        iscsi.tpgt[ii].is_enable := isTpgtEnable(tpgt_path)
        if (iscsi.tpgt[ii].is_enable) {
                iscsi.tpgt[ii].luns_path, _:= getLun(tpgt_path)
                for _, lun_path := range iscsi.tpgt[ii].luns_path {
                    backstore_type, object_name, type_number, err := getLunLinkTarget(lun_path)
        }
    }
    return matches, nil
}

func isTpgtEnable(tpgt_path string) (isEnable bool, error) {
    isEnable = false
    tmp, err := ioutil.ReadFile(filepath.Join(tpgt_path, "enable"))
    log.Debugf("isTpgtenable %s\n", tmp)
    if err != nil {
        return false, fmt.Errorf("isTpgtEnable error %v\n", err)
    }
    tmp_num, err := strconv.Atoi(strings.TrimSpace(string(tmp)))
    log.Debugf("getTpgt enable number %d\n", tmp_num)
    if (tmp_num > 0) {
        isEnable = true
    }
    return isEnable
}

func getLun(tpgt_path string) (matches []string, err error) {
    log.Debugf("getLun path %s\n", tpgt_path)
    matches, err := filepath.Glob(filepath.Join(tpgt_path, "lun/lun*"))
    if err != nil {
        return nil, fmt.Errorf("getLun error  %v\n", err )
    }
    return matches, nil
}

func getLunLinkTarget(lun_path string) ( backstore_type string, object_name string, type_number string, err error) {
    files, err := ioutil.ReadDir(lun_path)
    if err != nil {
        fmt.Errorf("getLunLinkTarget error  %v\n", err )
        return "", "", "", nil
    }
    for _, file := range files {
        log.Debugf("lun dir list file ->%s<-\n",file.Name())
        fileInfo, _:= os.Lstat(lun_path +  "/" + file.Name())
        if fileInfo.Mode() & os.ModeSymlink != 0 {
            target, err := os.Readlink( lun_path +  "/" + fileInfo.Name())
            if err != nil {
                return "", "", "", fmt.Errorf("Readlink err %v\n", err)
            }
            p1, object_name := filepath.Split(target)
            _, type_with_number := filepath.Split(filepath.Clean(p1))

            tmp := strings.Split(type_with_number, "_")
            backstore_type, type_number := tmp[0], tmp[1]
            if len(tmp) == 3 {
                backstore_type = fmt.Sprintf("%s_%s", tmp[0], tmp[1])
                type_number = tmp[2]
            }
            log.Debugf("object_name->%s-<, type->%s, type_number->%s<- \n", object_name, backstore_type, type_number)
            return backstore_type, object_name, type_number, nil
        }
    }
    return "", "", "", errors.New("getLunLinkTarget: Lun Link does not exist")
}
