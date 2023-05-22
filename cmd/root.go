package cmd

import (
	"fmt"
	"os"

	"github.com/levinion/fr/finder"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "fr",
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {

		finder.Cfg.JumpHiddenItem, _ = cmd.Flags().GetBool("all")
		finder.Cfg.JumpHiddenItem = !finder.Cfg.JumpHiddenItem
		finder.Cfg.Dir, _ = cmd.Flags().GetBool("dir")
		finder.Cfg.File, _ = cmd.Flags().GetBool("file")
		finder.Cfg.MatchFullPath, _ = cmd.Flags().GetBool("full")
		finder.Cfg.TimeRule, _ = cmd.Flags().GetString("time")
		finder.Cfg.Blur, _ = cmd.Flags().GetBool("blur")
		finder.Cfg.Fast, _ = cmd.Flags().GetBool("rush")
		excludeDirSlice, _ := cmd.Flags().GetStringSlice("exclude")
		if excludeDirSlice != nil {
			finder.Cfg.ExcludeDir = excludeDirSlice
		}
		finder.Cfg.Go, _ = cmd.Flags().GetBool("go")
		finder.Cfg.Num, _ = cmd.Flags().GetInt("num")

		//扩展名和类型相关
		ext, _ := cmd.Flags().GetString("ext")
		t, _ := cmd.Flags().GetString("type")
		if ext != "" && t != "" {
			fmt.Println("Error: Choose either file type or extension, not both.")
			return
		}
		if ext != "" {
			finder.Cfg.Ext = []string{ext}
		}
		if t != "" {
			e, ok := finder.TypeMap[t]
			if ok {
				finder.Cfg.Ext = e
			} else {
				fmt.Println("Error: invaild type")
				return
			}
		}

		//Debug用
		if os.Getenv("FR_DEBUG") == "true" {
			fmt.Println(finder.Cfg)
		}

		root, _ := cmd.Flags().GetString("root")

		length := len(args)
		switch length {
		case 0:
			if root == "" {
				root = "."
			}
			err := finder.Find(".*", root)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		case 1:
			if root == "" {
				root = "."
			}
			err := finder.Find(args[0], root)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		case 2:
			if root == "" {
				root = args[1]
			}
			err := finder.Find(args[0], root)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.Flags().BoolP("all", "a", false, "show all items along with hidden items")
	rootCmd.Flags().StringP("ext", "e", "", "extented name")
	rootCmd.Flags().BoolP("file", "f", false, "file type")
	rootCmd.Flags().BoolP("dir", "d", false, "dir type")
	rootCmd.Flags().BoolP("full", "F", false, "match full path")
	rootCmd.Flags().StringP("root", "R", "", "root dir that will be searched in")
	rootCmd.Flags().StringP("type", "t", "", "specific file type, can be video, audio, img and so on")
	rootCmd.Flags().StringP("time", "T", "", "modified time,accessed time and created time")
	rootCmd.Flags().StringSliceP("exclude", "E", nil, "exclude dir that will not be search")
	rootCmd.Flags().BoolP("blur", "b", false, "blur match")
	rootCmd.Flags().BoolP("rush", "r", false, "fast but not sorted")
	rootCmd.Flags().BoolP("go", "g", false, "")
	rootCmd.Flags().IntP("num", "n", 1, "")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
