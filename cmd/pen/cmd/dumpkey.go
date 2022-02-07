package cmd

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/penguintop/penguin_bsc/pkg/crypto"
	"github.com/penguintop/penguin_bsc/pkg/keystore"
	filekeystore "github.com/penguintop/penguin_bsc/pkg/keystore/file"
	memkeystore "github.com/penguintop/penguin_bsc/pkg/keystore/mem"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (c *command) initDumpKeyCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "dumpkey",
		Short: "Dump Penguin Private Wif Key",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return cmd.Help()
			}

			v := strings.ToLower(c.config.GetString(optionNameVerbosity))
			logger, err := newLogger(cmd, v)
			if err != nil {
				return fmt.Errorf("new logger: %v", err)
			}

			var keystore keystore.Service
			if c.config.GetString(optionNameDataDir) == "" {
				keystore = memkeystore.New()
				logger.Warning("A new keystore has been created in the memory, since the data directory with key store is not provided or missing.")
				return nil
			} else {
				keystore = filekeystore.New(filepath.Join(c.config.GetString(optionNameDataDir), "keys"))
			}

			var password string
			if p := c.config.GetString(optionNamePassword); p != "" {
				password = p
			} else if pf := c.config.GetString(optionNamePasswordFile); pf != "" {
				b, err := ioutil.ReadFile(pf)
				if err != nil {
					return err
				}
				password = string(bytes.Trim(b, "\n"))
			} else {
				exists, err := keystore.Exists("penguin")
				if err != nil {
					return err
				}
				if exists {
					password, err = terminalPromptPassword(cmd, c.passwordReader, "Password")
					if err != nil {
						return err
					}

					penguinPrivateKey, _, err := keystore.Key("penguin", password)
					if err != nil {
						return fmt.Errorf("penguin key: %w", err)
					}

					tempBytes := penguinPrivateKey.D.Bytes()
					var privKeyBytes [32]byte
					copy(privKeyBytes[32-len(tempBytes):], tempBytes)

					publicKey := &penguinPrivateKey.PublicKey

					signer := crypto.NewDefaultSigner(penguinPrivateKey)
					overlayEthAddress, err := signer.EthereumAddress()
					if err != nil {
						return err
					}

					address, err := crypto.NewOverlayAddress(*publicKey, c.config.GetUint64(optionNameNetworkID))
					if err != nil {
						return err
					}

					logger.Info("********************************************************************")
					logger.Infof("!!! PrivateKey: %s !!!", hex.EncodeToString(privKeyBytes[:]))
					logger.Infof("!!! Penguin Account Address: %s !!!", overlayEthAddress.String())
					logger.Infof("!!! Penguin Node Address: %s !!!", address.String())
					logger.Infof("!!! Please backup your PrivateKey, and do not tell it to anyone else !!!")
					logger.Info("********************************************************************")

				} else {
					return errors.New("Penguin private key file is missing or not existing, you may set it with the option of '--data-dir'.")
				}
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	c.setAllFlags(cmd)
	c.root.AddCommand(cmd)
	return nil
}
