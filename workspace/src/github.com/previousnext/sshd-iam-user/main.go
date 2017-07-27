package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliUser    = kingpin.Flag("user", "Raw user passed in by SSHD.").Required().String()
	cliAllowed = kingpin.Flag("allowed", "Comma separated list of users allow return of authorized_keys.").Required().String()
	cliRegion  = kingpin.Flag("region", "AWS Region.").Default("ap-southeast-2").String()
	cliGroup   = kingpin.Flag("group", "IAM group to load user list from.").Default("SSH").String()
)

func main() {
	kingpin.Parse()

	if !contains(*cliUser, strings.Split(*cliAllowed, ",")) {
		return
	}

	svc := iam.New(session.New(&aws.Config{Region: cliRegion}))

	list, err := svc.GetGroup(&iam.GetGroupInput{
		GroupName: cliGroup,
	})
	if err != nil {
		panic(err)
	}

	for _, user := range list.Users {
		keys, err := svc.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{
			UserName: user.UserName,
		})
		if err != nil {
			continue
		}

		for _, key := range keys.SSHPublicKeys {
			pub, err := svc.GetSSHPublicKey(&iam.GetSSHPublicKeyInput{
				Encoding:       aws.String(iam.EncodingTypeSsh),
				SSHPublicKeyId: key.SSHPublicKeyId,
				UserName:       user.UserName,
			})
			if err != nil {
				continue
			}

			fmt.Println(*pub.SSHPublicKey.SSHPublicKeyBody)
		}
	}
}
