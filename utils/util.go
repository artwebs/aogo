package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	DARWIN = iota
	LINUX
	OTHERSYSTEM
)

func FileExt(file string) string {
	return path.Ext(file)
}

func FileName(file string) string {
	name := strings.TrimSuffix(file, FileExt(file))
	if i := strings.LastIndex(file, DirSep()); i > -1 {
		name = name[i+1:]
	}
	return name
}

func FileBaseName(file string) string {
	return path.Base(file)
}

func FileBaseDir(file string) string {
	name := file
	if i := strings.LastIndex(file, DirSep()); i > -1 {
		name = name[:i+1]
	}
	return name
}

func FileIsExist(file string) bool {
	flag := false
	if _, err := os.Stat(file); err == nil {
		flag = true
	}
	return flag
}

func FileRemove(file string) error {
	return os.RemoveAll(file)
}

func StringSearch(s, ex string) (group []string) {
	reg, _ := regexp.Compile(ex)
	matched := reg.MatchString(string(s))
	if matched {
		group = reg.FindStringSubmatch(s)
	}
	return group
}

func NowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NowDateTimeFormater(fmt string) string {
	return TimeForamter(NowDateTime(), fmt)
}

func TimeForamter(s, fmt string) string {
	fmt = strings.Replace(fmt, "yy", "2006", -1)
	fmt = strings.Replace(fmt, "mm", "01", -1)
	fmt = strings.Replace(fmt, "dd", "02", -1)
	fmt = strings.Replace(fmt, "hh", "15", -1)
	fmt = strings.Replace(fmt, "mi", "04", -1)
	fmt = strings.Replace(fmt, "ss", "05", -1)
	timeformatdate, _ := time.Parse("2006-01-02 15:04:05", s)
	return timeformatdate.Format(fmt)
}

func TimeToInt(s string) int {
	rs := 0
	arr := strings.Split(s, ":")
	sep := 60 * 60
	for _, v := range arr {
		temp, _ := strconv.ParseInt(v, 10, 64)
		rs += int(temp) * sep
		sep = sep / 60
	}
	return rs
}

func TimeForInt(l int) string {
	sep := 60 * 60
	s := ""
	for i := 0; i < 3; i++ {
		if s != "" {
			s += ":"
		}
		temp := l / sep
		s += fmt.Sprintf("%02d", temp)
		l = l - temp*sep
		sep = sep / 60
	}
	return s
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

func DirSep() (s string) {
	switch os := runtime.GOOS; os {
	case "darwin":
		s = "/"
	case "linux":
		s = "/"
	default:
		s = "\\"
	}
	return s
}

func System() (i int) {
	switch os := runtime.GOOS; os {
	case "darwin":
		i = DARWIN
	case "linux":
		i = LINUX
	default:
		i = OTHERSYSTEM
	}
	return i
}

func ExecCMD(cmdstr string) (string, error) {
	fields := strings.Fields(cmdstr)
	cmd := exec.Command(fields[0], fields[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("exec", cmdstr, err)
	}
	return strings.TrimSpace(string(output)), err
}
