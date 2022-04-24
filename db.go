package brief

import (
	"fmt"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

var _dbs = map[string]*DB{}

// DB for storing nodes
type DB struct {
	bolt *bolt.DB
}

// NewDB for brief objects (aka BoltDB)
func NewDB(filepath string, mode os.FileMode, options *bolt.Options) (*DB, error) {
	boltdb, err := bolt.Open(filepath, mode, options)
	if err != nil {
		return nil, err
	}
	db := &DB{
		bolt: boltdb,
	}
	_dbs[boltdb.Path()] = db
	return db, nil
}

// Path where db is stored
func (db *DB) Path() string {
	return db.bolt.Path()
}

// PutNode into brief db
func (db *DB) PutNode(node *Node) error {
	if len(node.Name) == 0 {
		return fmt.Errorf("unable to store unnamed nodes")
	}
	return db.bolt.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(node.Type))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return bkt.Put([]byte(node.Name), node.Encode())
	})
}

// GetNode from brief db
func (db *DB) GetNode(Type, Name string) (*Node, error) {
	var node *Node
	return node, db.bolt.View(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(Type))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		out := bkt.Get([]byte(Name))
		if len(out) == 0 {
			return nil
		}
		nodes, err := Decode(strings.NewReader(string(out)), db.bolt.Path())
		if len(nodes) > 0 {
			node = nodes[0]
		}
		return err
	})
}

// Ref to a node
type Ref struct {
	Type, Name, Path string
	node             *Node
	db               *DB
}

// NewRef create reference to a node
func (db *DB) NewRef(node *Node) *Ref {
	return &Ref{
		Type: node.Type,
		Name: node.Name,
		Path: db.Path(),
		node: node,
		db:   db,
	}
}

// Node from referernce
func (ref *Ref) Node() *Node {
	if ref.node == nil {
		ref.node, _ = ref.db.GetNode(ref.Type, ref.Name)
	}
	return ref.node
}
