package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/bytom/node"
	"github.com/bytom/types"
	log "github.com/sirupsen/logrus"
	cmn "github.com/tendermint/tmlibs/common"
)

var runNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Run the bytomd",
	RunE:  runNode,
}

func init() {
	// p2p flags
	runNodeCmd.Flags().String("p2p.laddr", config.P2P.ListenAddress, "Node listen address. (0.0.0.0:0 means any interface, any port)")
	runNodeCmd.Flags().String("p2p.seeds", config.P2P.Seeds, "Comma delimited host:port seed nodes")
	runNodeCmd.Flags().Bool("p2p.skip_upnp", config.P2P.SkipUPNP, "Skip UPNP configuration")
	runNodeCmd.Flags().Bool("p2p.pex", config.P2P.PexReactor, "Enable Peer-Exchange ")
	runNodeCmd.Flags().Bool("wallet.enable", config.Wallet.Enable, "Enable wallet")

	RootCmd.AddCommand(runNodeCmd)
}

func runNode(cmd *cobra.Command, args []string) error {
	genDocFile := config.GenesisFile()
	if cmn.FileExists(genDocFile) {
		jsonBlob, err := ioutil.ReadFile(genDocFile)
		if err != nil {
			return fmt.Errorf("Couldn't read GenesisDoc file: %v", err)
		}
		genDoc, err := types.GenesisDocFromJSON(jsonBlob)
		if err != nil {
			return fmt.Errorf("Error reading GenesisDoc: %v", err)
		}
		if genDoc.ChainID == "" {
			return fmt.Errorf("Genesis doc %v must include non-empty chain_id", genDocFile)
		}
		config.ChainID = genDoc.ChainID
		config.PrivateKey = genDoc.PrivateKey
		config.Time = genDoc.GenesisTime
	} else {
		return fmt.Errorf("not find genesis.json")
	}

	// Create & start node
	n := node.NewNodeDefault(config)
	if _, err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	} else {
		log.WithField("nodeInfo", n.Switch().NodeInfo()).Info("Started node")
	}

	// Trap signal, run forever.
	n.RunForever()

	return nil
}
