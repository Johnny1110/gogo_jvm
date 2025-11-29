package method_area

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseMethodDescriptor(t *testing.T) {
	desc1 := "([[Ljava/lang/Object;)V" // -> void method(Object[][])
	md1 := parseMethodDescriptor(desc1)
	fmt.Println("MethodDescriptor-1: ", md1)
	assert.Equal(t, "[[Ljava/lang/Object;", md1.parameterTypes[0])
	assert.Equal(t, "V", md1.returnType)

	desc2 := "()V"
	md2 := parseMethodDescriptor(desc2)
	fmt.Println("MethodDescriptor-2: ", md2)
	assert.Equal(t, 0, len(md2.parameterTypes))
	assert.Equal(t, "V", md2.returnType)

	desc3 := "([I)V"
	md3 := parseMethodDescriptor(desc3)
	fmt.Println("MethodDescriptor-3: ", md3)
	assert.Equal(t, "[I", md3.parameterTypes[0])
	assert.Equal(t, "V", md3.returnType)
}
