/*
Package cmd do
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
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type provider struct {
	token string
}

type variable struct {
	resourceName string
	value        string
}

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Command to manage red team infrastructure within DigitalOcean",
	Long: `Recommended workflow for those new/unfamiliar with Terraform automation:
1. Variables - terragen do -v
2. Provider - terragen do -p
3. Droplets - terragen do -d
4. DNS Records - terragen do -n
5. Firewalls - terragen do -f

Assuming that terraform is installed some where in your PATH, terragen will run "terraform fmt" for you to keep 
the scripts organized and neat. Once the scripts are ready to deploy run:

terraform init - initializes and prepares the declared providers
terraform validate - quick check to find any potential errors in the scripts
terraform plan -out <name_of_plan> - generates a Terraform plan with a snapshot of what your provisioned infra looks like
terraform apply <name_of_plan> - provisions your infrastructure scripts in DigitalOcean 
`,
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
			fmt.Println("[x] Error: Unknown flag given. Try running `terragen do -h`")
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
	variableFile, err := os.OpenFile("variables.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer variableFile.Close()

	for {
		fmt.Printf("Enter the name of your variable resource name (enter quit when finished): ")
		var varName string
		fmt.Scanln(&varName)
		if varName == "quit" {
			break
		} else {
			fmt.Printf("Enter the value of your variable: ")
			var varValue string
			fmt.Scanln(&varValue)
			newVariable := variable{varName, varValue}
			variableFile.WriteString("variable \"" + newVariable.resourceName + "\" {\n" + "default = \"" + newVariable.value + "\"\n}\n")
		}
	}
	checkAndFormat()

}

func createProvider() {

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
	providerInfo :=
		`provider "digitalocean" {
		token = var.apikey
	}
	`
	providerFile.WriteString(providerInfo)

	fmt.Printf("Enter your API key here or press enter to have Terraform prompt you for it later (avoids hardcoding the key): ")
	var t string
	fmt.Scanln(&t)
	if len(strings.TrimSpace(t)) == 0 {
		providerFile.WriteString("variable \"apikey\" {}\n")
	} else {
		providerFile.WriteString("variable \"apikey\" {\n")
		providerFile.WriteString("default = \"" + t + "\"\n")
		providerFile.WriteString("}\n")
	}

	fmt.Printf("Enter the name of your SSH key specified in the DO SSH panel: ")
	var s string
	fmt.Scanln(&s)
	providerFile.WriteString("data \"digitalocean_ssh_key\" \"" + s + "\" {\n")
	providerFile.WriteString("name = \"" + s + "\"\n")
	providerFile.WriteString("}\n")

	checkAndFormat()
}

func createDroplet() {
	dropletFile, err := os.OpenFile("droplets.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer dropletFile.Close()

	fmt.Printf("Enter the name of the droplet resource: ")
	var resourceName string
	fmt.Scanln(&resourceName)
	dropletFile.WriteString(`resource "digitalocean_droplet" "` + resourceName + "\" {\n")

	fmt.Printf("Enter the image (OS) to deploy: ")
	var image string
	fmt.Scanln(&image)
	dropletFile.WriteString(`image = "` + image + "\"\n")

	fmt.Printf("Enter the name (hostname) of the droplet: ")
	var hostname string
	fmt.Scanln(&hostname)
	if strings.Contains(hostname, "var.") {
		dropletFile.WriteString(`name = ` + hostname + "\n")
	} else {
		dropletFile.WriteString(`name = "` + hostname + "\"\n")
	}

	fmt.Printf("Enter the region to deploy the droplet to: ")
	var region string
	fmt.Scanln(&region)
	dropletFile.WriteString(`region = "` + region + "\"\n")

	fmt.Printf("Enter the size of the droplet to deploy (hit enter to default to cheapest option): ")
	var size string
	fmt.Scanln(&size)
	if len(strings.TrimSpace(size)) == 0 {
		dropletFile.WriteString(`size = "s-1vcpu-1gb"` + "\n")
	} else {
		dropletFile.WriteString(`size = "` + size + "\"\n")
	}

	fmt.Printf("Enter the name of the SSH key to be used for authentication: ")
	var sshName string
	fmt.Scanln(&sshName)
	dropletFile.WriteString(`ssh_keys = [data.digitalocean_ssh_key.` + sshName + ".id]\n}\n")

	checkAndFormat()
}

func createDNS() {
	dnsFile, err := os.OpenFile("dns.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer dnsFile.Close()

	fmt.Printf("Enter the resource name of the record: ")
	var resourceName string
	fmt.Scanln(&resourceName)
	dnsFile.WriteString(`resource digitalocean_record "` + resourceName + `" {` + "\n")

	fmt.Printf("Enter the domain you want to apply the record to: ")
	var domain string
	fmt.Scanln(&domain)
	if strings.Contains(domain, "var.") {
		dnsFile.WriteString(`domain = ` + domain + "\n")
	} else {
		dnsFile.WriteString(`domain = "` + domain + "\"\n")
	}

	fmt.Printf("Enter the type of record (A, CNAME, etc): ")
	var dnsType string
	fmt.Scanln(&dnsType)
	dnsFile.WriteString("type = \"" + dnsType + "\"\n")

	fmt.Printf("Enter the name (hostname/subdomain) of the record: ")
	var hostname string
	fmt.Scanln(&hostname)
	dnsFile.WriteString(`name = "` + hostname + "\"\n")

	fmt.Printf("Enter the value (most likely IP address or interpolated string) of the record: ")
	var value string
	fmt.Scanln(&value)
	if dnsType == "CNAME" || dnsType == "MX" {
		dnsFile.WriteString(`value = "${` + value + `}."` + "\n")
	} else if strings.Contains(value, "digitalocean_droplet.") {
		dnsFile.WriteString(`value = ` + value + "\n")
	} else {
		dnsFile.WriteString(`value = "` + value + "\"\n")
	}

	if dnsType == "MX" {
		dnsFile.WriteString(`priority = 10` + "\n")
	}

	fmt.Printf("Enter the TTL value (pressing enter defaults the TTL to 600): ")
	var ttl string
	fmt.Scanln(&ttl)
	if len(strings.TrimSpace(ttl)) == 0 {
		dnsFile.WriteString(`ttl = 600` + "\n}\n")
	} else {
		dnsFile.WriteString("ttl = " + ttl + "\n}\n")
	}

	checkAndFormat()
}

func createFirewall() {
	fmt.Println("create firewall here")
	checkAndFormat()
}

func checkAndFormat() bool {
	// check if Terraform is in PATH
	_, err := exec.LookPath("terraform")
	if err != nil {
		fmt.Println("[x] Error: Terraform does not appear to be in PATH. Run `terraform fmt` wherever it's located")
		return false
	}
	cmd := exec.Command("terraform", "fmt")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	return true
}
