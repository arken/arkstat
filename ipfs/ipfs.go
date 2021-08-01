package ipfs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ipfsConfig "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi" // This package is needed so that all the preloaded plugins are loaded automatically.
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	migrate "github.com/ipfs/go-ipfs/repo/fsrepo/migrations"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type NodeConfArgs struct {
	PeerID         string
	PrivKey        string
	SwarmKey       string
	BootstrapPeers []string
}

type Node struct {
	api    icore.CoreAPI
	ctx    context.Context
	cancel context.CancelFunc
	node   *core.IpfsNode
}

func (n *Node) SetHandler(protocolID string, handler network.StreamHandler) {
	n.node.PeerHost.SetStreamHandler(protocol.ID(protocolID), handler)
}

// CreateNode creates an IPFS node and returns its coreAPI
func CreateNode(repoPath string, args NodeConfArgs) (node *Node, err error) {
	// Initialize Node ID
	id := ipfsConfig.Identity{
		PeerID:  args.PeerID,
		PrivKey: args.PrivKey,
	}
	// Initialize node structure
	node = &Node{}
	// Create IPFS node
	node.ctx, node.cancel = context.WithCancel(context.Background())
	// Create Swarm Key File
	if args.SwarmKey != "" {
		err = createSwarmKey(repoPath, args.SwarmKey)
		if err != nil {
			return nil, err
		}
	}
	// Open the repo
	fs, err := openFs(node.ctx, repoPath)
	if err != nil {
		err = createFs(node.ctx, repoPath, id, args.BootstrapPeers)
		if err != nil {
			return nil, err
		}
	}
	// Construct the node
	nodeOptions := &core.BuildCfg{
		Permanent: true,
		Online:    true,
		Routing:   libp2p.DHTOption,
		Repo:      fs,
	}
	node.node, err = core.NewNode(node.ctx, nodeOptions)
	if err != nil {
		return nil, err
	}
	node.node.IsDaemon = true
	// Attach the Core API to the constructed node
	node.api, err = coreapi.NewCoreAPI(node.node)
	return node, err

}

func openFs(ctx context.Context, repoPath string) (result repo.Repo, err error) {
	result, err = fsrepo.Open(repoPath)
	if err != nil && err == fsrepo.ErrNeedMigration {
		err = os.Setenv("IPFS_PATH", repoPath)
		if err != nil {
			return nil, err
		}
		err = migrate.RunMigration(ctx, migrate.NewHttpFetcher("", "", "ipfs", 0), fsrepo.RepoVersion, repoPath, false)
		if err != nil {
			return nil, err
		}
		result, err = fsrepo.Open(repoPath)
	}
	return result, err
}

// createFs builds the IPFS configuration repository.
func createFs(ctx context.Context, path string, id ipfsConfig.Identity, bootstrapPeers []string) (err error) {
	// Check if directory to configuration exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	// Create a ipfsConfig with default options and a 2048 bit key
	cfg, err := ipfsConfig.InitWithIdentity(id)
	if err != nil {
		return err
	}
	// Set default ipfsConfig values
	cfg.Datastore.StorageMax = "5GB"
	cfg.Reprovider.Strategy = "roots"
	cfg.Routing.Type = "dhtserver"
	cfg.Bootstrap = bootstrapPeers
	cfg.Swarm.ConnMgr.HighWater = 1200
	cfg.Swarm.ConnMgr.LowWater = 1000

	// Create the repo with the ipfsConfig
	err = fsrepo.Init(path, cfg)
	if err != nil {
		return fmt.Errorf("failed to init node: %s", err)
	}
	return nil
}

func createSwarmKey(path string, key string) (err error) {
	keyPath := filepath.Join(path, "swarm.key")
	// Check if directory to configuration exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	if _, err = os.Stat(keyPath); os.IsNotExist(err) {
		var file *os.File
		file, err = os.Create(keyPath)
		if err != nil {
			return err
		}
		_, err = file.WriteString("/key/swarm/psk/1.0.0/\n/base16/\n" + key)
	}
	return err
}
