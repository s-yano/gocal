package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"strings"
	"time"
)

func bright(s string) string {
	return fmt.Sprintf("\x1b[1;4m%v\x1b[0m", s)
}

func reverse(s string) string {
	return fmt.Sprintf("\x1b[7m%v\x1b[0m", s)
}

func centering(s string, w int) string {
	l := len(s)
	if l >= w {
		return s
	}
	s1 := (w - l) / 2
	s2 := w - (l + s1)
	return strings.Repeat(" ", s1) + s + strings.Repeat(" ", s2)
}

func is_include(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func mkbuf() []string {
	// return make([]string, 0, 7)
	return []string{"  ", "  ", "  ", "  ", "  ", "  ", "  "}
}

func index(d time.Time) int {
	return int((d.Weekday() + 6) % 7)
}

func firstDay(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func build_cal2(year int, month time.Month, holiday []string) []string {
	today_y, today_m, today_d := time.Now().Date()
	var buf []string
	buf = append(buf, centering(fmt.Sprintf("%v %d", month, year), 20))
	buf = append(buf, "Mo Tu We Th Fr Sa Su")
	s := mkbuf()
	for d := firstDay(year, month); d.Month() == month; d = d.AddDate(0, 0, 1) {
		i := (d.Weekday() + 6) % 7
		s[i] = fmt.Sprintf("%2d", d.Day())
		if d.Weekday() == 6 || d.Weekday() == 0 || is_include(fmt.Sprintf("%4d-%.2d-%.2d", year, month, d.Day()), holiday) {
			s[i] = bright(s[i])
		}
		if today_y == year && today_m == month && today_d == d.Day() {
			s[i] = reverse(s[i])
		}
		if d.Weekday() == 0 {
			buf = append(buf, strings.Join(s, " "))
			s = mkbuf()
		}
	}
	if len(s) > 0 {
		buf = append(buf, strings.Join(s, " "))
	}
	return buf
}

func build_cal(year int, month time.Month, label bool, holiday []string) []string {
	var buf []string
	// var str string
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	weekday := int(first.Weekday())
	last := first.AddDate(0, 1, -1)
	today_y, today_m, today_d := time.Now().Date()

	label_str := fmt.Sprintf("%v", first.Month())
	if label {
		label_str += fmt.Sprintf(" %v", year)
	}
	buf = append(buf, centering(label_str, 20))

	buf = append(buf, "Mo Tu We Th Fr Sa Su")
	str := strings.Repeat("   ", (weekday+6)%7)
	for d := 1; d <= last.Day(); d++ {
		day := fmt.Sprintf("%2d", d)
		if isatty.IsTerminal(os.Stdout.Fd()) {
			if weekday == 6 || weekday == 0 || is_include(fmt.Sprintf("%4d-%.2d-%.2d", year, month, d), holiday) {
				day = bright(day)
			}
			if today_y == year && today_m == month && today_d == d {
				day = reverse(day)
			}
		}
		str += day
		if weekday == 0 {
			buf = append(buf, str)
			str = ""
		} else {
			str += " "
		}
		weekday = (weekday + 1) % 7
	}
	if str != "" {
		buf = append(buf, str+strings.Repeat(" ", 20-len(str)))
	}
	return buf
}

func read_config() []string {
	var buf []string

	fp, err := os.Open(os.Getenv("HOME") + "/.gocal")
	if err != nil {
		return buf
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 || text[0] == '#' {
			continue
		}
		buf = append(buf, text)
	}
	return buf
}

func main() {
	holiday := read_config()

	var buf [3][]string
	cur := time.Now().AddDate(0, -1, 0)
	for i := 0; i < 3; i++ {
		y, m, _ := cur.Date()
		// buf[i] = build_cal(y, m, true, holiday)
		buf[i] = build_cal2(y, m, holiday)
		cur = cur.AddDate(0, 1, 0)
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 3; j++ {
			if i >= len(buf[j]) {
				fmt.Printf(strings.Repeat(" ", 22))
			} else {
				fmt.Printf("%s  ", buf[j][i])
			}
		}
		fmt.Printf("\n")
	}
}
