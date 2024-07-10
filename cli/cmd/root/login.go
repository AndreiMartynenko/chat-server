package root

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"

	desc "github.com/AndreiMartynenko/auth/pkg/auth_v1"
	"github.com/AndreiMartynenko/common/pkg/closer"
)

func login(ctx context.Context, address string, certPath string) error {
	client, err := authClient(address, certPath)
	if err != nil {
		return err
	}

	var username string

	fmt.Print(color.CyanString("Username: "))
	fmt.Scan(&username)

	fmt.Print(color.CyanString("Password (no echo): "))
	password, err := term.ReadPassword(0)
	if err != nil {
		return err
	}

	// Get refresh token
	res, err := client.Login(ctx, &desc.LoginRequest{
		Creds: &desc.Creds{
			Username: username,
			Password: string(password),
		},
	})
	if err != nil {
		return err
	}

	// Get access token
	resAccessToken, err := client.GetAccessToken(ctx, &desc.GetAccessTokenRequest{
		RefreshToken: res.GetRefreshToken(),
	})
	if err != nil {
		return err
	}

	// Save access token in file for later operations
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	closer.Add(file.Close)

	wr := bufio.NewWriter(file)
	_, err = wr.WriteString(resAccessToken.GetAccessToken())
	if err != nil {
		return err
	}

	err = wr.Flush()
	if err != nil {
		return err
	}
	return nil
}
