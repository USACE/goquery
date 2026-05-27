package goquery

type FluentInsert struct {
	store      DataStore
	ds         DataSet
	batch      bool
	batchSize  int
	batchQueue Batch
	tx         *Tx
	records    interface{}
	panicOnErr bool
	//return id
}

const defaultBatchSize = 100

func (i *FluentInsert) Tx(tx *Tx) *FluentInsert {
	i.tx = tx
	return i
}

func (i *FluentInsert) Batch(batch bool) *FluentInsert {
	i.batch = batch
	return i
}

func (i *FluentInsert) BatchSize(bs int) *FluentInsert {
	i.batchSize = bs
	return i
}

func (i *FluentInsert) BatchQueue(bq Batch) *FluentInsert {
	i.batchQueue = bq
	return i
}

func (i *FluentInsert) Records(recs interface{}) *FluentInsert {
	i.records = recs
	return i
}

func (i *FluentInsert) PanicOnErr(panicOnErr bool) *FluentInsert {
	i.panicOnErr = panicOnErr
	return i
}

func (i *FluentInsert) Execute() error {
	ii := InsertInput{
		Dataset:    i.ds,
		Records:    i.records,
		Batch:      i.batch,
		BatchSize:  i.batchSize,
		BatchQueue: i.batchQueue,
		PanicOnErr: i.panicOnErr,
	}
	return i.store.InsertRecs(i.tx, ii)
}
