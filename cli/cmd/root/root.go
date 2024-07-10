package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	filename = ".access_token"

	tokenHeader = "Bearer "

	lineUp = "\033[1A"

	addressFlag  = "address"
	certPathFlag = "cert"
	idFlag       = "id"
	usersFlag    = "users"

	addressFlagShort  = "a"
	certPathFlagShort = "c"
	idFlagShort       = "n"
	usersFlagShort    = "u"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "chat-client",
	Short: "Chat client",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to chat service",
	Run: func(cmd *cobra.Command, _ []string) {
		addr, err := cmd.Flags().GetString(addressFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", addressFlag, err)
		}

		certPath, err := cmd.Flags().GetString(certPathFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", certPathFlag, err)
		}

		err = login(context.Background(), addr, certPath)
		if err != nil {
			log.Printf("failed to login: %v", err)
		} else {
			fmt.Println(color.GreenString("\n\n[Successfully logged in]\n"))
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from chat service",
	Run: func(_ *cobra.Command, _ []string) {
		err := logout()
		if err != nil {
			log.Printf("failed to logout: %v", err)
		} else {
			fmt.Println(color.GreenString("\n\n[Successfully logged out]\n"))
		}
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create object",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete object",
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to object",
}

var createChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Create new chat",
	Run: func(cmd *cobra.Command, _ []string) {
		addr, err := cmd.Flags().GetString(addressFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", addressFlag, err)
		}

		certPath, err := cmd.Flags().GetString(certPathFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", certPathFlag, err)
		}

		users, err := cmd.Flags().GetString(usersFlag)
		if err != nil {
			log.Fatalf("failed to get users: %v", err)
		}

		id, err := createChat(context.Background(), addr, certPath, strings.Split(users, ","))
		if err != nil {
			log.Printf("failed to create chat: %v", err)
		} else {
			fmt.Printf("[%s %s]\n", color.CyanString("Created chat with id"), color.YellowString(id))
		}
	},
}

var deleteChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Delete chat",
	Run: func(cmd *cobra.Command, _ []string) {
		addr, err := cmd.Flags().GetString(addressFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", addressFlag, err)
		}

		certPath, err := cmd.Flags().GetString(certPathFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", certPathFlag, err)
		}

		id, err := cmd.Flags().GetString(idFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", idFlag, err)
		}

		err = deleteChat(context.Background(), addr, certPath, id)
		if err != nil {
			log.Printf("failed to delete chat: %v", err)
		} else {
			fmt.Printf("[%s %s]\n", color.CyanString("Deleted chat (if existed)"), color.YellowString(id))
		}
	},
}

var connectChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Connect to chat",
	Run: func(cmd *cobra.Command, _ []string) {
		addr, err := cmd.Flags().GetString(addressFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", addressFlag, err)
		}

		certPath, err := cmd.Flags().GetString(certPathFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", certPathFlag, err)
		}

		id, err := cmd.Flags().GetString(idFlag)
		if err != nil {
			log.Fatalf("failed to get %s: %v", idFlag, err)
		}

		err = connectChat(context.Background(), addr, certPath, id)
		if err != nil {
			log.Printf("failed to connect to chat: %v", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(connectCmd)

	createCmd.AddCommand(createChatCmd)
	deleteCmd.AddCommand(deleteChatCmd)
	connectCmd.AddCommand(connectChatCmd)

	// Add flags to commands
	loginCmd.Flags().StringP(addressFlag, addressFlagShort, "", "`IP:port` of authentication service")
	loginCmd.Flags().StringP(certPathFlag, certPathFlagShort, "", "path to authentication service certificate")
	createChatCmd.Flags().StringP(usersFlag, usersFlagShort, "", "chat participants")
	createChatCmd.Flags().StringP(addressFlag, addressFlagShort, "", "`IP:port` of chat service")
	createChatCmd.Flags().StringP(certPathFlag, certPathFlagShort, "", "path to chat service certificate")
	deleteChatCmd.Flags().StringP(idFlag, idFlagShort, "", "ID of the chat to delete")
	deleteChatCmd.Flags().StringP(addressFlag, addressFlagShort, "", "`IP:port` of chat service")
	deleteChatCmd.Flags().StringP(certPathFlag, certPathFlagShort, "", "path to chat service certificate")
	connectChatCmd.Flags().StringP(idFlag, idFlagShort, "", "ID of the chat to connect")
	connectChatCmd.Flags().StringP(addressFlag, addressFlagShort, "", "`IP:port` of chat service")
	connectChatCmd.Flags().StringP(certPathFlag, certPathFlagShort, "", "path to chat service certificate")

	// Mark required flags in commands
	err := loginCmd.MarkFlagRequired(addressFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", addressFlag, err)
	}

	err = loginCmd.MarkFlagRequired(certPathFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", certPathFlag, err)
	}

	err = createChatCmd.MarkFlagRequired(addressFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", idFlag, err)
	}

	err = createChatCmd.MarkFlagRequired(certPathFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", certPathFlag, err)
	}

	err = deleteChatCmd.MarkFlagRequired(idFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", idFlag, err)
	}

	err = deleteChatCmd.MarkFlagRequired(addressFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", addressFlag, err)
	}

	err = deleteChatCmd.MarkFlagRequired(certPathFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", certPathFlag, err)
	}

	err = connectChatCmd.MarkFlagRequired(idFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", idFlag, err)
	}

	err = connectChatCmd.MarkFlagRequired(addressFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", addressFlag, err)
	}

	err = connectChatCmd.MarkFlagRequired(certPathFlag)
	if err != nil {
		log.Fatalf("failed to mark %s flag as required: %v", certPathFlag, err)
	}
}
