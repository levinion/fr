package finder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

type config struct {
	JumpHiddenItem bool
	Ext            []string
	Dir            bool
	File           bool
	MatchFullPath  bool
	TimeRule       string
	timeRuleType   string
	ExcludeDir     []string
	Blur           bool
	Fast           bool
	Go             bool
	Num            int
}

var Cfg *config = &config{
	JumpHiddenItem: true,
	Ext:            nil,
	Dir:            false,
	File:           false,
	MatchFullPath:  false,
	TimeRule:       "",
	timeRuleType:   "",
	ExcludeDir:     nil,
	Blur:           false,
	Fast:           false,
	Go:             false,
	Num:            1,
}

type finderStruct struct {
	sync.Mutex
	items []*Item
}

type Item struct {
	item string
	info fs.FileInfo
}

var finderObj *finderStruct = &finderStruct{
	items: make([]*Item, 0),
}

func Find(pattern, root string) error {
	return find(pattern, root, func(foo, path string, info fs.FileInfo) (bool, error) {
		//附带File或Dir参数时进行文件类型判断
		if (Cfg.Dir || Cfg.File) && !(Cfg.Dir && Cfg.File) {
			if !Cfg.Dir { //若过滤文件夹
				if info.IsDir() {
					return false, nil
				}
			} else { //若过滤文件
				if !info.IsDir() {
					return false, nil
				}
			}
		}

		timePipeFlag, err := timePipe(info)
		if err != nil || !timePipeFlag {
			return false, err
		}

		//Ext与Type判断
		if Cfg.Ext != nil {
			if info.IsDir() {
				return false, nil
			}

			e := strings.TrimPrefix(filepath.Ext(path), ".")

			if itemNotInList(e, Cfg.Ext) {
				return false, nil
			}
		}

		//正则判断
		s := filepath.Base(path)
		if Cfg.MatchFullPath {
			s = path
		}
		if Cfg.Blur {
			pattern = strings.ToLower(pattern)
			s = strings.ToLower(s)
		}
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			return false, err
		}
		return matched, nil
	})
}

func find(foo, root string, f func(foo, path string, info fs.FileInfo) (bool, error)) error {
	return _find(foo, root, true, f)
}

var wg = sync.WaitGroup{}

func _find(foo, root string, main bool, f func(foo, path string, info fs.FileInfo) (bool, error)) error {

	dirEntry, err := os.ReadDir(root)
	if err != nil {
		//权限问题，跳过，直接返回
		if !main {
			wg.Done()
		}
		return nil
	}
	for _, entry := range dirEntry {
		info, _ := entry.Info()
		name := info.Name()
		path := filepath.Join(root, name)
		//忽略隐藏文件，排除项目
		if Cfg.JumpHiddenItem && strings.HasPrefix(name, ".") ||
			Cfg.ExcludeDir != nil && !itemNotInList(path, Cfg.ExcludeDir) {
			continue
		}

		//若是文件夹
		if info.IsDir() {
			wg.Add(1)
			go func(foo, path string, f func(foo string, path string, info fs.FileInfo) (bool, error)) {
				err := _find(foo, path, false, f)
				if err != nil {
					fmt.Println("Error: ", err)
					os.Exit(1)
				}
			}(foo, path, f)
		}
		//接口函数
		flag, err := f(foo, path, info)
		if err != nil {
			return err
		}
		if flag {
			if Cfg.Fast {
				displayItem(&Item{path, info})
			} else {
				finderObj.Lock()
				finderObj.items = append(finderObj.items, &Item{path, info})
				finderObj.Unlock()
			}
		}
	}
	if main { //若主线程则等待，待子进程全部返回后打印
		wg.Wait()
		if !Cfg.Fast {
			displayItems()
		}
	} else { //否则副线程完毕
		wg.Done()
	}
	return nil
}

func itemNotInList(ext string, list []string) bool {
	for i := range list {
		if ext == list[i] {
			return false
		}
	}
	return true
}

func displayItems() {
	sort.Slice(finderObj.items, func(i, j int) bool {
		return finderObj.items[i].item < finderObj.items[j].item
	})
	if !Cfg.Go {
		for _, item := range finderObj.items {
			displayItem(item)
		}
	} else {
		num := Cfg.Num - 1
		if num >= len(finderObj.items) || num < 0 {
			fmt.Printf("Error: invalid num %v\n", num)
			return
		}
		fmt.Println(finderObj.items[Cfg.Num-1].item)
	}
}

func displayItem(item *Item) {
	if Cfg.TimeRule != "" {

		switch Cfg.timeRuleType {
		case "m":
			t := item.info.ModTime()
			timeInfoString := t.Format(time.DateTime)
			fmt.Printf("%-40s\t\tModified Time: %s\n", item.item, timeInfoString)
		case "a":
			t := time.Unix(item.info.Sys().(*syscall.Stat_t).Atim.Sec, 0)
			timeInfoString := t.Format(time.DateTime)
			fmt.Printf("%-40s\t\tAccessed Time: %s\n", item.item, timeInfoString)
		case "c":
			t := time.Unix(item.info.Sys().(*syscall.Stat_t).Ctim.Sec, 0)
			timeInfoString := t.Format(time.DateTime)
			fmt.Printf("%-40s\t\tCreated Time: %s\n", item.item, timeInfoString)
		}
	} else {
		fmt.Println(item.item)
	}
}
