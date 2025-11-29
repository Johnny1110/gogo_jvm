package method_area

// SymRef 符號引用基類
// 所有符號引用（ClassRef, FieldRef, MethodRef）的共同基礎
type SymRef struct {
	cp        *RuntimeConstantPool // 所屬常量池
	className string               // 類名
	class     *Class               // 解析後的類引用（緩存）
}

// ResolvedClass 解析類引用
// 懶加載：第一次調用時解析，之後返回緩存
func (r *SymRef) ResolvedClass() *Class {
	if r.class == nil { // 第一次：r.class 是 nil
		r.resolveClassRef() // 執行加載（慢）
	}
	return r.class // 之後：直接返回緩存（快）
}

func (r *SymRef) resolveClassRef() {
	// 獲取當前類的 ClassLoader
	d := r.cp.Class()
	// 用 ClassLoader 加載目標類
	c := d.Loader().LoadClass(r.className)
	// TODO: 訪問權限檢查
	r.class = c
}

// ClassName 獲取類名
func (r *SymRef) ClassName() string {
	return r.className
}
