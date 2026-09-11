package read

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"property/src/config"
	"property/src/types"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tealeg/xlsx/v3"
)

// chan for attendance read from file
var attReadChan = make(chan types.Attendance)

// chan for attendance file read finish
var attFinishChan = make(chan string)

// attsMap is a map that stores slices of Attendance structs, keyed by their file MD5 hash.
var attsMap = make(map[string][]types.Attendance)

// counter is an atomic counter that keeps track of the number of attendance files being processed concurrently.
var counter atomic.Int32

// morning and evening are time.Time values representing the cutoff times for morning and evening attendance checks.
var morning, _ = time.ParseInLocation(time.TimeOnly, "10:00:00", time.Local)
var evening, _ = time.ParseInLocation(time.TimeOnly, "15:00:00", time.Local)

// getAttFilepath returns a slice of file paths for all attendance files in the configured folder.
func getAttFilepath() []string {
	files, err := os.ReadDir(config.GetConfig().AttFolder)
	if err != nil {
		panic(err)
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && isAttfile(file.Name()) {
			paths = append(paths, filepath.Join(config.GetConfig().AttFolder, file.Name()))
		}
	}
	return paths
}

// isAttfile checks if the given file path is an attendance file based on its extension and the configured target year and month.
func isAttfile(filepath string) bool {

	return strings.HasSuffix(filepath, ".xlsx") &&
		strings.Contains(filepath, fmt.Sprintf("%d", config.GetConfig().TargetYear)) &&
		strings.Contains(filepath, fmt.Sprintf("%d", config.GetConfig().TargetMonth))

}

// ReadAtt reads all attendance files in the configured folder and returns a slice of Attendent structs.
func ReadAtt() *[]types.Attendance {
	paths := getAttFilepath()
	atts := make([]types.Attendance, 0)
	if len(paths) == 0 {
		log.Println("no att file found")
		return &atts
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(paths))
	counter.Add(int32(len(paths)))
	//read att file concurrently
	for _, path := range paths {
		go ReadSingleAtt(path, attReadChan)
	}

	//	handle att from file concurrently
	go func() {
		for {
			select {
			//handle att from file
			case att := <-attReadChan:
				if len(att.FileMd5) == 0 {
					continue
				}
				if atts, found := attsMap[att.FileMd5]; found {
					atts = append(atts, att)
					attsMap[att.FileMd5] = atts
				} else {
					attsMap[att.FileMd5] = []types.Attendance{att}
				}
				//signal that att read finish
			case s := <-attFinishChan:
				wg.Done()
				log.Printf("att file %s read finish", s)
				counter.Add(-1)
				if counter.Load() == 0 {
					log.Default().Println("all att file read finish")
					return
				}
			}
		}
	}()

	wg.Wait()
	combineAtts(attsMap, &atts)

	return &atts
}

// combineAtts combines attendance records from multiple files into a single slice of Attendent structs based on lookupkey.
func combineAtts(attMap map[string][]types.Attendance, all *[]types.Attendance) {
	if len(attMap) == 0 {
		return
	}
	//combine att if is empty or conbine is 0, else combine att with same key
	if config.GetConfig().Conbine == 0 || len(attMap) == 1 {
		for _, atts := range attMap {
			*all = append(*all, atts...)
		}
	} else if config.GetConfig().Conbine == 1 {
		for _, atts := range attMap {
			if len(*all) == 0 {
				*all = append(*all, atts...)
				continue
			}
			keyToAttMap := make(map[string]types.Attendance)
			for _, att := range atts {
				keyToAttMap[att.LookUpKey] = att
			}
			//combine att with same key
			for index, att := range *all {
				if att, found := keyToAttMap[att.LookUpKey]; found {

					for i := range len(att.Att) {
						(*all)[index].Att[i] = (*all)[index].Att[i] | att.Att[i]
					}
					(*all)[index] = att

					delete(keyToAttMap, att.LookUpKey)
				}
			}

			//new attendent not in all, add to all
			for _, att := range keyToAttMap {
				*all = append(*all, att)
			}
		}
	}
}

func ReadSingleAtt(path string, attReadchan chan types.Attendance) {
	// log.Printf("reading att file %s", path)
	excel, err := xlsx.OpenFile(path)

	if err != nil {
		panic(err)
	}
	sheet := excel.Sheet["刷卡记录"]

	maxRows := sheet.MaxRow
	md5, err := getFileMd5(path)
	for row := range maxRows {
		switch row {
		//title row
		case types.ATTENDENT_ROW_TITLE:
			val, err := sheet.Cell(row, 0)
			if err != nil {
				panic(err)
			}
			title := val.String()
			if title != "刷卡记录表" {
				panic("sheet err " + title)
			}
			//time row
		// case types.ATTENDENT_ROW_TIME:
		// for col := range maxCols {
		// 	val, err := sheet.Cell(row, col)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	println(row, ":", col, "=", val.String())
		// }
		case types.ATTENDENT_ROW_TIME, types.ATTENDENT_ROW_TABLE_TITLE_1, types.ATTENDENT_ROW_TABLE_TITLE_2:
			//ignore
		default:
			att := types.Attendance{Att: [31]int{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0}}
			readDataRow(sheet, row, &att)
			att.LookUpKey = fmt.Sprintf("%s-%s", att.Name, att.Dept)
			att.FileMd5 = md5

			attReadChan <- att
		}
	}
	attFinishChan <- fmt.Sprintf("%s:%s", path, md5)
}

// readDataRow reads a single row of attendance data from the given Excel sheet and populates the provided Attendent struct with the extracted information.
func readDataRow(sheet *xlsx.Sheet, row int, att *types.Attendance) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for col := range sheet.MaxCol {
		val, err := sheet.Cell(row, col)
		if err != nil {
			log.Println(err)
		}

		switch col {
		case types.ATTENDENT_COL_NO:
			no, err := val.Int()
			if err != nil {
				panic(err)
			}
			att.No = no
		case types.ATTENDENT_COL_NAME:
			att.Name = val.String()
		case types.ATTENDENT_COL_DEPT:
			att.Dept = val.String()
		default:
			ts := val.String()
			if len(ts) == 0 {
				continue
			}

			t := strings.Split(ts, "\n")

			for _, s := range t {
				if len(strings.Trim(s, " ")) == 0 {
					continue
				}
				checkTime, err := time.ParseInLocation(time.TimeOnly, s+":00", loc)
				if err != nil {
					log.Panic(err)
				}

				if checkTime.Before(morning) {
					att.Att[col-3] = att.Att[col-3] | types.BYTE_MORNING
				}

				if checkTime.After(evening) {
					att.Att[col-3] = att.Att[col-3] | types.BYTE_EVENING
				}
			}
		}
	}
}

func getFileMd5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	buf := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		hash.Write(buf[:n])
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
