package suiron

// IntStack - LIFO stack for integers. Used by the tokenizer.
//
// Reference:
// https://www.educative.io/edpresso/how-to-implement-a-stack-in-golang
//
// Cleve Lendon

type IntStack []int

// IsEmpty - check for empty stack.
func (s *IntStack) IsEmpty() bool {
    return len(*s) == 0
}

// Push - pushes an integer onto the stack.
func (s *IntStack) Push(i int) {
    *s = append(*s, i)
}

// Pop - pops an integer from the top of the stack.
// Return: popped integer
//         success flag - true if integer popped
func (s *IntStack) Pop() (int, bool) {
    if len(*s) == 0 {
        return 0, false
    } else {
        index := len(*s) - 1 // Top index.
        element := (*s)[index]
        *s = (*s)[:index] // Remove last item.
        return element, true
    }
} // Pop

// Peek - returns the top integer on the stack.
// Does not pop integer.
// Return: top integer
//         success flag - true if the stack was not empty
func (s *IntStack) Peek() (int, bool) {
    if len(*s) == 0 {
        return 0, false
    } else {
        index := len(*s) - 1 // Top index.
        element := (*s)[index]
        return element, true
    }
} // Peek
