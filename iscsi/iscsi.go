// Copyright 2017 The Prometheus Authors
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
)

// TARGETPATH is static sys path for iscsi target 
const TARGETPATH = "/sys/kernel/config/target/iscsi"
// TARGETCORE static sys path for backstore 
const TARGETCORE = "/sys/kernel/config/target/core"

// TPGT struct for sys target portal group tag info
type TPGT struct {
    Name        string  // name of the tpgt group
    TpgtPath    string  // file path of tpgt 
    IsEnable    bool    // is the tpgt enable
    Luns        []LUN   // the Luns that tpgt has 
}

// LUN struct for sys logical unit number info
type LUN struct {
    Name        string  // name of the lun
    LunPath     string  // file path of the lun
    Backstore   string  // backstore of the lun
    ObjectName  string  // place holder for object 
    TypeNumber  string  // place holder for number of the device
}

// FILEIO struct for backstore info
type FILEIO struct { 
    Name        string  // name of the fileio 
    Fnumber     string  // number related to the backstore
    ObjectName  string  // place holder for object in iscsi object
    Filename    string  // link to the actual file being export
}

// IBLOCK struct for backstore info
type IBLOCK struct {
    Name        string  // name of the iblock
    Bnumber     string  // number related to the backstore 
    ObjectName  string  // place holder for object in iscsi object
    Iblock      string  // link to the actual block being export 
}

// RBD struct for backstore info
type RBD struct {
    Name    string  // name of the rbd 
    Rnumber string  // number related to the backstore 
    Pool    string  // place holder for the rbd pool
    Image   string  // place holder for the rbd image
}

// RDMCP struct for backstore info
type RDMCP struct {
    Name        string  // name of the rdm_cp 
    ObjectName  string  // place holder for object name 
}

// Stats struct for all targets info
type Stats struct {
    Name    string
    Tpgt    []TPGT
}

// GetStats is the main iscsi status information func
// building the path and prepare info for enable iscsi 
func GetStats(iqnPath string) (*Stats, error) {
    var istats Stats

    istats.Name = filepath.Base(iqnPath)
    matches, err := filepath.Glob(filepath.Join(iqnPath, "tpgt*"))
    if err != nil {
        return nil, fmt.Errorf("lio: get TPGT error %v", err )
    }
    istats.Tpgt = make([]TPGT, len(matches))

    for ii, tpgtPath := range matches {
        istats.Tpgt[ii].Name = filepath.Base(tpgtPath)
        istats.Tpgt[ii].TpgtPath = tpgtPath
        istats.Tpgt[ii].IsEnable, _ = isPathEnable(tpgtPath)
        if (istats.Tpgt[ii].IsEnable) {
                matchesLunsPath, _ := getLun(tpgtPath)
                istats.Tpgt[ii].Luns = make([]LUN, len(matchesLunsPath))

                for ll, lunPath := range matchesLunsPath {
                    backstore, objectName, typeNumber, err := getLunLinkTarget(lunPath)
                    if err != nil {
                        continue
                    }
                    istats.Tpgt[ii].Luns[ll].Name       = filepath.Base(lunPath)
                    istats.Tpgt[ii].Luns[ll].LunPath    = lunPath
                    istats.Tpgt[ii].Luns[ll].Backstore  = backstore
                    istats.Tpgt[ii].Luns[ll].ObjectName = objectName
                    istats.Tpgt[ii].Luns[ll].TypeNumber = typeNumber
                }
        }
    }
    return &istats, nil
}

// isPathEnable is a utility function 
// check if the file "enable" contain enable message
func isPathEnable(path string) (bool, error) {
    var isEnable bool
    isEnable = false
    tmp, err := ioutil.ReadFile(filepath.Join(path, "enable"))
    if err != nil {
        return false, fmt.Errorf("is Path Enable error %v", err)
    }
    tmpNum, err := strconv.Atoi(strings.TrimSpace(string(tmp)))
    if (tmpNum > 0) {
        isEnable = true
    }
    return isEnable, nil
}

func getLun(tpgtPath string) (matches []string, err error) {
    matches, err = filepath.Glob(filepath.Join(tpgtPath, "lun/lun*"))
    if err != nil {
        return nil, fmt.Errorf("getLun error  %v", err )
    }
    return matches, nil
}

func getLunLinkTarget(lunPath string) (backstoreType string,
    objectName string, typeNumber string, err error) {
    files, err := ioutil.ReadDir(lunPath)
    if err != nil {
        return "", "", "", fmt.Errorf("lio getLunLinkTarget error  %v", err )
    }
    for _, file := range files {
        fileInfo, _:= os.Lstat(lunPath +  "/" + file.Name())
        if fileInfo.Mode() & os.ModeSymlink != 0 {
            target, err := os.Readlink( lunPath +  "/" + fileInfo.Name())
            if err != nil {
                return "", "", "", fmt.Errorf("Readlink err %v", err)
            }
            p1, objectName := filepath.Split(target)
            _, typeWithNumber := filepath.Split(filepath.Clean(p1))

            tmp := strings.Split(typeWithNumber, "_")
            backstoreType, typeNumber := tmp[0], tmp[1]
            if len(tmp) == 3 {
                backstoreType = fmt.Sprintf("%s_%s", tmp[0], tmp[1])
                typeNumber = tmp[2]
            }
            return backstoreType, objectName, typeNumber, nil
        }
    }
    return "", "", "", errors.New("lio getLunLinkTarget: Lun Link does not exist")
}

// ReadWriteOPS read and return the stat of read and write in megabytes, 
// and total commands that send to the target
func ReadWriteOPS(iqn string, tpgt string, lun string) (readmb uint64,
    writemb uint64, iops uint64, err error){

    readmbPath := filepath.Join(TARGETPATH, iqn, tpgt, "lun", lun,
        "statistics/scsi_tgt_port/read_mbytes")

    if _, err := os.Stat(readmbPath); os.IsNotExist(err) {
        return 0, 0, 0, fmt.Errorf("lio: file %s is missing", readmbPath)
    }
    readmb, err = readUintFromFile(readmbPath) 
    if err != nil {
        return 0, 0, 0, fmt.Errorf("lio: read_mbytes error %s", err)
    }

    writembPath := filepath.Join(TARGETPATH, iqn, tpgt, "lun", lun,
        "statistics/scsi_tgt_port/write_mbytes")

    if _, err := os.Stat(writembPath); os.IsNotExist(err) {
        return 0, 0, 0, fmt.Errorf("lio: file %s is missing", readmbPath)
    }
    writemb, err = readUintFromFile(writembPath) 
    if err != nil {
        return 0, 0, 0, fmt.Errorf("lio: write_mbytes error %s", err)
    }

    iopsPath := filepath.Join(TARGETPATH, iqn, tpgt, "lun", lun,
        "statistics/scsi_tgt_port/in_cmds")

    if _, err := os.Stat(iopsPath); os.IsNotExist(err) {
        return 0, 0, 0, fmt.Errorf("lio: file %s is missing", iopsPath)
    }
    iops, err = readUintFromFile(iopsPath) 
    if err != nil {
        return 0, 0, 0, fmt.Errorf("lio: in_cmds error %s", err)
    }

    return readmb, writemb, iops, nil
}

// GetFileioUdev is getting the actual info to build up 
// the FILEIO data and match with the enable target 
func (fileio FILEIO) GetFileioUdev(fileioNumber string, 
    objectName string) (fio *FILEIO, err error) {

    fileio.Name         = "fileio_" + fileioNumber
    fileio.Fnumber      = fileioNumber
    fileio.Object_name  = objectName
    
    udevPath := filepath.Join(TARGETCORE, fileio.Name, fileio.Object_name, "udev_path")

    if _, err := os.Stat(udevPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("lio: fileio_%s is missing file name", fileio.Fnumber)
    }
    filename, err := ioutil.ReadFile(udevPath)
    if err != nil {
        return nil, fmt.Errorf("lio: Cannot read filename from udev link :%s", udevPath)
    }
    fileio.Filename = strings.TrimSpace(string(filename))
    
    return &fileio, nil
}

// GetIblockUdev is getting the actual info to build up 
// the IBLOCK data and match with the enable target 
func (iblock IBLOCK) GetIblockUdev(iblockNumber string,
    objectName string) (ib *IBLOCK, err error) {

    iblock.Name         = "iblock_" + iblockNumber
    iblock.Bnumber      = iblockNumber
    iblock.ObjectName   = objectName
    
    udevPath := filepath.Join(TARGETCORE, iblock.Name, iblock.ObjectName, "udev_path")

    if _, err := os.Stat(udevPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("lio: iblock_%s is missing file name", iblock.Bnumber)
    }
    filename, err := ioutil.ReadFile(udevPath)
    if err != nil {
        return nil, fmt.Errorf("lio: Cannot read iblock from udev link :%s", udevPath)
    }
    iblock.Iblock = strings.TrimSpace(string(filename))
    
    return &iblock, nil
}

// GetRBDMatch is getting the actual info to build up 
// the RBD data and match with the enable target 
func (rbd RBD) GetRBDMatch(rbdNumber string, poolImage string) (r *RBD, err error) {

    rbd.Name    = "rbd_" + rbdNumber
    rbd.Rnumber = rbdNumber

    systemRbds, err := filepath.Glob("/sys/devices/rbd/[0-9]*")
    if err != nil {
        return nil, fmt.Errorf("lio: Cannot find any rbd block")
    }

    for systemRbdNumber, systemRbdPath := range systemRbds { 
        var systemPool, systemImage string = "", "" 
        systemPoolPath := filepath.Join(systemRbdPath, "pool")
        if _, err := os.Stat(systemPoolPath); os.IsNotExist(err) {
            continue
        }
        bSystemPool, err := ioutil.ReadFile(systemPoolPath)
        if err != nil {
            continue
        } else { 
            system_pool = strings.TrimSpace(string(bSystemPool))
        }

        systemImagePath := filepath.Join(systemRbdPath, "name")
        if _, err := os.Stat(systemImagePath); os.IsNotExist(err) {
            continue
        }
        bSystemImage, err := ioutil.ReadFile(systemImagePath)
        if err != nil {
            continue
        } else { 
            system_image = strings.TrimSpace(string(bSystemImage))
        }

        if matchRBD(fmt.Sprintf("%d", systemRbdNumber), rbd.Rnumber) &&
        matchPoolImage(systemPool, systemImage, poolImage) {
            rbd.Pool = systemPool
            rbd.Image= systemImage
            return &rbd, nil 
        }
    }
    return nil, nil
}

// GetRDMCPPath is getting the actual info to build up RDMCP data 
func (rdmcp RDMCP) GetRDMCPPath(rdmcpNumber string, objectName string) (r *RDMCP, err error) {
    rdmcp.Name          = "rd_mcp_" + rdmcpNumber
    rdmcp.ObjectName   = objectName

    rdmcpPath := filepath.Join(TARGETCORE, rdmcp.Name, rdmcp.ObjectName)

    if _, err := os.Stat(rdmcpPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("lio: %s does not exist", rdmcpPath)
    }
    isEnable, err := isPathEnable(rdmcpPath)
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

func matchRBD(rbdNumber string, rbdName string) (isEqual bool) {
    isEqual = false
    if strings.Compare(rbdName, rbdNumber) == 0 { 
        isEqual = true
    }
    return isEqual 
}

func matchPoolImage(pool string, image string, matchPoolImage string) (isEqual bool) { 
    isEqual = false
    var poolImage = fmt.Sprintf("%s-%s", pool, image)
    if strings.Compare(poolImage, matchPoolImage) == 0 { 
        isEqual = true
    }
    return isEqual 
}
