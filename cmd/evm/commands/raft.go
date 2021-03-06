package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Fantom-foundation/go-evm/src/consensus/raft"
	"github.com/Fantom-foundation/go-evm/src/engine"
)

//AddRaftFlags adds flags to the Raft command
func AddRaftFlags(cmd *cobra.Command) {

	cmd.Flags().String("raft.dir", config.Raft.RaftDir, "Base directory for Raft data")
	cmd.Flags().String("raft.snapshot-dir", config.Raft.SnapshotDir, "Snapshot directory")
	cmd.Flags().String("raft.node-addr", config.Raft.NodeAddr, "IP:PORT of Raft node")
	cmd.Flags().String("raft.server-id", string(config.Raft.LocalID), "Unique ID of this server")

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic("Unable to bind viper flags")
	}
}

//NewRaftCmd returns the command that starts EVM-Lite with Raft consensus
func NewRaftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raft",
		Short: "Run the evm node with Raft consensus",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			config.SetDataDir(config.BaseConfig.DataDir)

			logger.WithFields(logrus.Fields{
				"Raft": config.Raft,
			}).Debug("Config")

			return nil
		},
		RunE: runRaft,
	}

	AddRaftFlags(cmd)

	return cmd
}

func runRaft(cmd *cobra.Command, args []string) error {
	raftConsensus := raft.NewRaft(*config.Raft, logger)
	consensusEngine, err := engine.NewConsensusEngine(*config, raftConsensus, logger)
	if err != nil {
		return fmt.Errorf("error building Engine: %s", err)
	}

	if err := consensusEngine.Run(); err != nil {
		return fmt.Errorf("error running Engine: %s", err)
	}

	return nil
}
