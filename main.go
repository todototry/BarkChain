package main

import (
	"BarkChain/blockchain"
	"fmt"

)

func main()  {
	//db, err := bolt.Open(blockchain.DB_PATH, 0600, &bolt.Options{Timeout: 1 * time.Second})
	//if err != nil {
	//	fmt.Println("open db error")
	//}
	//
	//
	//err = db.View(func(tx *bolt.Tx) error {
	//
	//	bucket := tx.Bucket([]byte(blockchain.BUCKET_NAME))
	//
	//	cur := bucket.Cursor()
	//
	//	k,v := cur.First()
	//
     //   fmt.Println(k, v)
	//
	//return nil
	//})
	//
	//if err != nil {
	//	fmt.Println("read db view.error ")
	//}
	//os.Exit(1)


	bc := blockchain.NewBlockChain()
	bc.NewBlock("send 3 Yuan to Eva")
	bc.NewBlock("send 2 Yuan to Fan")

	// view
	for index, b := range bc.Blocks {
		fmt.Printf("%d, %x , %x,  %s\n",index, string(b.BlockHeader.PrevHash), string(b.BlockHeader.CurHash), b.Data)
	}

	// validate
	for index, b := range bc.Blocks {
		proof := blockchain.NewProofOfWork(b, blockchain.POW_TARGET)
		fmt.Printf("%d, %x , %x,  %s  , %v \n",index, string(b.BlockHeader.PrevHash), string(b.BlockHeader.CurHash), b.Data, proof.Validate() )
	}

	bc.Db.Close()

}
