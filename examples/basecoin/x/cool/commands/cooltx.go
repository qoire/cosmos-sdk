package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/cosmos/cosmos-sdk/examples/basecoin/x/cool"
)

// what cool transaction
func WhatCoolTxCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "whatcool [answer]",
		Short: "What's cooler than being cool?",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide an answer")
			}

			name := viper.GetString(client.FlagName)
			if name == "" {
				return errors.Errorf("must provide a name using --name")
			}

			// get the from address from the name flag
			from, err := builder.GetFromAddress(name)
			if err != nil {
				return err
			}

			// create the message
			msg := cool.NewWhatCoolMsg(from, args[0])

			// get password
			buf := client.BufferStdin()
			prompt := fmt.Sprintf("Password to sign with '%s':", name)
			passphrase, err := client.GetPassword(prompt, buf)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			res, err := builder.SignBuildBroadcast(name, passphrase, msg, cdc)
			if err != nil {
				return err
			}

			fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			return nil
		},
	}
}

// set what cool transaction
func SetWhatCoolTxCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "setwhatcool [answer]",
		Short: "You're so cool, tell us what is cool!",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide an answer")
			}

			name := viper.GetString(client.FlagName)
			if name == "" {
				return errors.Errorf("must provide a name using --name")
			}

			// get the from address from the name flag
			from, err := builder.GetFromAddress(name)
			if err != nil {
				return err
			}

			// get password
			buf := client.BufferStdin()
			prompt := fmt.Sprintf("Password to sign with '%s':", name)
			passphrase, err := client.GetPassword(prompt, buf)
			if err != nil {
				return err
			}

			// create the message
			msg := cool.NewSetWhatCoolMsg(from, args[0])

			// build and sign the transaction, then broadcast to Tendermint
			res, err := builder.SignBuildBroadcast(name, passphrase, msg, cdc)
			if err != nil {
				return err
			}

			fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			return nil
		},
	}
}
