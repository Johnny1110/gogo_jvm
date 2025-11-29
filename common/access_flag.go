package common

const (
	ACC_PUBLIC       = 0x0001 // public
	ACC_PRIVATE      = 0x0002 // private
	ACC_PROTECTED    = 0x0004 // protected
	ACC_STATIC       = 0x0008 // static
	ACC_FINAL        = 0x0010 // final
	ACC_SYNCHRONIZED = 0x0020 // synchronized (方法)
	ACC_VOLATILE     = 0x0040 // volatile (字段)
	ACC_TRANSIENT    = 0x0080 // transient (字段)
	ACC_NATIVE       = 0x0100 // native (方法)
	INTERFACE        = 0x0200
	ACC_ABSTRACT     = 0x0400 // abstract
	ACC_SYNTHETIC    = 0x1000
	ACC_ENUM         = 0x4000
	ACC_SUPER        = 0x0020
	ACC_INTERFACE    = 0x0200
	ACC_ANNOTATION   = 0x2000
	ACC_STRICT       = 0x0800
	ACC_VARARGS      = 0x0080
	ACC_BRIDGE       = 0x0040
)
