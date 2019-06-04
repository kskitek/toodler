package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const regex = "TODO"

type SearchResult struct {
	Hash       string
	Author     string
	Date       time.Time
	FileName   string
	LineNumber int
	Line       string
}

type Searcher struct {
}

func (s *Searcher) Search(repositoryUrl string) ([]SearchResult, error) {
	// TODO create log file under /var/logs/projectName and log with log.New(w..)
	projectName := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(repositoryUrl, ":", ""), "/", ""), ".", "")

	projectPath, err := s.cloneRepository(repositoryUrl, projectName)
	if err != nil {
		return nil, err
	}
	resultFile, err := s.runCommand(projectPath, projectName)
	if err != nil {
		return nil, err
	}
	result, err := s.parseResults(resultFile)

	return result, nil
}

func (s *Searcher) cloneRepository(url, projectName string) (projectPath string, err error) {
	log.Printf("Cloning repository: " + url)
	projectPath = "/tmp/repos/" + projectName // TODO config or const
	var cmd *exec.Cmd
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		cmd = exec.Command("git", "clone", url, projectPath)
	} else {
		cmd = exec.Command("git", "pull")
		cmd.Dir = projectPath
		// TODO if "Already up to date" I don't have to look for TODOs.. return an error and handle in Search
    // TODO maybe in case of update it would be easier to find TODOs on diff?
	}
	cmd.Stdout = os.Stdout // TODO write me to log file per project
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf(err.Error())
	}
	return
}

func (s *Searcher) runCommand(path, projectName string) (resultFile string, err error) {
	log.Printf("Looking for TODOs in: " + path)
	resultFile = "/tmp/blames/" + projectName
	cmd := exec.Command("/usr/local/bin/findIt.sh", path, regex, resultFile)
	cmd.Stdout = os.Stdout // TODO write me to log file per project
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf(err.Error())
		return "", err
	}

	return
}

func (s *Searcher) parseResults(resultFile string) ([]SearchResult, error) {
	result := make([]SearchResult, 0)
	log.Printf("Reading results: " + resultFile)
	f, err := os.Open(resultFile)
	if err != nil {
		log.Printf(err.Error())
		return make([]SearchResult, 0), err
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		res, err := s.parseLine(line)
		if err != nil {
			continue
		}
		result = append(result, *res)
	}
	log.Printf("Found [%d] TODOs", len(result))
	return result, nil
}

func (s *Searcher) parseLine(line string) (*SearchResult, error) {
  // TODO verify why empty TODOs are returned
	res := &SearchResult{}
	res.Hash = line[:40]
	line = line[41 : len(line)-1]

	sp := strings.Index(line, " ")
	res.FileName = line[:sp]
	line = line[sp : len(line)-1]

	ob := strings.Index(line, "(")
	cb := strings.Index(line, ")")
	lineNumber, err := strconv.Atoi(strings.Trim(line[:ob-1], " "))
	if err != nil {
		log.Printf("Cannot parse linenumber: " + line[:ob-1])
		return nil, err
	}
	res.LineNumber = lineNumber

	lt := strings.Index(line, "<")
	gt := strings.Index(line, ">")
	res.Author = line[lt+1 : gt]

	dateChunks := strings.Split(strings.Trim(line[gt+1:cb], " "), " ")
	dateStr := strings.Join(dateChunks[:3], " ")
	date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
	if err != nil {
		log.Printf("Cannot parse date: " + dateStr)
		return nil, err
	}
	res.Date = date

	res.Line = strings.Trim(line[cb+1:], " ")
	return res, nil
}
