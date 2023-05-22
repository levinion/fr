package finder

var TypeMap map[string][]string = map[string][]string{
	"video": {"mp4", "mkv", "avi", "wmv", "mov", "rmvb"},
	"audio": {"mp3", "wav", "flac", "aac", "wma"},
	"img":   {"jpg", "png", "gif", "bmp", "tiff"},
	"doc":   {"md", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "txt"},
	"exec":  {"exe", "msi", "bin"},
	"comp":  {"zip", "rar", "7z", "tar", "gz", "bz2"},
	"code": {"c", "cpp", "java", "go", "js", "py", "rb",
		"swift", "ts", "sh", "vue", "dart", "php", "cs", "asm", "lua"},
	"font":  {"ttf", "otf", "woff", "woff2"},
	"other": {"iso", "dmg", "dll"},
}
