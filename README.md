# Repository Scoring Algorithm Documentation

## Overview

This Go program calculates the activity scores of repositories based on data
from a commits.csv file. It ranks the repositories by their scores and outputs
the top 10 in descending order. The score is computed using the following formula:

```math
Score = 1 + (0.10 * Files_Modified) + (0.01 * (Additions + Deletions))
```

## Features

### Multithreading

Utilizes goroutines and sync primitives (sync.Mutex and sync.WaitGroup) for
concurrent processing of the CSV file and score computation.

### CSV Parsing

Reads and processes commit data from a CSV file containing fields: Timestamp,
Username, Repository, Files, Additions, Deletions.


| Field        | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `timestamp`  | The unix timestamp of the commit                             |
| `user`       | The GitHub username of the commit author. This data is unreliable and blank indicates the author's identity is unknown. |
| `repository` | The repository name which the commit was pushed to.          |
| `files`      | The number of files changes by the commit                    |
| `additions`  | The number of line additions in this commit                  |
| `deletions`  | The number of deletions in this commit                       |


### Sorting

Scores are calculated and sorted in descending order, with additional
statistics for modified files, additions, and deletions.

### Performance Measurement

The program logs the elapsed time for execution.

## Input

The program expects a CSV file named commits.csv in the following format:

```csv
Timestamp,Username,Repository,Files,Additions,Deletions
1607531861,user1,repo1,5,100,50
1607531261,user2,repo2,3,150,100
1607124554,user3,repo1,2,50,25
```

## Output

The program outputs the top 10 repositories with the highest scores in a
table format:

```bash
--------------------------------------------------------------------------------
Top10 Repository Scores (Descending Order):
--------------------------------------------------------------------------------
# - REPO       -    SCORE - FILES -  ADDITIONS -  DELETIONS
1 - repo1      -   123.45 -    10 -        200 -        150
2 - repo2      -    98.76 -     8 -        150 -        100
--------------------------------------------------------------------------------
Elapsed time: 10ms
```

## How to Run

### Prerequisites

- Install Go (version 1.18+ recommended).
- Ensure the CSV file commits.csv exists in the same directory as the
- program.
- Setup: Save the program code to a file named main.go.
- Run the Program: Execute the following commands:

```bash
go run main.go
```

## Performance Optimization

The program uses:

- Goroutines: Concurrently processes CSV rows and computes scores for
improved performance on large datasets.
- Mutexes: Safely updates shared data structures across threads.

## Code Explanation

### File Parsing

The parseCSV function reads and processes the commits.csv file using
goroutines to parse each row into a Row struct.

### Score Calculation

The calculateScores function computes repository scores concurrently,
aggregating results into maps for scores, files, additions, and deletions.

### Sorting the score

The sortScores function organizes repositories by descending scores into
a slice of scoreRepository structs.
