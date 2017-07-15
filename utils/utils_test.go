package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	// "log"
	// "os"
	"testing"
)

func TestUtil(t *testing.T) {
	filename := "/data/temp/abc.txt"
	Convey("获取文件扩张名", t, func() {
		So(FileExt(filename), ShouldEqual, ".txt")
	})
	Convey("获取文件，不包括扩张名", t, func() {
		So(FileName(filename), ShouldEqual, "abc")
	})
	Convey("获取全名文件", t, func() {
		So(FileBaseName(filename), ShouldEqual, "abc.txt")
	})
	Convey("获取目录名称", t, func() {
		So(FileBaseDir(filename), ShouldEqual, "/data/temp/")
	})

	Convey("文件及文件夹是否存在", t, func() {
		So(FileIsExist("util.go"), ShouldEqual, true)
		So(FileIsExist("util11.go"), ShouldEqual, false)
		So(FileIsExist("./test"), ShouldEqual, false)
	})

	Convey("删除文件及文件夹", t, func() {
		So(FileRemove("xxx.x"), ShouldEqual, nil)
	})

	Convey("字符串搜索", t, func() {
		So(StringSearch("file '/data/temp/6.flv'", "'(.+)'")[1], ShouldEqual, "/data/temp/6.flv")
	})

	Convey("时间格式化", t, func() {
		So(TimeForamter("2015-08-12 00:00:00", "yymmddhhmiss"), ShouldEqual, "20150812000000")
	})

	Convey("时间转换为int", t, func() {
		So(TimeToInt("00:02:01"), ShouldEqual, 121)
	})

	Convey("int转换为时间", t, func() {
		So(TimeForInt(121), ShouldEqual, "00:02:01")
	})

	Convey("是否含有中文", t, func() {
		So(IsChineseChar("是否含有中文"), ShouldEqual, true)
	})

	Convey("文件夹分隔", t, func() {
		So(DirSep(), ShouldEqual, "/")
	})

	Convey("获取系统类型", t, func() {
		So(System(), ShouldEqual, DARWIN)
	})

	// log.Println("-->", os.Environ()["GOCHAR"])
	Convey("运行命令", t, func() {
		value, _ := ExecCMD("echo hello world")
		So(value, ShouldEqual, "hello world")
	})

}
