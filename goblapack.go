package goblapack

import (
	"errors"
	"syscall"
	"unsafe"
)

// Matrix is a dense matrix structure that holds the data in column-major order
// in a linear array. Because of this lay
type Matrix struct {
	Rows int
	Cols int
	Data []float64
}

// NewMatrix creates a new matrix of the give size.
func NewMatrix(rows, cols int) *Matrix {
	size := rows * cols
	m := Matrix{Rows: rows, Cols: cols}
	m.Data = make([]float64, size, size)
	return &m
}

// Get returns the value at the given row and column.
func (m *Matrix) Get(row, col int) float64 {
	i := row + m.Rows*col
	return m.Data[i]
}

// Set sets the matrix cell at the given row and column to the given value.
func (m *Matrix) Set(row, col int, value float64) {
	i := row + m.Rows*col
	m.Data[i] = value
}

// Copy creates a copy of the matrix.
func (m *Matrix) Copy() *Matrix {
	c := NewMatrix(m.Rows, m.Cols)
	copy(c.Data, m.Data)
	return c
}

// Multiply calculates the matrix-matrix-product C = A * B where A is the matrix
// on which the method is called, B the method parameter, and C the return value.
func (m *Matrix) Multiply(b *Matrix) (*Matrix, error) {
	if m.Cols != b.Rows {
		return nil, errors.New("Cannot multiply matrices: shapes do not match")
	}
	handle, err := syscall.GetProcAddress(nativeLib, "goblapack_mmult")
	if err != nil {
		return nil, err
	}
	c := NewMatrix(m.Rows, b.Cols)
	aPtr := uintptr(unsafe.Pointer(&m.Data[0]))
	bPtr := uintptr(unsafe.Pointer(&b.Data[0]))
	cPtr := uintptr(unsafe.Pointer(&c.Data[0]))
	syscall.Syscall6(uintptr(handle), 6,
		uintptr(m.Rows),
		uintptr(b.Cols),
		uintptr(m.Cols),
		aPtr,
		bPtr,
		cPtr)
	return c, nil
}

// Invert calculates the inverse of the matrix.
func (m *Matrix) Invert() (*Matrix, error) {
	inverse := m.Copy()
	err := inverse.InvertInPlace()
	return inverse, err
}

// InvertInPlace calculates the inverse of the matrix which is stored directly
// in the original matrix.
func (m *Matrix) InvertInPlace() error {
	if m.Cols != m.Rows {
		return errors.New("The matrix is not square")
	}
	handle, err := syscall.GetProcAddress(nativeLib, "goblapack_invert")
	if err != nil {
		return err
	}
	dataPtr := uintptr(unsafe.Pointer(&m.Data[0]))
	r, _, _ := syscall.Syscall(uintptr(handle), 2, uintptr(m.Rows), dataPtr, 0)
	info := int(r)
	if info > 0 {
		return errors.New("Matrix is singular")
	}
	if info < 0 {
		return errors.New("Invalid data input")
	}
	return nil
}
