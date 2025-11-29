package marea

// MethodDescriptor method descriptor parse result
// ex: (IDLjava/lang/String;)V
// parameterTypes: ["I", "D", "Ljava/lang/String;"]
// → returnType: "V"
// type encode:
// B - byte      C - char      D - double    F - float
// I - int       J - long      S - short     Z - boolean
// V - void      L; - ref type    [ - slice
type MethodDescriptor struct {
	parameterTypes []string
	returnType     string
}

func parseMethodDescriptor(descriptor string) MethodDescriptor {
	md := MethodDescriptor{}

	// skip '('
	i := 1

	// parse params type
	for descriptor[i] != ')' {
		paramType, nextIndex := parseFieldType(descriptor, i)
		md.parameterTypes = append(md.parameterTypes, paramType)
		i = nextIndex
	}

	// skip ')'
	i++

	// parse return type
	md.returnType, _ = parseFieldType(descriptor, i)

	return md
}

// parseFieldType parse next type, return paramType and nextIndex
func parseFieldType(descriptor string, startIndex int) (paramType string, nextIndex int) {
	switch descriptor[startIndex] {
	// basic type or void (1 char):
	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 'V':
		paramType = string(descriptor[startIndex])
		nextIndex = startIndex + 1
		break
	// ref - ex: Ljava/lang/Object;
	case 'L':
		endIndex := startIndex + 1
		for descriptor[endIndex] != ';' {
			endIndex++
		}
		paramType = descriptor[startIndex : endIndex+1]
		nextIndex = endIndex + 1
		break
		// slice - ex: ([I)V -> should return "[I"
	case '[':
		elementType, idx := parseFieldType(descriptor, startIndex+1)
		paramType = "[" + elementType
		nextIndex = idx
	default:
		panic("Invalid descriptor: " + descriptor)
	}
	return
}

// GetParameterTypes getter
func (md *MethodDescriptor) GetParameterTypes() []string {
	return md.parameterTypes
}

// GetReturnType getter
func (md *MethodDescriptor) GetReturnType() string {
	return md.returnType
}
