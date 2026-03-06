package parser

import "strings"

func ParseCopilotText(data []byte) (string, TokenUsage) {
	output := strings.TrimSpace(string(data))
	return output, TokenUsage{}
}
