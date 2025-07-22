package handler

import (
	"fmt"
	"os"
	"strings"
)

func isLikelyJSON(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 1)
	for {
		n, err := f.Read(buf)
		if n == 0 || err != nil {
			return false
		}
		// Пропускаем пробелы и переносы строк
		if buf[0] == ' ' || buf[0] == '\n' || buf[0] == '\r' || buf[0] == '\t' {
			continue
		}
		return buf[0] == '{' || buf[0] == '['
	}
}

// versionMatches проверяет, подходит ли версия под условие
func versionMatches(ver, cond string) bool {
	if cond == "" || cond == ver {
		return true
	}
	cond = strings.TrimSpace(cond)
	if strings.HasPrefix(cond, ">=") {
		return !versionLess(ver, cond[2:])
	}
	if strings.HasPrefix(cond, "<=") {
		return !versionLess(cond[2:], ver)
	}
	if strings.HasPrefix(cond, ">") {
		return versionLess(cond[1:], ver)
	}
	if strings.HasPrefix(cond, "<") {
		return versionLess(ver, cond[1:])
	}
	if strings.HasPrefix(cond, "=") {
		return ver == cond[1:]
	}
	return ver == cond
}

// versionLess сравнивает версии как числа с точками
func versionLess(a, b string) bool {
	aa := strings.Split(a, ".")
	bb := strings.Split(b, ".")
	for i := 0; i < len(aa) || i < len(bb); i++ {
		var ai, bi int
		if i < len(aa) {
			fmt.Sscanf(aa[i], "%d", &ai)
		}
		if i < len(bb) {
			fmt.Sscanf(bb[i], "%d", &bi)
		}
		if ai < bi {
			return true
		}
		if ai > bi {
			return false
		}
	}
	return false
}
