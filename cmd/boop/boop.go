package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type config struct {
	Endpoint string `yaml:"endpoint"`
	Phone    string `yaml:"phone"`

	Reversed bool
	Limit    int
}

type Boop struct {
	Raw       string
	Id        int64
	Timestamp time.Time
	Text      string
}

// Extract and return custom boop configuration file path (-c flag) from
// `os.Args`. If none provided, return default boop configuration file path
// (`~/.config/boop/confi.yaml` or `~/Library/Application Support/...` on macOS)
func configPath() string {
	for i, arg := range os.Args[1:] {
		if arg == "-c" {
			if i < len(os.Args)-2 {
				return os.Args[i+2]
			} else if i == len(os.Args)-2 {
				log.Fatal("no value specified for -c flag")
			} else {
				continue
			}
		}
	}
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("failed to get user configuration directory: ", err)
	}

	return path.Join(userConfigDir, "boop", "config.yaml")
}

func (c *config) load(path string) {
	// Configuration defaults
	c.Reversed = false
	c.Limit = 10

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("failed to read boop configuration file ", path, ": ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("failed to parse boop configuration file: ", err)
	}
}

func (c *config) getBoops() []Boop {
	resp, err := http.Get(fmt.Sprint(c.Endpoint, "?From=", url.QueryEscape(c.Phone)))
	if err != nil {
		log.Fatal("failed to initiate GET request: ", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed to read response body: ", err)
	}

	var rawBoops []string
	json.Unmarshal(body, &rawBoops)

	var boops []Boop

	for _, rawBoop := range rawBoops {
		strings.SplitN(rawBoop, " ", 2)
		extracted := regexp.MustCompile(`^(\d+)\s(\d+)\s([\s\S]*)`).FindStringSubmatch(rawBoop)

		id, err := strconv.ParseInt(extracted[1], 10, 64)
		if err != nil {
			log.Fatalf("failed to parse id %v: %v", extracted[1], err)
		}

		timestamp, err := strconv.ParseInt(extracted[2], 10, 64)
		if err != nil {
			log.Fatalf("failed to parse timestamp %v: %v", extracted[2], err)
		}

		ts := time.Unix(timestamp, 0)

		text := extracted[3]

		boops = append(boops, Boop{
			Raw:       rawBoop,
			Id:        id,
			Timestamp: ts,
			Text:      text,
		})
	}

	return boops
}

func (c *config) postBoop(text string) error {
	resp, err := http.PostForm(c.Endpoint, url.Values{
		"From": []string{c.Phone},
		"Body": []string{text},
	})
	if err != nil {
		return fmt.Errorf("failed to initiate POST request: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("boop failed (%v)", resp.StatusCode)
	}
}

// Parse values provided in `os.Args` and populate this config with them
func (c *config) parse(args []string) {
	if len(args) == 0 {
		return
	}

	flag, args := args[0], args[1:]

	switch flag {
	case "-r":
		c.Reversed = true
	case "-p":
		c.Phone = args[0]
		args = args[1:]
	case "-n":
		limit, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("failed to parse ", args[0], ": ", err)
		}
		c.Limit = limit
		args = args[1:]
	case "-e":
		c.Endpoint = args[0]
		args = args[1:]
	default:
		term := fmt.Sprint(flag, " ", strings.Join(args, " "))
		args = []string{}
		log.Printf("searching for \"%v\"\n", term)
	}

	c.parse(args)
}

// Return true if boop is being invoked with some data being piped into it
func openedWithPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func printBoops(boops []Boop) {
	prevDate, prevTime := "", ""
	for _, boop := range boops {
		date := boop.Timestamp.Format("2006 Jan 02")
		time := boop.Timestamp.Format("15:04")
		id := strconv.FormatInt(boop.Id, 10)

		if prevDate != date {
			if prevDate != "" {
				fmt.Println()
			}
			prevDate = date
			prevTime = time // Both because it's a first line
			fmt.Printf("%v @%v%*v\n", time, id, 80-len(time)-len(id)-2, date)
		}

		if prevTime != time {
			prevTime = time
			fmt.Println()
			fmt.Printf("%v @%v\n", time, id)
		}

		fmt.Println(strings.TrimSpace(boop.Text))
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func reverse(boops []Boop) []Boop {
	reversed := make([]Boop, len(boops))
	for i, boop := range boops {
		reversed[len(boops)-i-1] = boop
	}
	return reversed
}

func main() {
	c := config{}
	c.load(configPath())
	c.parse(os.Args[1:])
	if openedWithPipe() {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("failed to read stdin: ", err)
		}
		err = c.postBoop(string(bytes))
		if err != nil {
			log.Fatal("failed to post boop: ", err)
		}
	} else {
		boops := c.getBoops()
		if c.Reversed {
			printBoops(boops[0:min(len(boops), c.Limit)])
		} else {
			boops = reverse(boops)
			printBoops(boops[max(len(boops)-c.Limit, 0):])
		}
	}
}
