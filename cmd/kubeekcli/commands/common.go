package commands

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/urfave/cli/v3"

	"kubeekcli/internal/prompt"
)

// resolveGenerateDest determines the destination path for generated files based on CLI flags and arguments.
// It prompts the user if the destination exists and may overwrite, unless the "force" flag is set.
// Parameters:
//
//	c - the CLI command containing flags and arguments.
//
// Returns:
//
//	string - the resolved destination path.
//	error  - any error encountered during resolution or user interaction.
func resolveGenerateDest(c *cli.Command) (string, error) {
	dst := strings.TrimSpace(c.String("out"))
	if dst == "" {
		dst = strings.TrimSpace(c.String("dst"))
	}
	if dst == "" {
		if args := c.Args().Slice(); len(args) > 0 && strings.TrimSpace(args[0]) != "" {
			dst = strings.TrimSpace(args[0])
		}
	}
	if dst == "" {
		name := strings.TrimSpace(c.String("name"))
		if name == "" {
			name = "generated"
		}
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
			dst = name
		} else {
			dst = "./" + name
		}
	}

	for {
		if info, err := os.Stat(dst); err == nil && info.IsDir() && !c.Bool("force") {
			yes, err := prompt.Confirm(fmt.Sprintf("Destination %q exists. Overwrite?", dst), false)
			if err != nil {
				return "", err
			}
			if yes {
				break
			}
			for {
				newDst, err := prompt.Prompt("Enter a NEW destination folder (or blank to cancel): ")
				if err != nil {
					return "", err
				}
				if newDst == "" {
					return "", fmt.Errorf("aborted by user")
				}
				if _, err := os.Stat(newDst); os.IsNotExist(err) {
					dst = newDst
					break
				}
				fmt.Printf("Path %q already exists. ", newDst)
			}
			continue
		}
		break
	}
	return dst, nil
}

func parseDefaults(s string) map[string]string {
	res := map[string]string{}
	if s == "" {
		return res
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		k := strings.TrimSpace(kv[0])
		v := ""
		if len(kv) == 2 {
			v = strings.TrimSpace(kv[1])
		}
		if k != "" {
			res[k] = v
		}
	}
	return res
}

func All() []*cli.Command {
	return []*cli.Command{
		GenerateCmd(),
	}
}
