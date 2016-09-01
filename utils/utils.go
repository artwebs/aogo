package utils

import (
	"crypto/md5"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"reflect"
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
	name := file
	if i := strings.LastIndex(file, DirSep()); i > -1 {
		name = name[i+1:]
	}
	return name
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

func FileCopy(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		log.Println(err)
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
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

//str -24h(s,m,h)
func DateTimeAdd(str string, len time.Duration) time.Time {
	tmp, _ := time.ParseDuration(str)
	return time.Now().Add(tmp * len)
}

func DataTimeForString(str string) time.Time {
	tmp, _ := time.Parse("2006-01-02 15:04:05", str)
	return tmp
}

func DateTime(t time.Time) string {
	return DataTimeForamter(t, "yy-mm-dd hh:mi:ss")
}

func DataTimeForamter(t time.Time, f string) string {
	f = strings.Replace(f, "yy", "2006", -1)
	f = strings.Replace(f, "mm", "01", -1)
	f = strings.Replace(f, "dd", "02", -1)
	f = strings.Replace(f, "hh", "15", -1)
	f = strings.Replace(f, "mi", "04", -1)
	f = strings.Replace(f, "ss", "05", -1)
	return t.Format(f)
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

func RandomBytes(size int) []byte {
	ikind, kinds, result := 0, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		ikind = rand.Intn(3)
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func MachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(crypto_rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

func Hex(b []byte) string {
	return hex.EncodeToString(b)
}

func Tag(obj interface{}) string {
	return reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
}

func StrUpperUnderline(s string) string {
	reg := regexp.MustCompile(`([A-Z])`)
	return strings.TrimPrefix(strings.ToLower(reg.ReplaceAllString(s, "_"+"$1")), "_")
}

func UrlEncode(s string) string {
	return url.QueryEscape(s)
}

func UrlDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

func HttpGet(url string) string {
	response, _ := http.Get(url)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}

func HttpClientIP(r *http.Request) (string, string, error) {
	var ip string
	var port string
	var err error
	ip, port, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ip, port, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)

	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return ip, port, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}
	return ip, port, nil
}
