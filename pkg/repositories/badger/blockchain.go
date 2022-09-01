package badger

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/tclaudel/golang_blockchain/config"
	entity2 "github.com/tclaudel/golang_blockchain/internal/entity"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/entity"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
)

const LastHashKey = "lh"

type PublicKey struct {
	PkX string `json:"x"`
	PkY string `json:"y"`
}

type Signature struct {
	SigR string `json:"r"`
	SigS string `json:"s"`
}

type Transaction struct {
	TSenderPublicKey  PublicKey `json:"sender_public_key"`
	TSenderAddress    string    `json:"sender_address"`
	TRecipientAddress string    `json:"recipient_address"`
	TAmount           float64   `json:"amount"`
	TSignature        Signature `json:"signature"`
}

func (t Transaction) SenderPublicKey() (values.PublicKey, error) {
	pk, err := values.PublicKeyFromStrings(t.TSenderPublicKey.PkX, t.TSenderPublicKey.PkY)
	if err != nil {
		return values.PublicKey{}, err
	}

	return pk, nil
}

func (t Transaction) SenderAddress() values.Address {
	return values.AddressFromString(t.TSenderAddress)
}

func (t Transaction) RecipientAddress() values.Address {
	return values.AddressFromString(t.TRecipientAddress)
}

func (t Transaction) Amount() values.Amount {
	return values.AmountFromFloat64(t.TAmount)
}

func (t Transaction) Signature() (values.Signature, error) {
	sig, err := values.SignatureFromStrings(t.TSignature.SigR, t.TSignature.SigS)
	if err != nil {
		return values.Signature{}, err
	}

	return sig, nil
}

func (t Transaction) Entity() entity.Transaction {
	return t
}

type Block struct {
	BTimestamp    time.Time     `json:"timestamp"`
	BNonce        int           `json:"nonce"`
	BPreviousHash []byte        `json:"previous_hash"`
	BTransactions []Transaction `json:"transactions"`
	BHash         []byte        `json:"hash"`
}

func (b *Block) Hash() values.Hash {
	return values.HashFromBytes(b.BHash)
}

func (b *Block) Timestamp() values.Timestamp {
	return values.TimestampFromTime(b.BTimestamp)
}

func (b *Block) Nonce() values.Nonce {
	return values.NonceFromInt(b.BNonce)
}

func (b *Block) PreviousHash() values.Hash {
	return values.HashFromBytes(b.BPreviousHash)
}

func (b Block) Transactions() ([]values.Transaction, error) {
	var transactions = make([]values.Transaction, len(b.BTransactions))
	for i, transaction := range b.BTransactions {
		pk, err := transaction.SenderPublicKey()
		if err != nil {
			return nil, err
		}

		sig, err := transaction.Signature()
		if err != nil {
			return nil, err
		}

		transactions[i] = values.TransactionFromValues(
			pk,
			transaction.SenderAddress(),
			transaction.RecipientAddress(),
			transaction.Amount(),
			sig,
		)
	}

	return transactions, nil
}

func (b *Block) Entity() entity.Block {
	return b
}

type Blockchain struct {
	chain []*Block
}

func (b Blockchain) Chain() []entity.Block {
	var chain = make([]entity.Block, len(b.chain))
	for i, block := range b.chain {
		chain[i] = block.Entity()
	}

	return chain
}

func blockchainFromEntity(bc entity.Blockchain) (Blockchain, error) {
	var blockchain Blockchain
	for _, block := range bc.Chain() {
		var transactions []Transaction
		txs, err := block.Transactions()
		if err != nil {
			return Blockchain{}, err
		}
		for _, transaction := range txs {
			pk, err := transaction.SenderPublicKey()
			if err != nil {
				return Blockchain{}, err
			}

			sig, err := transaction.Signature()
			if err != nil {
				return Blockchain{}, err
			}

			x, y := pk.Strings()
			r, s := sig.Strings()
			transactions = append(transactions, Transaction{
				TSenderPublicKey: PublicKey{
					PkX: x,
					PkY: y,
				},
				TSenderAddress:    transaction.SenderAddress().String(),
				TRecipientAddress: transaction.RecipientAddress().String(),
				TAmount:           transaction.Amount().Float64(),
				TSignature: Signature{
					SigR: r,
					SigS: s,
				},
			})
		}

		blockchain.chain = append(blockchain.chain, &Block{
			BTimestamp:    block.Timestamp().Time(),
			BNonce:        block.Nonce().Int(),
			BPreviousHash: block.PreviousHash().Bytes(),
			BTransactions: transactions,
			BHash:         block.Hash().Bytes(),
		})
	}

	return Blockchain{
		chain: blockchain.chain,
	}, nil
}

type BlockchainRepository struct {
	logger *zap.Logger
	badger *badger.DB
}

func (b BlockchainRepository) Append(entityBlock entity.Block) error {
	txs, err := entityBlock.Transactions()
	if err != nil {
		b.logger.Error("failed to get BTransactions from block", zap.Error(err))
		return err
	}

	var block = Block{
		BTimestamp:    entityBlock.Timestamp().Time(),
		BNonce:        entityBlock.Nonce().Int(),
		BPreviousHash: entityBlock.PreviousHash().Bytes(),
		BTransactions: func() []Transaction {
			var transactions = make([]Transaction, len(txs))
			for i, transaction := range txs {
				pk, err := transaction.SenderPublicKey()
				if err != nil {
					b.logger.Error("failed to get sender public key from transaction", zap.Error(err))
					return nil
				}

				sig, err := transaction.Signature()
				if err != nil {
					b.logger.Error("failed to get TSignature from transaction", zap.Error(err))
					return nil
				}

				pkX, pkY := pk.Strings()
				sigR, sigS := sig.Strings()
				transactions[i] = Transaction{
					TSenderPublicKey: PublicKey{
						PkX: pkX,
						PkY: pkY,
					},
					TSenderAddress:    transaction.SenderAddress().String(),
					TRecipientAddress: transaction.RecipientAddress().String(),
					TAmount:           transaction.Amount().Float64(),
					TSignature: Signature{
						SigR: sigR,
						SigS: sigS,
					},
				}
			}

			return transactions
		}(),
		BHash: entityBlock.Hash().Bytes(),
	}

	err = b.badger.Update(func(txn *badger.Txn) error {
		enc, err := json.Marshal(block)
		if err != nil {
			b.logger.Error("failed to serialize block", zap.Error(err))
			return err
		}

		b.logger.Info("appending block", zap.String("hash", block.Hash().String()), zap.String("previous hash", entityBlock.PreviousHash().String()))
		err = txn.Set(block.Hash().Bytes(), enc)
		if err != nil {
			b.logger.Error("failed to set block", zap.Error(err))
			return err
		}

		b.logger.Debug("setting last hash", zap.Any(LastHashKey, block.Hash()))
		err = txn.Set([]byte(LastHashKey), block.Hash().Bytes())
		if err != nil {
			b.logger.Error("failed to set last BHash", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil

}

func (b BlockchainRepository) Get() (entity.Blockchain, error) {
	var blockchain = new(Blockchain)

	iterator, err := b.newIterator()
	if err != nil {
		return nil, err
	}

	var i = 0
	for {
		block, err := iterator.next()
		if err != nil {
			return nil, err
		}

		blockchain.chain = append(blockchain.chain, block)
		if block.Hash().Equal(values.GenesisHash()) {
			break
		}
		if i > 3 {
			break
		}
		i++
	}

	bc, err := blockchainFromEntity(*blockchain)
	if err != nil {
		b.logger.Error("cannot convert entity.Blockchain to Blockchain", zap.Error(err))
		return nil, err
	}
	return bc, nil
}

type blockchainIterator struct {
	currentHash values.Hash
	badger      *badger.DB
	logger      *zap.Logger
}

func (b BlockchainRepository) newIterator() (*blockchainIterator, error) {
	var lastHash values.Hash
	err := b.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LastHashKey))
		if err != nil {
			b.logger.Error("failed to get last hash", zap.Error(err))
			return err
		}

		err = item.Value(func(val []byte) error {
			lastHash = values.HashFromBytes(val)
			return nil
		})
		if err != nil {
			b.logger.Error("failed to get last hash", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &blockchainIterator{
		badger:      b.badger,
		logger:      b.logger,
		currentHash: lastHash,
	}, nil
}

func (b *blockchainIterator) next() (*Block, error) {
	var nextBlock = new(Block)
	err := b.badger.View(func(txn *badger.Txn) error {
		b.logger.Debug("getting block", zap.String("hash", b.currentHash.String()))
		item, err := txn.Get(b.currentHash.Bytes())
		if err != nil {
			b.logger.Error("failed to get value from block", zap.String("key", string(b.currentHash.Bytes())), zap.Error(err))
			return err
		}

		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, nextBlock)
			if err != nil {
				b.logger.Error("failed to deserialize block", zap.Error(err))
				return err
			}
			return nil
		})

		if err != nil {
			b.logger.Error("failed to get last BHash value", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	b.currentHash = values.HashFromBytes(nextBlock.BPreviousHash)
	return nextBlock, nil
}

func (b BlockchainRepository) Close() error {
	return b.badger.Close()
}

func NewBlockchainRepository(logger *zap.Logger, blockchainCfg config.Blockchain) (repositories.Blockchain, error) {
	opts := badger.DefaultOptions(blockchainCfg.Badger.Path)
	if blockchainCfg.Badger.Logger {
		opts.Logger = NewBadgerLogger(logger)
	} else {
		opts.Logger = nil
	}
	db, err := badger.Open(opts)
	if err != nil {
		logger.Error("failed to open badger", zap.Error(err))
		return nil, err
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LastHashKey)); err == badger.ErrKeyNotFound {
			logger.Info("no last Hash found, creating genesis block")
			genesisBlock := Block{
				BTimestamp:    entity2.Genesis.Timestamp().Time(),
				BNonce:        entity2.Genesis.Nonce().Int(),
				BPreviousHash: entity2.Genesis.PreviousHash().Bytes(),
				BTransactions: nil,
				BHash:         values.GenesisHash().Bytes(),
			}
			data, err := json.Marshal(genesisBlock)
			if err != nil {
				logger.Error("failed to serialize genesis block", zap.Error(err))
				return err
			}

			logger.Info("appending genesis block", zap.Any(values.GenesisHash().String(), genesisBlock))
			err = txn.Set(values.GenesisHash().Bytes(), data)
			if err != nil {
				return err
			}

			logger.Debug("setting genesis last hash", zap.Any(LastHashKey, values.GenesisHash().String()))
			err = txn.Set([]byte(LastHashKey), values.GenesisHash().Bytes())
			if err != nil {
				return err
			}
		} else {
			item, err := txn.Get([]byte(LastHashKey))
			if err != nil {
				return err
			}
			err = item.Value(func(val []byte) error {
				return nil
			})

			return err

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	bcRepo := &BlockchainRepository{
		logger: logger,
		badger: db,
	}

	return bcRepo, nil
}
