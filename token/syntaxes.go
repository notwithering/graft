package token

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	HTMLTagSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`(</?g-.+?\s*/?>)`),
		Close:     regexp.MustCompile(`</g-.+?>`),
		Parse: func(s string) (command string, data any, err error) {
			tokenizer := html.NewTokenizer(strings.NewReader(s))

			for {
				tokenType := tokenizer.Next()

				switch tokenType {
				case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
					name, hasAttr := tokenizer.TagName()
					command = string(bytes.TrimLeft(name, "g-"))

					args := make(map[string]any)

					for hasAttr {
						key, val, more := tokenizer.TagAttr()
						keyStr, valStr := string(key), string(val)

						if valStr == "" {
							args[keyStr] = true
						} else {
							args[keyStr] = valStr
						}

						hasAttr = more
					}

					return command, args, nil
				case html.ErrorToken:
					return "", nil, nil
				}
			}
		},
	}
	DoubleBraceSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`\{\{.+?\}\}`),
		Close:     regexp.MustCompile(`\{\{end\}\}`),
		Parse: func(s string) (string, any, error) {
			s = strings.TrimPrefix(s, "{{")
			s = strings.TrimSuffix(s, "}}")
			s = strings.TrimSpace(s)

			fields := strings.Fields(s)
			if len(fields) == 0 {
				return "", nil, nil
			}

			command := fields[0]
			attributes := make(map[string]any)

			i := 1
			for i < len(fields) {
				field := fields[i]
				if strings.Contains(field, "=") {
					keyValuePair := strings.SplitN(field, "=", 2)
					key := keyValuePair[0]
					value := keyValuePair[1]

					if strings.HasPrefix(value, `"`) && !strings.HasSuffix(value, `"`) {
						for j := i + 1; j < len(fields); j++ {
							value += " " + fields[j]
							i = j
							if strings.HasSuffix(fields[j], `"`) {
								break
							}
						}
					}

					value = strings.Trim(value, `"`)
					attributes[key] = value
				} else {
					attributes[field] = true
				}
				i++
			}

			return command, attributes, nil
		},
	}
)
