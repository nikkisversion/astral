package main

// MultiStruct block tests clustered type declarations
type (
	User   struct{ Name string }
	Config struct{ Port int }
)

// SingleStruct tests a standard isolated declaration
type SingleStruct struct{ ID int }

// EmptyStruct tests zero-field syntax
type EmptyStruct struct{}

// inlineComment tests preservation
func inlineComment() {
	// Inside block
}

// InlineFunc tests one-liner declarations
func InlineFunc() int { return 42 }

// NLTestFunc tests handling of newlines in function bodies
func NLTestFunc(nlString string) string {
	for _, char := range nlString {
		if char == '\n' {
			return "Contains newline"
		}
	}
	return "No newline found"
}
