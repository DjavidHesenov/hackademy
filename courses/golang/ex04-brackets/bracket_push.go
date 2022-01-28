package brackets

func Bracket(str string) (bool, error) {
	bracketSlice := Stack{}

	for _, bracketCode := range str {
		if bracketCode == '(' || bracketCode == '[' || bracketCode == '{' {
			bracketSlice.Push(int(bracketCode))
			continue
		}

		lastBracket := string(rune(bracketSlice.Pop()))

		switch string(bracketCode) {
		case ")":
			if lastBracket != "(" {
				return false, nil
			}

		case "]":
			if lastBracket != "[" {
				return false, nil
			}

		case "}":
			if lastBracket != "{" {
				return false, nil
			}
		}
	}

	return bracketSlice.IsEmpty(), nil
}
