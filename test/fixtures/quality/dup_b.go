//go:build ignore

package fixtures

func computeBeta(m int) int {
	s1 := m + 9
	s2 := s1 * 2
	s3 := s2 - 3
	s4 := s3
	if s3 > 10 {
		s4 = s3 + 1
	} else {
		s4 = s3 - 1
	}
	s5 := s4
	for j := 0; j < 3; j++ {
		s5 = s5 + j
	}
	var s6 string
	switch {
	case s5 > 100:
		s6 = "big"
	case s5 > 50:
		s6 = "medium"
	default:
		s6 = "small"
	}
	s7 := s5
	switch s6 {
	case "big":
		s7 = s5 * 2
	case "medium":
		s7 = s5 + 10
	default:
		s7 = s5 - 5
	}
	s8 := 0
	if s7 > 0 {
		s8 = s7 * s7
	}
	answer := s1 + s2 + s3 + s4 + s5 + s7 + s8
	return answer
}

func farewellBeta(name string) string {
	return name + ", thanks for visiting beta!"
}
