package instructions

type BytecodeReader struct {
	code []byte // bytecode (from Code Attribute)
	pc   int    // program counter
}
