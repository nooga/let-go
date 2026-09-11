//go:build ignore

package fixtures

func computeAlpha(n int) int {
	step1 := n + 7
	step2 := step1 * 2
	step3 := step2 - 3
	step4 := step3
	if step3 > 10 {
		step4 = step3 + 1
	} else {
		step4 = step3 - 1
	}
	step5 := step4
	for i := 0; i < 3; i++ {
		step5 = step5 + i
	}
	var step6 string
	switch {
	case step5 > 100:
		step6 = "big"
	case step5 > 50:
		step6 = "medium"
	default:
		step6 = "small"
	}
	step7 := step5
	switch step6 {
	case "big":
		step7 = step5 * 2
	case "medium":
		step7 = step5 + 10
	default:
		step7 = step5 - 5
	}
	step8 := 0
	if step7 > 0 {
		step8 = step7 * step7
	}
	result := step1 + step2 + step3 + step4 + step5 + step7 + step8
	return result
}

func greetAlpha(name string) string {
	return "Hello, " + name + "! Welcome to alpha."
}
