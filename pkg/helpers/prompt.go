package helpers

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func PromptRemovalOf(subject string, forceOverride bool) bool {
	if forceOverride {
		return true
	}

	return PromptConfirm(fmt.Sprintf("Do you want to remove %s", subject), forceOverride)
}

func PromptConfirm(q string, forceOverride bool) bool {
	if forceOverride {
		return true
	}

	p := promptui.Prompt{
		IsConfirm: true,
		Label:     q,
	}

	res, err := p.Run()
	if err != nil {
		return false
	}

	resConv := strings.ToUpper(res)
	return resConv == "Y" || resConv == "YES"
}
