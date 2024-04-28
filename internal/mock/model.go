package mock



// MockDB is a mock implementation of PersonalAllowanceInterface
type MockDB struct {
	Amounts []float64 // Slice to simulate database records
	DefaultAmount float64
}




// Create simulates creating a personal allowance record
func (m *MockDB) Create(amount float64) error {
	m.Amounts = append(m.Amounts, amount) // Append the new amount as a new record
	return nil
}

// Read simulates reading the latest personal allowance record
func (m *MockDB) Read() (float64, error) {
	if len(m.Amounts) == 0 {
		return m.DefaultAmount, nil
	}
	return m.Amounts[len(m.Amounts)-1], nil // Return the last element as the latest record
}

