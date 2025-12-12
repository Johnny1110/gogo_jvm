package native

// define init() func for native package.
// import native package in main.go, and this func will trigger all sub package like java/io, java/lang ...
import (
	// using _ import, only for auto trigger init() func
	_ "github.com/Johnny1110/gogo_jvm/native/java/io" // this will auto trigger all .go file init() func
	_ "github.com/Johnny1110/gogo_jvm/native/java/lang"
	// _ "github.com/Johnny1110/gogo_jvm/native/java/util"
)
