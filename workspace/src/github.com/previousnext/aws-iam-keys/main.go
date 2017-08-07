package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliRegion = kingpin.Flag("region", "AWS Region.").Default("ap-southeast-2").String()
	cliGroup  = kingpin.Flag("iam-group", "IAM group to load user list from.").Default("SSH").String()
	cliFile   = kingpin.Flag("file", "Authorized keys file to write to.").Required().String()
	cliOwner  = kingpin.Flag("owner", "Enforce this owner").Required().String()
)

func main() {
	kingpin.Parse()

	svc := iam.New(session.New(&aws.Config{Region: cliRegion}))

	var authorized []string

	list, err := svc.GetGroup(&iam.GetGroupInput{
		GroupName: cliGroup,
	})
	if err != nil {
		log.Println("failed to get group:", err)
		os.Exit(1)
	}

	for _, user := range list.Users {
		keys, err := svc.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{
			UserName: user.UserName,
		})
		if err != nil {
			log.WithFields(log.Fields{"user": *user.UserName}).Info(err)
			continue
		}

		for _, key := range keys.SSHPublicKeys {
			pub, err := svc.GetSSHPublicKey(&iam.GetSSHPublicKeyInput{
				Encoding:       aws.String(iam.EncodingTypeSsh),
				SSHPublicKeyId: key.SSHPublicKeyId,
				UserName:       user.UserName,
			})
			if err != nil {
				log.WithFields(log.Fields{"user": *user.UserName}).Info(err)
				continue
			}

			authorized = append(authorized, *pub.SSHPublicKey.SSHPublicKeyBody)
		}
	}

	err = ioutil.WriteFile(*cliFile, []byte(strings.Join(authorized, "\n")), 0600)
	if err != nil {
		log.Println("failed to write authorized file:", err)
		os.Exit(1)
	}

	log.Println("file has been written:", *cliFile)

	userData, err := user.Lookup(*cliOwner)
	if err != nil {
		log.Println("failed to lookup users uid/gid:", err)
		os.Exit(1)
	}

	uid, err := strconv.Atoi(userData.Uid)
	if err != nil {
		log.Println("failed to marshall users uid:", err)
		os.Exit(1)
	}

	gid, err := strconv.Atoi(userData.Gid)
	if err != nil {
		log.Println("failed to marshall users gid:", err)
		os.Exit(1)
	}

	err = os.Chown(*cliFile, uid, gid)
	if err != nil {
		log.Println("failed to chown authorized file:", err)
		os.Exit(1)
	}

	log.Println("user permissions have been updated")
}
