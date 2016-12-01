package main

import (
	"bufio"
	"fmt"
	"os"
	"flag"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
	"log"

	"regexp"
	"errors"
)

var input_file = flag.String("i", "in.template", "input template")
var output_file = flag.String("o", "out.file", "output template")
var zk_host = flag.String("zk", "127.0.0.1", "input template")
var system_name = flag.String("sn", "local", "system_name")

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
func process(lines []string, vars map[string]string) ([]string, error) {
	var lines_new []string
	lines_new = make([]string, len(lines))
	for index, line := range lines {
		match := regexp.MustCompile(`\{{(.*?)\}}`).FindAllStringSubmatch(line, -1)
		if match != nil {
			if _, ok := vars[match[0][1]];ok  {
				lines_new[index] = strings.Replace(line, match[0][0], vars[match[0][1]], -1)
			} else {

				return lines_new, errors.New("cannot find variable " + match[0][1] + " in /" + *system_name + "/env" )
			}
		} else {
			lines_new[index] = line
		}
	}
	err := new(error)
	return lines_new, *err
}
func get_child(conn zk.Conn, path string) map[string]string{
	defer timeTrack(time.Now(), "get_child")
	defer conn.Close()
	data, _, err := conn.Get(path)

	if err != nil {
		panic(err)
	}
	vars := make(map[string]string)
	for _, param := range strings.Split(string(data[:]),"\n") {
		kv := strings.Split(param, "=")
		if len(kv) == 2 {
			fmt.Printf("%s=%s\n", kv[0], kv[1])
			vars[kv[0]]=kv[1]
		}
	}

	return vars
}

func main() {
	flag.Parse()
	zk_hosts := strings.Split(*zk_host, ",")

	conn, _, err := zk.Connect(zk_hosts, time.Second * 5)
	if err != nil {
		panic(err)
	}
	vars := get_child(*conn, "/" + *system_name + "/env")
	template, err := readLines(*input_file)
	if err != nil {
		panic(err)
	}
	result, err := process(template,vars)
	if err != nil {
		panic(err)
	}
	err = writeLines(result, *output_file)
	if err != nil {
		panic(err)
	}
}