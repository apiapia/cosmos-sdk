package server

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
)

// ShowNodeIDCmd - ported from Tendermint, dump node ID to stdout
func ShowNodeIDCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "show_node_id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidator - ported from Tendermint, show this node's validator info
func ShowValidatorCmd(ctx *Context) *cobra.Command {
	flagJSON := "json"
	cmd := cobra.Command{
		Use:   "show_validator",
		Short: "Show this node's tendermint validator info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := ctx.Config
			privValidator := pvm.LoadOrGenFilePV(cfg.PrivValidatorFile())
			pubKey := privValidator.PubKey

			if viper.GetBool(flagJSON) {

				cdc := wire.NewCodec()
				wire.RegisterCrypto(cdc)
				pubKeyJSONBytes, err := cdc.MarshalJSON(pubKey)
				if err != nil {
					return err
				}
				fmt.Println(string(pubKeyJSONBytes))
				return nil
			}
			pubKeyHex := strings.ToUpper(hex.EncodeToString(pubKey.Bytes()))
			fmt.Println(pubKeyHex)
			return nil
		},
	}
	cmd.Flags().Bool(flagJSON, false, "get machine parseable output")
	return &cmd
}

// UnsafeResetAllCmd - extension of the tendermint command, resets initialization
func UnsafeResetAllCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe_reset_all",
		Short: "Reset blockchain database, priv_validator.json file, and the logger",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			tcmd.ResetAll(cfg.DBDir(), cfg.PrivValidatorFile(), ctx.Logger)
			return nil
		},
	}
}
