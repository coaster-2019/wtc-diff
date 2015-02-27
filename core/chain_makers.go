package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethutil"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/pow"
	"github.com/ethereum/go-ethereum/state"
	"math/big"
)

// So we can generate blocks easily
type FakePow struct{}

func (f FakePow) Search(block pow.Block, stop <-chan struct{}) []byte { return nil }
func (f FakePow) Verify(block pow.Block) bool                         { return true }
func (f FakePow) GetHashrate() int64                                  { return 0 }
func (f FakePow) Turbo(bool)                                          {}

func NewBlockFromParent(addr []byte, parent *types.Block) *types.Block {
	return newBlockFromParent(addr, parent)
}

func MakeBlock(bman *BlockProcessor, parent *types.Block, i int, db ethutil.Database) *types.Block {
	return makeBlock(bman, parent, i, db)
}

func MakeChain(bman *BlockProcessor, parent *types.Block, max int, db ethutil.Database) types.Blocks {
	return makeChain(bman, parent, max, db)
}

func NewChainMan(block *types.Block, eventMux *event.TypeMux, db ethutil.Database) *ChainManager {
	return newChainManager(block, eventMux, db)
}

func NewBlockProc(db ethutil.Database, txpool *TxPool, cman *ChainManager, eventMux *event.TypeMux) *BlockProcessor {
	return newBlockProcessor(db, txpool, cman, eventMux)
}

func NewCanonical(n int, db ethutil.Database) (*BlockProcessor, error) {
	return newCanonical(n, db)
}

func newBlockFromParent(addr []byte, parent *types.Block) *types.Block {
	block := types.NewBlock(parent.Hash(), addr, parent.Root(), ethutil.BigPow(2, 32), nil, "")

	block.SetUncles(nil)
	block.SetTransactions(nil)
	block.SetReceipts(nil)

	header := block.Header()
	header.Difficulty = CalcDifficulty(block, parent)
	header.Number = new(big.Int).Add(parent.Header().Number, ethutil.Big1)
	header.GasLimit = CalcGasLimit(parent, block)

	block.Td = parent.Td

	return block
}

// Actually make a block by simulating what miner would do
func makeBlock(bman *BlockProcessor, parent *types.Block, i int, db ethutil.Database) *types.Block {
	addr := ethutil.LeftPadBytes([]byte{byte(i)}, 20)
	block := newBlockFromParent(addr, parent)
	state := state.New(block.Root(), db)
	cbase := state.GetOrNewStateObject(addr)
	cbase.SetGasPool(CalcGasLimit(parent, block))
	cbase.AddBalance(BlockReward)
	state.Update(ethutil.Big0)
	block.SetRoot(state.Root())
	return block
}

// Make a chain with real blocks
// Runs ProcessWithParent to get proper state roots
func makeChain(bman *BlockProcessor, parent *types.Block, max int, db ethutil.Database) types.Blocks {
	bman.bc.currentBlock = parent
	blocks := make(types.Blocks, max)
	for i := 0; i < max; i++ {
		block := makeBlock(bman, parent, i, db)
		td, err := bman.processWithParent(block, parent)
		if err != nil {
			fmt.Println("process with parent failed", err)
			panic(err)
		}
		block.Td = td
		blocks[i] = block
		parent = block
		fmt.Printf("New Block: %x\n", block.Hash())
	}
	fmt.Println("Done making chain")
	return blocks
}

// Create a new chain manager starting from given block
// Effectively a fork factory
func newChainManager(block *types.Block, eventMux *event.TypeMux, db ethutil.Database) *ChainManager {
	bc := &ChainManager{db: db, genesisBlock: GenesisBlock(db), eventMux: eventMux}
	if block == nil {
		bc.Reset()
	} else {
		bc.currentBlock = block
		bc.td = block.Td
	}
	return bc
}

// block processor with fake pow
func newBlockProcessor(db ethutil.Database, txpool *TxPool, cman *ChainManager, eventMux *event.TypeMux) *BlockProcessor {
	bman := NewBlockProcessor(db, txpool, newChainManager(nil, eventMux, db), eventMux)
	bman.Pow = FakePow{}
	return bman
}

// Make a new canonical chain by running InsertChain
// on result of makeChain
func newCanonical(n int, db ethutil.Database) (*BlockProcessor, error) {
	eventMux := &event.TypeMux{}
	txpool := NewTxPool(eventMux)

	bman := newBlockProcessor(db, txpool, newChainManager(nil, eventMux, db), eventMux)
	bman.bc.SetProcessor(bman)
	parent := bman.bc.CurrentBlock()
	if n == 0 {
		return bman, nil
	}
	lchain := makeChain(bman, parent, n, db)
	bman.bc.InsertChain(lchain)
	return bman, nil
}
