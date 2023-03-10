package global

import (
	"os"
	"path/filepath"
	"sync"
)

var once = new(sync.Once)

var RootDir string

// inferRootDir 推断出项目根目录
//在项目中任意一个目录执行编译然后运行或直接 go run，该函数都能正确找到项目的根目录。
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		//判断目录 d 下面是否存在 template 目录（只要是项目根目录下存在的目录即可，并非一定是 template）
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}
