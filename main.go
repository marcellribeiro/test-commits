package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Row struct {
	Timestamp  string
	Username   string
	Repository string
	Files      int
	Additions  int
	Deletions  int
}

type scoreRepository struct {
	Repository string
	Score      float64
	Files      int
	Additions  int
	Deletions  int
}

func parseCSV(fileName string) ([]Row, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var commits []Row
	var mu sync.Mutex
	wg := sync.WaitGroup{}

	for i, line := range lines {
		if i == 0 {
			continue
		}

		wg.Add(1)
		go func(line []string) {
			defer wg.Done()

			files := getIntFromStr(line[3])
			additions := getIntFromStr(line[4])
			deletions := getIntFromStr(line[5])

			commit := Row{
				Timestamp:  line[0],
				Username:   line[1],
				Repository: line[2],
				Files:      files,
				Additions:  additions,
				Deletions:  deletions,
			}

			mu.Lock()
			commits = append(commits, commit)
			mu.Unlock()
		}(line)
	}

	wg.Wait()

	return commits, nil
}

func getIntFromStr(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func calculateScores(commits []Row) (map[string]float64, map[string]float64, map[string]float64, map[string]float64) {
	scores := make(map[string]float64)
	files := make(map[string]float64)
	additions := make(map[string]float64)
	deletions := make(map[string]float64)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, commit := range commits {
		wg.Add(1)

		go func(c Row) {
			defer wg.Done()

			// get score
			score := 1 + (0.10 * float64(c.Files)) + (0.01 * float64(c.Additions+c.Deletions))

			mu.Lock()
			scores[c.Repository] += score
			files[c.Repository] += float64(c.Files)
			additions[c.Repository] += float64(c.Additions)
			deletions[c.Repository] += float64(c.Deletions)
			mu.Unlock()
		}(commit)
	}

	wg.Wait()
	return scores, files, additions, deletions
}

func sortScores(scores map[string]float64, files map[string]float64, additions map[string]float64, deletions map[string]float64) []scoreRepository {
	var sortedScores []scoreRepository

	for repo, score := range scores {
		sortedScores = append(sortedScores, scoreRepository{
			Repository: repo,
			Score:      score,
			Files:      int(files[repo]),
			Additions:  int(additions[repo]),
			Deletions:  int(deletions[repo]),
		})
	}

	sort.Slice(sortedScores, func(i, j int) bool {
		return sortedScores[i].Score > sortedScores[j].Score
	})

	return sortedScores
}

func main() {
	start := time.Now()

	commits, err := parseCSV("commits.csv")
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	scores, files, additions, deletions := calculateScores(commits)

	sortedScores := sortScores(scores, files, additions, deletions)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Top10 Repository Scores (Descending Order):")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf(
		"%2s - %-9s - %10s - %5s - %10s - %10s\n",
		"#", "REPO", "SCORE", "FILES", "ADDITIONS", "DELETIONS",
	)
	for i, score := range sortedScores {
		if i == 10 {
			break
		}
		fmt.Printf(
			"%2d - %-9s - %10.2f - %5d - %10d - %10d\n",
			i+1, score.Repository, score.Score, score.Files, score.Additions, score.Deletions,
		)
	}

	elapsed := time.Since(start)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Elapsed time: %s\n", elapsed)
}
