package main

import "bufio"
import "fmt"
import "os"
import "strings"
import "strconv"

type Mode int
const (
	Access Mode = iota
	Trunk
)

func configPort(intfName string, opts *Opts) (uint8, Mode) {
	if *opts.config == "" {
		return 1, Access
	}


	file, err := os.Open(*opts.config)
	if err != nil {
		fmt.Println(err)
		return 1, Access
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		// should check input errors
		if line[0] != intfName {
			continue
		}
		var vlan int
		var mode Mode
		switch line[1] {
		case "trunk":
			vlan = 1
			mode = Trunk
		case "access":
			vlan, err = strconv.Atoi(line[2])
			if err != nil {
				fmt.Println(err)
				break
			}
			mode = Access
		}
		return uint8(vlan), mode

	}
	if err := scanner.Err(); err!= nil {
		fmt.Println(err)
		return 1, Access
	}
	return 1, Access
}
	

