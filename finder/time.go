package finder

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

func timePipe(info fs.FileInfo) (bool, error) {
	//TimeRule判断
	if Cfg.TimeRule != "" {
		// 确认表达式合法性，之后解析每部分的值
		if ok, _ := regexp.MatchString(`^([mac])?[><][\d]+[smhdMy]?$`, Cfg.TimeRule); ok {
			ra, _ := regexp.Compile(`^[mac]`)
			sa := ra.FindString(Cfg.TimeRule)
			rb, _ := regexp.Compile(`[><]`)
			sb := rb.FindString(Cfg.TimeRule)
			rc, _ := regexp.Compile(`[\d]+`)
			sc := rc.FindString(Cfg.TimeRule)
			rd, _ := regexp.Compile(`[smhdMy]$`)
			sd := rd.FindString(Cfg.TimeRule)

			//定义变量
			nowTime := time.Now()
			unixTime := nowTime.Unix()
			var controllerTime int64
			num, err := strconv.Atoi(sc)
			if err != nil {
				return false, err
			}
			switch sd {
			case "s":
				controllerTime = int64(num)
			case "m":
				controllerTime = int64(num * 60)
			case "h":
				controllerTime = int64(num * 60 * 60)
			case "d","":
				controllerTime = int64(num * 60 * 60 * 24)
			case "M":
				controllerTime = int64(num * 60 * 60 * 24 * 30)
			case "y":
				controllerTime = int64(num * 60 * 60 * 24 * 30 * 12)
			}
			ft := unixTime - controllerTime
			ct := time.Unix(ft, 0)

			switch sa {
			case "m","":
				Cfg.timeRuleType = "m"
				modTime := info.ModTime()
				switch sb {
				case "<":
					if modTime.Before(ct) {
						return false, nil
					}
				case ">":
					if modTime.After(ct) {
						return false, nil
					}
				}
			case "a":
				Cfg.timeRuleType = "a"
				accessedTime := time.Unix(info.Sys().(*syscall.Stat_t).Atim.Sec, 0)
				switch sb {
				case "<":
					if accessedTime.Before(ct) {
						return false, nil
					}
				case ">":
					if accessedTime.After(ct) {
						return false, nil
					}
				}
			case "c":
				Cfg.timeRuleType = "c"
				createdTime := time.Unix(info.Sys().(*syscall.Stat_t).Ctim.Sec, 0)
				switch sb {
				case "<":
					if createdTime.Before(ct) {
						return false, nil
					}
				case ">":
					if createdTime.After(ct) {
						return false, nil
					}
				}
			}
		} else {
			fmt.Println("Error: invalid expression")
			os.Exit(1)
		}
	}
	return true, nil
}
