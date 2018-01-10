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
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv" 
    "strings" 

    // "github.com/prometheus/common/log"
)

const TARGET_PATH = "/sys/kernel/config/target/iscsi"
const TARGET_CORE = "/sys/kernel/config/target/core"

type TPGT struct {
    Name        string
    Tpgt_path   string
    Is_enable   bool
    Luns        []LUN
}

type LUN struct {
    Name        string
    Lun_path    string
    Backstore   string
    Object_name string
    Type_number string
}

type FILEIO struct { 
    Name        string
    Fnumber     string
    Object_name string
    Filename    string
}

type IBLOCK struct {
    Name        string
    Bnumber     string
    Object_name string
    Iblock      string
}

type RBD struct {
    Name    string
    Rnumber string
    Pool    string
    Image   string
}

type RDMCP struct {
    Name        string
    Object_name string
}

type Stats struct {
    Name    string
    Tpgt    []TPGT
}

// main iscsi status information 
// building the path and prepare info for enable iscsi 
func GetStats(iqn_path string) (*Stats, error) {
    var istats Stats

    // log.Debugf("lio: GetStats path :%s", iqn_path)
    istats.Name = filepath.Base(iqn_path)
    // log.Debugf("lio: GetStats name :%s", istats.Name)

    matches, err := filepath.Glob(filepath.Join(iqn_path, "tpgt*"))
    if err != nil {
        return nil, fmt.Errorf("lio: get TPGT error %v\n", err )
    }
    istats.Tpgt = make([]TPGT, len(matches))

    for ii, tpgt_path := range matches {
        // log.Debugf("lio: GetStats tpgt_path :%s", tpgt_path)
        istats.Tpgt[ii].Name = filepath.Base(tpgt_path)
        istats.Tpgt[ii].Tpgt_path = tpgt_path
        istats.Tpgt[ii].Is_enable, _ = isPathEnable(tpgt_path)
        if (istats.Tpgt[ii].Is_enable) {
                matches_luns_path, _ := getLun(tpgt_path)
                istats.Tpgt[ii].Luns = make([]LUN, len(matches_luns_path))

                for ll, lun_path := range matches_luns_path {
                    backstore, object_name, type_number, err := getLunLinkTarget(lun_path)
                    if err != nil {
                        // log.Errorf("lio: get TPGT Lun error %v\n", err )
                        continue
                    }
                    istats.Tpgt[ii].Luns[ll].Name       = filepath.Base(lun_path)
                    istats.Tpgt[ii].Luns[ll].Lun_path   = lun_path
                    istats.Tpgt[ii].Luns[ll].Backstore  = backstore
                    istats.Tpgt[ii].Luns[ll].Object_name= object_name
                    istats.Tpgt[ii].Luns[ll].Type_number= type_number
                }
        }
    }
    return &istats, nil
}

// utility function 
// check if the file "enable" contain enable message
func isPathEnable(path string) (bool, error) {
    var isEnable bool
    isEnable = false
    tmp, err := ioutil.ReadFile(filepath.Join(path, "enable"))
    // log.Debugf("lio: is Path %s enable?", tmp)
    if err != nil {
        return false, fmt.Errorf("is Path Enable error %v\n", err)
    }
    tmp_num, err := strconv.Atoi(strings.TrimSpace(string(tmp)))
    // log.Debugf("lio: isPathEnable number %d", tmp_num)
    if (tmp_num > 0) {
        isEnable = true
    }
    return isEnable, nil
}

func getLun(tpgt_path string) (matches []string, err error) {
    // log.Debugf("lio: getLun path %s", tpgt_path)
    matches, err = filepath.Glob(filepath.Join(tpgt_path, "lun/lun*"))
    if err != nil {
        return nil, fmt.Errorf("getLun error  %v\n", err )
    }
    return matches, nil
}

func getLunLinkTarget(lun_path string) (backstore_type string,
    object_name string, type_number string, err error) {
    files, err := ioutil.ReadDir(lun_path)
    if err != nil {
        return "", "", "", fmt.Errorf("lio getLunLinkTarget error  %v\n", err )
    }
    for _, file := range files {
        // log.Debugf("lio: lun dir list file ->%s<-",file.Name())
        // fileInfo, _:= os.Lstat(lun_path +  "/" + file.Name())
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
            //  log.Debugf("lio: object_name->%s-<, type->%s, type_number->%s<- ", 
            //  object_name, backstore_type, type_number)
            return backstore_type, object_name, type_number, nil
        }
    }
    return "", "", "", errors.New("lio getLunLinkTarget: Lun Link does not exist")
}

func ReadWriteOPS(iqn string, tpgt string, lun string) (readmb uint64,
    writemb uint64, iops uint64, err error){

    readmb_path := filepath.Join(TARGET_PATH, iqn, tpgt, "lun", lun,
        "statistics/scsi_tgt_port/read_mbytes")
    // log.Debugf("lio: Read File path: %s\n", readmb_path)

    if _, err := os.Stat(readmb_path); os.IsNotExist(err) {
        return 0, 0, 0, fmt.Errorf("lio: file %s is missing!", readmb_path)
    }
    readmb, err = readUintFromFile(readmb_path) 
    if err != nil {
        return 0, 0, 0, fmt.Errorf("lio: read_mbytes error %s!", err)
    }

    writemb_path := filepath.Join(TARGET_PATH, iqn, tpgt, "lun", lun,
        "statistics/scsi_tgt_port/write_mbytes")
    // log.Debugf("lio: Write File path: %s\n", readmb_path)

    if _, err := os.Stat(writemb_path); os.IsNotExist(err) {
        return 0, 0, 0, fmt.Errorf("lio: file %s is missing!", readmb_path)
    }
    writemb, err = readUintFromFile(writemb_path) 
    if err != nil {
        return 0, 0, 0, fmt.Errorf("lio: write_mbytes error %s!", err)
    }

    iops_path := filepath.Join(TARGET_PATH, iqn, tpgt, "lun", lun,
        "statistics/scsi_tgt_port/in_cmds")
    // log.Debugf("lio: Write File path: %s\n", iops_path)

    if _, err := os.Stat(iops_path); os.IsNotExist(err) {
        return 0, 0, 0, fmt.Errorf("lio: file %s is missing!", iops_path)
    }
    iops, err = readUintFromFile(iops_path) 
    if err != nil {
        return 0, 0, 0, fmt.Errorf("lio: in_cmds error %s!", err)
    }

    return readmb, writemb, iops, nil
}

func (fileio FILEIO) GetFileioUdev(fileio_number string, 
    object_name string) (fio *FILEIO, err error) {

    fileio.Name         = "fileio_" + fileio_number
    fileio.Fnumber      = fileio_number
    fileio.Object_name  = object_name
    
    // log.Debugf("lio: Fileio udev_path ->%s<-", filepath.Join(TARGET_CORE,
    // fileio.Name, fileio.Object_name, "udev_path"))
    udev_path := filepath.Join(TARGET_CORE, fileio.Name, fileio.Object_name, "udev_path")

    if _, err := os.Stat(udev_path); os.IsNotExist(err) {
        return nil, fmt.Errorf("lio: fileio_%s is missing file name ...!", fileio.Fnumber)
    }
    filename, err := ioutil.ReadFile(udev_path)
    if err != nil {
        return nil, fmt.Errorf("lio: Cannot read filename from udev link! :%s", udev_path)
    }
    fileio.Filename = strings.TrimSpace(string(filename))
    
    return &fileio, nil
}

func (iblock IBLOCK) GetIblockUdev(iblock_number string,
    object_name string) (ib *IBLOCK, err error) {

    iblock.Name         = "iblock_" + iblock_number
    iblock.Bnumber      = iblock_number
    iblock.Object_name  = object_name
    
    // log.Debugf("lio: IBlock udev_path ->%s<-", filepath.Join(TARGET_CORE,
    // iblock.Name, iblock.Object_name, "udev_path"))
    udev_path := filepath.Join(TARGET_CORE, iblock.Name, iblock.Object_name, "udev_path")

    if _, err := os.Stat(udev_path); os.IsNotExist(err) {
        return nil, fmt.Errorf("lio: iblock_%s is missing file name ...!",
        iblock.Bnumber)
    }
    filename, err := ioutil.ReadFile(udev_path)
    if err != nil {
        return nil, fmt.Errorf("lio: Cannot read iblock from udev link! :%s", udev_path)
    }
    iblock.Iblock = strings.TrimSpace(string(filename))
    
    return &iblock, nil
}

func (rbd RBD) GetRBDMatch(rbd_number string, pool_image string) (r *RBD, err error) {

    rbd.Name    = "rbd_" + rbd_number
    rbd.Rnumber = rbd_number
    // log.Debugf("lio: RBD info Name ->%s<-", rbd.Name)
    // log.Debugf("lio: RBD info Rnumber ->%s<-", rbd.Rnumber)
    // log.Debugf("lio: RBD info Pool-Image->%s<-", pool_image)

    system_rbds, err := filepath.Glob("/sys/devices/rbd/[0-9]*")
    if err != nil {
        return nil, fmt.Errorf("lio: Cannot find any rbd block!")
    }

    for system_rbd_number, system_rbd_path := range system_rbds { 
        var system_pool, system_image string = "", "" 
        // log.Debugf("lio: rbd_path ->%s<-", system_rbd_path)
        system_pool_path := filepath.Join(system_rbd_path, "pool")
        if _, err := os.Stat(system_pool_path); os.IsNotExist(err) {
            // log.Errorf("lio: rbd%d pool file %s is missing!",
            // system_rbd_number, system_pool_path )
            continue
        }
        b_system_pool, err := ioutil.ReadFile(system_pool_path)
        if err != nil {
            // log.Errorf("lio: Cannot read pool name from %s!", system_pool_path)
            continue
        } else { 
            system_pool = strings.TrimSpace(string(b_system_pool))
        }

        system_image_path := filepath.Join(system_rbd_path, "name")
        if _, err := os.Stat(system_image_path); os.IsNotExist(err) {
            // log.Errorf("lio: rbd%d image file %s is missing!",
            // system_rbd_number, system_image_path )
            continue
        }
        b_system_image, err := ioutil.ReadFile(system_image_path)
        if err != nil {
            // log.Errorf("lio: Cannot read image name from %s!", system_image_path)
            continue
        } else { 
            system_image = strings.TrimSpace(string(b_system_image))
        }
        // log.Debugf("lio: System rbd_%d :", system_rbd_number)
        // log.Debugf("lio: Matching label->rbd%s", rbd.Rnumber)
        // log.Debugf("lio: System pool --->%s<--- image --->%s<---", system_pool, system_image)
        // log.Debugf("lio: Matching pool-image->%s", pool_image)

        if matchRBD(fmt.Sprintf("%d", system_rbd_number), rbd.Rnumber) &&
        matchPoolImage(system_pool, system_image, pool_image) {
            rbd.Pool = system_pool
            rbd.Image= system_image
            return &rbd, nil 
        }
    }
    return nil, nil
}

func (rdmcp RDMCP) GetRDMCPPath(rdmcp_number string, object_name string) (r *RDMCP, err error) {
    rdmcp.Name          = "rd_mcp_" + rdmcp_number
    rdmcp.Object_name   = object_name

    rdmcp_path := filepath.Join(TARGET_CORE, rdmcp.Name, rdmcp.Object_name)
    // log.Debugf("lio: RDMCP path ->%s<-", rdmcp_path)

    if _, err := os.Stat(rdmcp_path); os.IsNotExist(err) {
        return nil, fmt.Errorf("lio: %s does not exist!", rdmcp_path)
    }
    isEnable, err := isPathEnable(rdmcp_path)
    if err != nil {
        return nil, fmt.Errorf("lio: error %v", err)
    }
    if isEnable { 
        return &rdmcp, nil
    }
    return nil, nil
}

func readUintFromFile(path string) (uint64, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return 0, err
    }
    value, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
    if err != nil {
        return 0, err
    }
    return value, nil
}

func matchRBD(rbd_number string, rbd_name string) (isEqual bool) {
    isEqual = false
    if strings.Compare(rbd_name, rbd_number) == 0 { 
        isEqual = true
    }
    return isEqual 
}

func matchPoolImage(pool string, image string, match_pool_image string) (isEqual bool) { 
    isEqual = false
    var pool_image = fmt.Sprintf("%s-%s", pool, image)
    // log.Debugf("lio: compare ->%s<- with ->%s<- ", pool_image, match_pool_image)
    if strings.Compare(pool_image, match_pool_image) == 0 { 
        isEqual = true
    }
    return isEqual 
}
