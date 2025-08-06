package main

import (
	"fmt"
	"log"
	"log/slog"
	"slices"
	"strings"

	"github.com/xuri/excelize/v2"
)

type course struct {
	Abbr     string
	FullName string
}

func getAllCourcesByPriorty(filename string, tags []string, passedCourses map[string]struct{}) ([][]course, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := f.GetRows("Table 1")
	if err != nil {
		return nil, err
	}
	if len(rows) <= 2 {
		return nil, fmt.Errorf(" rows less than 2 in the file")
	}
	rows = rows[2:]

	prior := make([][]course, 4)

	for _, row := range rows {
		if len(row) < 6 {
			continue
		}

		abbr := escapeNewLine(row[1])
		fullname := escapeNewLine(row[2])

		if req := escapeNewLine(row[5]); req != "" { //req
			req = simplifyReq(req)
			if !Resolve(req, passedCourses) {
				continue
			}
		}

		if req := escapeNewLine(row[7]); req != "" { //anti-req
			req = simplifyReq(req)
			if Resolve(req, passedCourses) {
				slog.Info("Skipping course due to anti-req", "course", abbr, "anti-req", req)
				continue
			}
		}

		foundAtPrior := 3
		for i := 8; i < min(len(row), 12); i++ {
			priorTags := row[i]
			priorTags = escapeNewLine(priorTags)
			if contains(priorTags, tags) {
				foundAtPrior = i - 8
				break
			}
		}

		prior[foundAtPrior] = append(prior[foundAtPrior], course{
			Abbr:     abbr,
			FullName: fullname,
		})
	}

	return prior, nil
}

func contains(s string, targets []string) bool {
	for tag := range strings.SplitSeq(s, ", ") {
		if slices.Contains(targets, tag) {
			return true
		}
	}
	return false
}

func main() {
	// req := "LING 273 Survey of Research Methods in Linguistics (2158) (C- and above) AND (LING 375 The Art and Science of Analyzing Languages: Morphosyntax of the World's Languages (6237) (C- and above) OR LING 377 Historical Linguistics (7972) (C- and above) OR LING 461 Experimental semantics (7196) (C- and above) OR LING 473 Advanced Empirical Methods in Linguistics (6244) (C- and above))"

	// fmt.Println(simplifyReq(req))

	// return
	reqFileName := "req.xlsx"
	passedCourses := ToMap([]string{
		"WCS 150",
		"HST 100",
		"MATH 161",
		"MATH 162",
		"PHYS 161",
		"PHYS 162",
		"CSCI 151",
		"CSCI 152",
		"MATH 251",
		"MATH 273",
	})

	tags := []string{
		"2 year UG SEDS",
		"SEDS",
		"Computer Science",
	}

	priors, err := getAllCourcesByPriorty(reqFileName, tags, passedCourses)
	if err != nil {
		log.Fatal(err)
	}
	for i, prior := range priors {
		fmt.Printf("Prior %d: %d courses\n", i+1, len(prior))
		for _, course := range prior {
			fmt.Printf("\t%-20s (%s)\n", course.Abbr, course.FullName)
		}
		fmt.Println()
	}
}

func ToMap(ss []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range ss {
		m[s] = struct{}{}
	}
	return m
}

func escapeNewLine(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}
