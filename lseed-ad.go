package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-ldap/ldap"
)

var ouDN, bindUser, bindPassword, bindHost, photoPath string
var bindPort, numGroups, numMembersPerGroup int
var help bool

func main() {
	flag.StringVar(&ouDN, "ou", "OU=loadtest,DC=example,DC=com", "the organizational unit that will contain the seeded data")
	flag.StringVar(&bindUser, "user", "example\\administrator", "the bind user")
	flag.StringVar(&bindPassword, "password", "3HVqjA5auLew7u3&be=%2rQathqjnsHi", "the bind password")
	flag.StringVar(&bindHost, "host", "example.com.dev.spinmint.com", "the bind host")
	flag.IntVar(&bindPort, "port", 389, "the bind port")
	flag.IntVar(&numGroups, "groups", 2, "the number of groups")
	flag.IntVar(&numMembersPerGroup, "members", 10, "the number of members per group")
	flag.StringVar(&photoPath, "photo", "", "the path to the profile photo")
	flag.BoolVar(&help, "help", false, "show help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", bindHost, bindPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// fmt.Printf(bindUser)
	// fmt.Printf(bindPassword)
	err = l.Bind(bindUser, bindPassword)
	if err != nil {
		log.Fatal(err)
	}

	// create org. unit
	err = l.Add(&ldap.AddRequest{
		DN: ouDN,
		Attributes: []ldap.Attribute{
			{Type: "objectclass", Vals: []string{"organizationalunit"}},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(".")

	// get profile photo data
	var strData string
	if len(photoPath) > 0 {
		imageData, err := ioutil.ReadFile(photoPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		strData = string(imageData)
	}

	// create users
	fmt.Printf("\n Creating Users")
	for i := 0; i < (numGroups * numMembersPerGroup); i++ {
		attributes := []ldap.Attribute{
			{Type: "objectclass", Vals: []string{"iNetOrgPerson"}},
			{Type: "cn", Vals: []string{fmt.Sprintf("test.%d", i)}},
			{Type: "sn", Vals: []string{"User"}},
			{Type: "mail", Vals: []string{fmt.Sprintf("success+test%d@simulator.amazonses.com", i)}},
			{Type: "userPassword", Vals: []string{"Password1"}},
		}
		if len(strData) > 0 {
			attributes = append(attributes, ldap.Attribute{Type: "thumbnailPhoto", Vals: []string{strData}})
		}

		dn := fmt.Sprintf("CN=Test.%d,%s", i, ouDN)
		_, err := ldap.ParseDN(dn)

		if err != nil {
			log.Fatal(err)
		}

		err = l.Add(&ldap.AddRequest{
			DN:         fmt.Sprintf("cn=test.%d,%s", i, ouDN),
			Attributes: attributes,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(".")
		if i%1000 == 0 {
			fmt.Printf("%v", i)
			fmt.Printf("%v", time.Now())
		}

	}

	fmt.Printf("\n Creating Groups")
	for i := 0; i < numGroups; i++ {
		groupDN := fmt.Sprintf("cn=tgroup-%d,%s", i, ouDN)

		var uniqueMembers []string
		for j := 0; j < numMembersPerGroup; j++ {
			uniqueMembers = append(uniqueMembers, fmt.Sprintf("cn=test.%d,%s", j+(numMembersPerGroup*i), ouDN))
		}

		err = l.Add(&ldap.AddRequest{
			DN: groupDN,
			Attributes: []ldap.Attribute{
				{Type: "objectclass", Vals: []string{"groupOfUniqueNames"}},
				{Type: "uniqueMember", Vals: []string{uniqueMembers[0]}},
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(".")

		for _, member := range uniqueMembers[1:] {
			err = l.Modify(&ldap.ModifyRequest{
				DN: groupDN,
				Changes: []ldap.Change{
					{Operation: ldap.AddAttribute, Modification: ldap.PartialAttribute{Type: "uniqueMember", Vals: []string{member}}},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(".")
		}
	}

	fmt.Println("\nSuccessfully completed.")
}
