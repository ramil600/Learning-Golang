package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//The interface of internal layer that is beyond this unit test, needs to be mocked
type Calculator interface {
	Calculate(int) int
}

//We enclose Mock object in the struct that will be mocking Calculator interface
type calculatorMock struct {
	mock.Mock
}

//Defining mock Calculate method parameters and outputs, implementing Calculator interface
func (c *calculatorMock) Calculate(j int) int {
	args := c.Called(j)
	return args.Int(0)
}

// The interface Calculator is internal layer that we are not interested in testing
//so this is the reason for mocking this interface
type Person struct {
	calculator Calculator
}

//Function to be tested, it calls internal layer's Calculate method
func (p Person) AddTwo(j int) int {
	return p.calculator.Calculate(j)
}

//Dependency injection of internal layer and creating outer layer object
func NewPerson(calculator Calculator) *Person {
	return &Person{
		calculator: calculator,
	}
}

//Testing whether expected and actual result of AddTwo opeation is Equal
func TestAddTwo(t *testing.T) {
	assertions := require.New(t)

	tests := []struct {
		a    int
		want int
	}{
		{
			a:    2,
			want: 4,
		},
		{
			a:    3,
			want: 5,
		},
		{
			a:    4,
			want: 6,
		},
	}

	mock := calculatorMock{}

	//Benefit of mock object is that we wire it inside the test, so all test resides within the function
	//and only visible here. New scenarios can be created in other tests according to testing needs
	mock.On("Calculate", 2).Return(4)
	mock.On("Calculate", 3).Return(5)
	mock.On("Calculate", 4).Return(6)

	//We substitute real Calculator with mock Calculator by implementing Calculator interface
	p := NewPerson(&mock)

	for _, tt := range tests {
		//testify require package provides clean and readable assertions for testing
		assertions.Equal(p.AddTwo(tt.a), tt.want)
	}
}
