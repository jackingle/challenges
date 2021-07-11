package main

import (
	"fmt"
	"net"
	"time"
)

// Write a program that will check to see if port 443, 80, and 22 are open on these sites.   Sites = [‘google.com’, ‘aol.com’, ‘aws.net’, ‘facebook.com’].
// Display the result as formatted json written to the console.

var sites = []string{"google.com", "aol.com", "aws.net", "facebook.com"}
var ports = []int{443, 80, 22}

type TotalResult struct {
	Websites []Website `json:"websites"`
}
type Website struct {
	HostName string `json:"host_name"`
	Ports    []Port `json:"ports"`
}
type Port struct {
	PortNumber int    `json:"port_number"`
	Response   string `json:"response"`
}

func main() {
	result := checkPorts()
	fmt.Println(result)
}

func checkPorts() TotalResult {
	result := &TotalResult{}

	for i := range sites {
		site := Website{
			HostName: sites[i],
		}
		result.Websites = append(result.Websites, site)

		for j := range ports {
			port := Port{
				PortNumber: ports[j],
			}
			result.Websites[i].Ports = append(result.Websites[i].Ports, port)
			// result.Websites[i].Ports[j].PortNumber = ports[j]
			if err := checkPort(sites[i], ports[j]); err != nil {
				result.Websites[i].Ports[j].Response = "Failed to connect!"
			} else {
				result.Websites[i].Ports[j].Response = "Successfully connected!"
			}
		}
	}
	return *result
}

func checkPort(site string, port int) error {
	host := fmt.Sprintf("%v:%v", site, port)
	conn, err := net.DialTimeout("tcp", host, time.Duration(3*time.Second))
	if err != nil {
		return err
	}

	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		fmt.Printf("Timeout error: %s\n", err)
		return err
	}
	defer conn.Close()

	return nil
}
