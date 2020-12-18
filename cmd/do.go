/*
Copyright © 2020 b1gbroth3r <b1gbroth3r@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type provider struct {
	token string
}

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		varStatus, _ := cmd.Flags().GetBool("vars")
		providerStatus, _ := cmd.Flags().GetBool("provider")
		dropletStatus, _ := cmd.Flags().GetBool("droplet")
		dnsStatus, _ := cmd.Flags().GetBool("dns")
		firewallStatus, _ := cmd.Flags().GetBool("firewall")

		if varStatus {
			createVars()
		} else if providerStatus {
			createProvider()
		} else if dropletStatus {
			createDroplet()
		} else if dnsStatus {
			createDNS()
		} else if firewallStatus {
			createFirewall()
		} else {
			fmt.Println("[x] Error: Unknown flag given.")
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	doCmd.Flags().BoolP("vars", "v", false, "Declare variables to be referenced across Terraform scripts")
	doCmd.Flags().BoolP("provider", "p", false, "Creates Digital Ocean provider")
	doCmd.Flags().BoolP("droplet", "d", false, "Declares droplets to provision")
	doCmd.Flags().BoolP("dns", "n", false, "Declares DNS records associated with infrastructure")
	doCmd.Flags().BoolP("firewall", "f", false, "Declares firewall configurations if necessary for droplets")
}

func createVars() {
	fmt.Println("create variables here")
}

func createProvider() {
	fmt.Println("create provider here")
	providerFile, err := os.OpenFile("provider.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer providerFile.Close()

	reqProviderInfo :=
		`terraform {
			required_providers {
				digitalocean = {
					source = "digitalocean/digitalocean"
					version = "2.3.0"
				}
			}
		}
	
		`
	providerFile.WriteString(reqProviderInfo)

	fmt.Printf("Enter your API key here or press enter to have Terraform prompt you for it later (avoids hardcoding the key): ")
	var t string
	fmt.Scanln(&t)
	if len(strings.TrimSpace(t)) == 0 {
		providerFile.WriteString("variable \"apikey\" {}\n")
	}

	providerInfo :=
		`provider "digitalocean" {
		token = var.apikey
		}
		
		`
	providerFile.WriteString(providerInfo)
}

func createDroplet() {
	fmt.Println("create Droplets here")
}

func createDNS() {
	fmt.Println("create DNS records here")
}

func createFirewall() {
	fmt.Println("create firewall here")
}
