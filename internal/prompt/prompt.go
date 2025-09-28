package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func Confirm(prompt string, def bool) (bool, error) {
	suffix := "y/N"
	if def {
		suffix = "Y/n"
	}

	ans, err := Prompt(fmt.Sprintf("%s [%s]: ", prompt, suffix))

	if err != nil {
		return def, err
	}
	a := strings.ToLower(strings.TrimSpace(ans))

	switch a {
	case "", "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return def, nil
	}
}

func ConfirmOverwrite(dst string) (bool, error) {
	return Confirm(fmt.Sprintf("Destination %q exists. Overwrite?", dst), false)
}
