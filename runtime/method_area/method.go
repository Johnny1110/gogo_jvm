package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/common"
)

// Method 運行時方法結構
//
// 方法是執行的核心單位，包含：
// - 方法簽名（名稱 + 描述符）
// - 訪問標誌
// - 字節碼
// - 操作數棧大小
// - 局部變量表大小
type Method struct {
	accessFlags  uint16
	name         string
	descriptor   string
	class        *Class // 所屬類
	maxStack     uint16 // 操作數棧最大深度
	maxLocals    uint16 // 局部變量表大小
	code         []byte // 方法字節碼
	argSlotCount uint   // 參數佔用的 rtcore 數量
}

// newMethods 從 ClassFile 的 MemberInfo 創建 Method 列表
func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class
	method.accessFlags = cfMethod.AccessFlags()
	method.name = cfMethod.Name()
	method.descriptor = cfMethod.Descriptor()
	method.copyAttributes(cfMethod)
	method.calcArgSlotCount()
	return method
}

// copyAttributes 複製 Code 屬性
func (m *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil {
		m.maxStack = codeAttr.MaxStack()
		m.maxLocals = codeAttr.MaxLocals()
		m.code = codeAttr.Code()
	}
}

// calcArgSlotCount 計算參數佔用的 rtcore 數量
// 根據方法描述符解析參數類型
// 例如：(II)V → 2 個 int → 2 slots
//
//	(JD)V → 1 個 long + 1 個 double → 4 slots
//	(Ljava/lang/String;I)V → 1 個引用 + 1 個 int → 2 slots
func (m *Method) calcArgSlotCount() {
	parsedDescriptor := parseMethodDescriptor(m.descriptor)
	for _, paramType := range parsedDescriptor.parameterTypes {
		m.argSlotCount++
		// long 和 double 佔 2 個 rtcore
		if paramType == "J" || paramType == "D" {
			m.argSlotCount++
		}
	}
	// 非靜態方法需要額外一個 rtcore 存 this
	if !m.IsStatic() {
		m.argSlotCount++
	}
}

// =============== Getters ===============

func (m *Method) Name() string        { return m.name }
func (m *Method) Descriptor() string  { return m.descriptor }
func (m *Method) Class() *Class       { return m.class }
func (m *Method) MaxStack() uint16    { return m.maxStack }
func (m *Method) MaxLocals() uint16   { return m.maxLocals }
func (m *Method) Code() []byte        { return m.code }
func (m *Method) ArgSlotCount() uint  { return m.argSlotCount }
func (m *Method) AccessFlags() uint16 { return m.accessFlags }

// =============== Access Flags ===============

func (m *Method) IsPublic() bool       { return m.accessFlags&common.ACC_PUBLIC != 0 }
func (m *Method) IsPrivate() bool      { return m.accessFlags&common.ACC_PRIVATE != 0 }
func (m *Method) IsProtected() bool    { return m.accessFlags&common.ACC_PROTECTED != 0 }
func (m *Method) IsStatic() bool       { return m.accessFlags&common.ACC_STATIC != 0 }
func (m *Method) IsFinal() bool        { return m.accessFlags&common.ACC_FINAL != 0 }
func (m *Method) IsSynchronized() bool { return m.accessFlags&common.ACC_SYNCHRONIZED != 0 }
func (m *Method) IsBridge() bool       { return m.accessFlags&common.ACC_BRIDGE != 0 }
func (m *Method) IsVarargs() bool      { return m.accessFlags&common.ACC_VARARGS != 0 }
func (m *Method) IsNative() bool       { return m.accessFlags&common.ACC_NATIVE != 0 }
func (m *Method) IsAbstract() bool     { return m.accessFlags&common.ACC_ABSTRACT != 0 }
func (m *Method) IsStrict() bool       { return m.accessFlags&common.ACC_STRICT != 0 }

// =============== Helper ===============

// IsPSVM 是否是 public static void main(String[] args)
func (m *Method) IsPSVM() bool {
	return m.IsPublic() && m.IsStatic() &&
		m.name == "main" && m.descriptor == "([Ljava/lang/String;)V"
}
