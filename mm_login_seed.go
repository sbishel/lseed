package main

import (
	"bytes"
	dbsql "database/sql"
	"encoding/base32"
	"flag"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pborman/uuid"

	"fmt"
	"log"
)

// var ouDN, bindUser, bindPassword, bindHost, photoPath string
// var bindPort, numGroups, numMembersPerGroup int
// var help bool
var numMembers int
var help bool

func main() {

	flag.IntVar(&numMembers, "members", 100000, "the number of members to create")

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	db, err := dbsql.Open("mysql", "mmuser:mostest@tcp(127.0.0.1:3306)/mattermost_test?charset=utf8mb4,utf8\u0026readTimeout=30s\u0026writeTimeout=30s")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// var id string

	insertStatement := "INSERT INTO Users (" +
		"Id, CreateAt, UpdateAt, DeleteAt, Username," +
		"Password, AuthData, AuthService, Email, EmailVerified," +
		"Nickname, Firstname, Lastname, Position, Roles," +
		"AllowMarketing, Props, " +
		"NotifyProps," +
		"LastPasswordUpdate, LastPictureUpdate, FailedAttempts, Locale," +
		"Timezone," +
		"MfaActive, MfaSecret" +
		") VALUES (" +
		"?, ?, ?, 0, ?," +
		"'', ?, 'ldap', ?, true," +
		"?, ?, ?, 805306368, 'system_user'," +
		"false, '{ }', " +
		"'{ \"channel\" :\"true\", \"comments\" :\"never\", \"desktop\" :\"mention\", \"desktop_sound\" :\"true\", \"email\" :\"true\", \"first_name\" :\"false\", \"mention_keys\" :\"\", \"push\" :\"mention\", \"push_status\" :\"away\" }'," +
		"?, 0, 0, 'en'," +
		"'{ \"automaticTimezone\" :\"America/Toronto\", \"manualTimezone\" :\"\", \"useAutomaticTimezone\" :\"true\" }'," +
		"false, ''" +
		");"

	fmt.Print(insertStatement)

	insert, err := db.Prepare(insertStatement)
	if err != nil {
		log.Fatal(err)
	}

	// create users
	for i := 10000; i < (numMembers); i++ {

		var id = NewId()
		var createdAt = time.Now().Unix()
		var userName = fmt.Sprintf("test.%d", i)
		var email = fmt.Sprintf("success+test%d@simulator.amazonses.com", i)
		var firstName = fmt.Sprintf("test.%d", i)
		var lastName = "User"

		_, err := insert.Exec(id, createdAt, createdAt, userName, userName, email, userName, firstName, lastName, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		// _, err = db.Exec(insertStatement, params)
		fmt.Print(".")
	}
	fmt.Println("\nSuccessfully completed.")

}

func NewId() string {
	var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}

// flag.StringVar(&ouDN, "ou", "ou=loadtest,dc=mm,dc=test,dc=com", "the organizational unit that will contain the seeded data")
// flag.StringVar(&bindUser, "user", "cn=admin,dc=mm,dc=test,dc=com", "the bind user")
// flag.StringVar(&bindPassword, "password", "mostest", "the bind password")
// flag.StringVar(&bindHost, "host", "0.0.0.0", "the bind host")
// flag.IntVar(&bindPort, "port", 389, "the bind port")
// flag.IntVar(&numGroups, "groups", 2, "the number of groups")
// flag.IntVar(&numMembersPerGroup, "members", 10, "the number of members per group")
// flag.StringVar(&photoPath, "photo", "", "the path to the profile photo")
// flag.BoolVar(&help, "help", false, "show help")
// flag.Parse()

// l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", bindHost, bindPort))
// if err != nil {
// 	log.Fatal(err)
// }
// defer l.Close()

// fmt.Printf(bindUser)
// fmt.Printf(bindPassword)
// err = l.Bind(bindUser, bindPassword)
// if err != nil {
// 	log.Fatal(err)
// }

// // create org. unit
// err = l.Add(&ldap.AddRequest{
// 	DN: ouDN,
// 	Attributes: []ldap.Attribute{
// 		{Type: "objectclass", Vals: []string{"organizationalunit"}},
// 	},
// })
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf(".")

// // get profile photo data
// var strData string
// if len(photoPath) > 0 {
// 	imageData, err := ioutil.ReadFile(photoPath)
// 	if err != nil {
// 		log.Fatal(err)
// 		os.Exit(1)
// 	}
// 	strData = string(imageData)
// }

// // create users
// for i := 0; i < (numGroups * numMembersPerGroup); i++ {
// 	attributes := []ldap.Attribute{
// 		{Type: "objectclass", Vals: []string{"iNetOrgPerson"}},
// 		{Type: "cn", Vals: []string{fmt.Sprintf("Test%d", i)}},
// 		{Type: "sn", Vals: []string{"User"}},
// 		{Type: "mail", Vals: []string{fmt.Sprintf("success+test%d@simulator.amazonses.com", i)}},
// 		{Type: "userPassword", Vals: []string{"Password1"}},
// 	}
// 	if len(strData) > 0 {
// 		attributes = append(attributes, ldap.Attribute{Type: "jpegPhoto", Vals: []string{strData}})
// 	}
// 	err = l.Add(&ldap.AddRequest{
// 		DN:         fmt.Sprintf("uid=test.%d,%s", i, ouDN),
// 		Attributes: attributes,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf(".")
// }

// for i := 0; i < numGroups; i++ {
// 	groupDN := fmt.Sprintf("cn=tgroup-%d,%s", i, ouDN)

// 	var uniqueMembers []string
// 	for j := 0; j < numMembersPerGroup; j++ {
// 		uniqueMembers = append(uniqueMembers, fmt.Sprintf("uid=test.%d,%s", j+(numMembersPerGroup*i), ouDN))
// 	}

// 	err = l.Add(&ldap.AddRequest{
// 		DN: groupDN,
// 		Attributes: []ldap.Attribute{
// 			{Type: "objectclass", Vals: []string{"groupOfUniqueNames"}},
// 			{Type: "uniqueMember", Vals: []string{uniqueMembers[0]}},
// 		},
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf(".")

// 	for _, member := range uniqueMembers[1:] {
// 		err = l.Modify(&ldap.ModifyRequest{
// 			DN: groupDN,
// 			Changes: []ldap.Change{
// 				{Operation: ldap.AddAttribute, Modification: ldap.PartialAttribute{Type: "uniqueMember", Vals: []string{member}}},
// 			},
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf(".")
// 	}
// }

// }

// func queryDatabase() {

// 	db, err := sql.Open("mysql", "mmuser@tcp(127.0.0.1:3306)/mattermost_test")

// 	result, err := db.ExecContext(ctx,
// 		"INSERT INTO users (name, age) VALUES ($1, $2)",
// 		"gopher",
// 		27,
// 	)
// }
