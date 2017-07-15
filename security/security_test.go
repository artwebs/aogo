package security

import (
	"fmt"
	// "github.com/artwebs/aogo/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecurity(t *testing.T) {
	Convey("DES加密", t, func() {
		key := "Y8gyxetKJ68N3d35Lass72GP"
		data := "a12*&1c中文"
		// key := "Y8gyxetK"
		// fmt.Println(string(util.RandomBytes(24)))
		// desObj := NewSecurityDES()
		// crypted1, _ := desObj.EncryptString(key, data)
		// fmt.Println(crypted1)
		// fmt.Println(desObj.DecryptString(key, crypted1))
		// So("1", ShouldEqual, "1")

		des3Obj := NewSecurityTripleDES()
		crypted2, _ := des3Obj.EncryptString(key, data)
		fmt.Println(crypted2)
		fmt.Println(des3Obj.DecryptString(key, crypted2))

		aesObj := NewSecurityAES()
		crypted3, _ := aesObj.EncryptString(key, data)
		fmt.Println(crypted3)
		fmt.Println(aesObj.DecryptString(key, crypted3))
	})
}
