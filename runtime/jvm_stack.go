package runtime

type JVMStack struct {
	maxSize uint   // max size of jvm stack
	size    uint   // current stack size
	top     *Frame // top of the frame stack
}

// NewJVMStack create a new JVMStack
func NewJVMStack(maxSize uint) *JVMStack {
	return &JVMStack{
		maxSize: maxSize,
		size:    0,
	}
}

func (s *JVMStack) Push(frame *Frame) {
	// check StackOverFlow
	if s.size >= s.maxSize {
		panic("java.lang.StackOverflowError")
	}

	currentFrame := s.top
	if currentFrame != nil {
		// link new Frame with current Frame
		frame.lower = currentFrame
	}

	s.top = currentFrame
	s.size++
}

func (s *JVMStack) Pop() *Frame {
	if s.IsEmpty() {
		panic("java.lang.StackUnderflowError")
	}

	poppedFrame := s.top
	s.top = poppedFrame.lower

	// disconnect popped frame and previous frame
	poppedFrame.lower = nil
	s.size--

	return poppedFrame
}

func (s *JVMStack) Top() *Frame {
	if s.IsEmpty() {
		panic("java.lang.StackUnderflowError")
	}

	return s.top
}

func (s *JVMStack) IsEmpty() bool {
	return s.size == 0 && s.top == nil
}

func (s *JVMStack) Size() uint {
	return s.size
}

// Clear clear stack
func (s *JVMStack) Clear() {
	for !s.IsEmpty() {
		s.Pop()
	}
}

// GetFrames get all frame (for tracing exception call stack)
// control order: top at first
func (s *JVMStack) GetFrames() []*Frame {
	frames := make([]*Frame, 0, s.size)
	for frame := s.top; frame != nil; frame = frame.lower {
		frames = append(frames, frame)
	}
	return frames
}
