package postfix

import (
	"fmt"
	"strings"
)

var opPrecedence = map[rune]int{
	'|': 1,
	'.': 2,
	'*': 3,
	'+': 3,
}

func ToPostfix(regex string) (string, error) {
	var output strings.Builder
	var operators []rune

	formattedRegex := insertConcat(regex)

	for _, token := range formattedRegex {
		if isOperand(token) {
			output.WriteRune(token)
		} else if token == '(' {
			operators = append(operators, token)
		} else if token == ')' {
			foundParen := false
			for len(operators) > 0 {
				top := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				if top == '(' {
					foundParen = true
					break
				}
				output.WriteRune(top)
			}
			if !foundParen {
				return "", fmt.Errorf("sintax error: wrong brackets number or order")
			}
		} else {
			for len(operators) > 0 {
				top := operators[len(operators)-1]
				if top == '(' || opPrecedence[top] < opPrecedence[token] {
					break
				}
				output.WriteRune(top)
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		}
	}

	for len(operators) > 0 {
		top := operators[len(operators)-1]
		if top == '(' {
			return "", fmt.Errorf("sintax error: wrong brackets number or order")
		}
		output.WriteRune(top)
		operators = operators[:len(operators)-1]
	}

	return output.String(), nil
}

func isOperand(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == 'ε'
}

func insertConcat(regex string) string {
	var result strings.Builder
	for i := 0; i < len(regex); i++ {
		result.WriteByte(regex[i])
		if i+1 < len(regex) {
			curr := rune(regex[i])
			next := rune(regex[i+1])

			if (isOperand(curr) || curr == ')' || curr == '*' || curr == '+') && (isOperand(next) || next == '(') {
				result.WriteRune('.')
			}
		}
	}
	return result.String()
}
