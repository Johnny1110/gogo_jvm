package method_area

// ============================================================
// v0.3.3: Object Lifecycle Support Methods
// ============================================================
// This file contains methods for:
// - Cloneable interface checking (for Object.clone())
// - Finalizer detection (for GC finalization support)

// ============================================================
// Cloneable Interface Check
// ============================================================

// IsCloneable checks if this class (or any of its superclasses) implements
// the java.lang.Cloneable interface.
//
// According to JLS:
// - All array types are implicitly Cloneable
// - For other types, the class hierarchy must be checked
//
// This is used by Object.clone() to determine if cloning is permitted.
func (c *Class) IsCloneable() bool {
	// Rule 1: All arrays are Cloneable (JLS requirement)
	if c.IsArray() {
		return true
	}

	// Rule 2: Check if this class or any superclass implements Cloneable
	// recursive check all []interfaces to find java/lang/Cloneable
	cloneableClass := c.Loader().LoadClass("java/lang/Cloneable", false)
	if cloneableClass == nil {
		panic("java class not found: java/lang/Cloneable")
	}
	return c.IsSubInterfaceOf(cloneableClass)
}

// ============================================================
// Finalizer Detection (for v0.5.x GC)
// ============================================================

// HasNonTrivialFinalizer checks if this class has a finalize() method
// that is NOT the default Object.finalize() (which does nothing).
//
// This is used to determine if objects of this class need special
// handling during garbage collection (finalization queue).
//
// Rules:
// - If the class doesn't have finalize()V method -> false
// - If the finalize() method is from java/lang/Object -> false (trivial)
// - If the finalize() method is overridden by this class or a superclass
func (c *Class) HasNonTrivialFinalizer() bool {
	// Look for finalize()V method in the class hierarchy
	method := c.GetMethod("finalize", "()V")
	if method == nil {
		return false
	}

	// Check if it's the trivial Object.finalize()
	declaringClass := method.Class()
	if declaringClass != nil && declaringClass.Name() == "java/lang/Object" {
		// Object.finalize() is considered trivial because it does nothing
		return false
	}

	// The class has a non-trivial finalize() method
	return true
}

// ============================================================
// Additional Helper Methods
// ============================================================

// IsSerializable checks if this class implements java.io.Serializable
// This is useful for deep copy via serialization
func (c *Class) IsSerializable() bool {
	// Arrays of primitives are Serializable
	// Arrays of Serializable elements are Serializable
	if c.IsArray() {
		componentClass := c.ComponentClass()
		if componentClass.IsPrimitive() {
			return true // primitive class
		}

		return componentClass.IsSerializable()
	}

	serializable := c.Loader().LoadClass("java/io/Serializable", false)
	return c.IsSubInterfaceOf(serializable)
}
