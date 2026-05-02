package diff

import (
	"strings"
)

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// DiffWords computes an LCS-based space-separated text diff.
// Returns a formatted string where differing words are highlighted.
// + indicates unexpected word (red), - indicates missing expected word (green).
func DiffWords(expected, actual string) string {
	expWords := strings.Fields(expected)
	actWords := strings.Fields(actual)

	m := len(expWords)
	n := len(actWords)

	// LCS table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if expWords[i-1] == actWords[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// Backtrack to find diff
	i, j := m, n
	var result []string

	for i > 0 && j > 0 {
		if expWords[i-1] == actWords[j-1] {
			// Match, normal text
			result = append([]string{expWords[i-1]}, result...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			// Expected word missing in actual (Green)
			result = append([]string{"\033[32m-" + expWords[i-1] + "\033[0m"}, result...)
			i--
		} else {
			// Actual word not in expected (Red)
			result = append([]string{"\033[31m+" + actWords[j-1] + "\033[0m"}, result...)
			j--
		}
	}

	for i > 0 {
		result = append([]string{"\033[32m-" + expWords[i-1] + "\033[0m"}, result...)
		i--
	}
	for j > 0 {
		result = append([]string{"\033[31m+" + actWords[j-1] + "\033[0m"}, result...)
		j--
	}

	return strings.Join(result, " ")
}
