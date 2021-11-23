package main

import (
	"strings"

	"golang.org/x/net/html"
)

func escapedIndex(s, substr string) int {
	i := strings.Index(s, substr)
	if i == -1 {
		return i
	}

	if i != 0 && s[i-1] == '\\' {
		i += len(substr)
		return i + escapedIndex(s[i:], substr)
	}

	return i
}

func getValueInside(s, sepa, sepb string) string {
	// fmt.Println(s)
	i := escapedIndex(s, sepa)
	if i == -1 {
		return ""
	}
	i += len(sepa)
	s = s[i:]
	// fmt.Println(s)
	i = escapedIndex(s, sepb)
	if i == -1 {
		return ""
	}
	return s[:i]
}

func getValueInsideQuotes(s string) string {
	return getValueInside(s, `"`, `";`)
}

func parseT1Assertion(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		i := strings.Index(src, "t1Assertion")
		if i != -1 {
			return getValueInsideQuotes(src[i:])
		}
	}
	return ""
}

func parseKeyboardContent(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		i := strings.Index(src, `KEYCONTENT = "`)
		if i != -1 {
			content := getValueInsideQuotes(src[i:])
			if strings.Contains(content, "keyboard") {
				return content
			}
		}
	}
	return ""
}

func parseKeyboardMap(node *html.Node) map[string]string {
	keyMap := map[string]string{}
	elms := getAllElemenstByTag(node, "td")

	for _, elm := range elms {
		text := strings.TrimSpace(getInnerText(elm))
		// fmt.Println(text)
		// number, err := strconv.Atoi(text)
		// if err != nil {
		// 	continue
		// }
		// fmt.Println("n:", number)
		onclick := getAttribute(elm, "onclick")
		if onclick == "" {
			continue
		}
		value := getValueInside(onclick, `\"`, `\"`)
		// onclick = strings.ReplaceAll(onclick, `\`, "")
		// fmt.Println(onclick)
		// value := getValueInsideQuotes(onclick)
		// fmt.Println(value)
		// fmt.Printf("%d: %s\n", number, value)
		// keyMap[strconv.Itoa(number)] = value
		if value != "" {
			keyMap[text] = value
		}
	}

	return keyMap

}

func parsePasswordInputName(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "'PASSWORD':"
		i := strings.Index(src, ptn)
		if i != -1 {
			return getValueInside(src[i+len(ptn):], `'`, `'`)
		}
	}
	return ""
}

func parseUrlRedirect(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "urlRedirect ="
		i := strings.Index(src, ptn)
		if i != -1 {
			urlRedirect := getValueInside(src[i+len(ptn):], `'`, `'`)
			if urlRedirect != "" {
				return urlRedirect
			}
		}
	}
	return ""
}

func parseTokenMua(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "tokenMua="
		i := strings.Index(src, ptn)
		if i != -1 {
			value := getValueInside(src[i+len(ptn):], `'`, `'`)
			if value != "" {
				return value
			}
		}
	}
	return ""
}

func parseCodeRedirect(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "codeRedirect="
		i := strings.Index(src, ptn)
		if i != -1 {
			value := getValueInside(src[i+len(ptn):], `'`, `'`)
			if value != "" {
				return value
			}
		}
	}
	return ""
}

// parseLocationReplace searches for a script tag with a location.replace statement
// and returns the replacement string
func parseLocationReplace(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "location.replace"
		i := strings.Index(src, ptn)
		if i != -1 {
			value := getValueInside(src[i+len(ptn):], `"`, `"`)
			if value != "" {
				return value
			}
		}
	}
	return ""
}

func parseCstParam(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "'cst="
		i := strings.Index(src, ptn)
		if i != -1 {
			value := getValueInside(src[i:], `=`, `'`)
			if value != "" {
				return value
			}
		}
	}
	return ""
}

func parseCsrfToken(node *html.Node) string {
	elms := getAllElemenstByTag(node, "input")
	for _, elm := range elms {
		name := getAttribute(elm, "name")
		value := getAttribute(elm, "value")
		if name == "CSRF_TOKEN" && value != "" {
			return value
		}
	}
	return ""
}

// parseTokenValue extracts the value of the token argument
// present in the given node with the following format: {'TOKEN': 'value'}
func parseTokenValue(node *html.Node) string {
	scripts := getAllElemenstByTag(node, "script")
	for _, script := range scripts {
		src := getInnerText(script)
		ptn := "'TOKEN':"
		i := strings.Index(src, ptn)
		if i != -1 {
			value := getValueInside(src[i+len(ptn):], `'`, `'}`)
			if value != "" {
				return value
			}
		}
	}
	return ""
}
