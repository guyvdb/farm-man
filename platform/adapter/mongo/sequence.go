package mongo

import (
	"context"
	"fmt"
	"github.com/guyvdb/farm-man/platform/model/sequence"
	"github.com/guyvdb/farm-man/platform/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SEQPADLEN int32 = 7
)

type MongoSequenceRepository struct {
	database *mongo.Database
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewMongoSequenceRepository(database *mongo.Database) repository.SequenceRepository {
	return &MongoSequenceRepository{
		database: database,
	}
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func padnum(seq int, padlen int) string {
	// convert seq to a string and pad it to padlen digits long
	v := fmt.Sprintf("%d", seq)
	len := int(padlen) - len(v)

	for i := 0; i < len; i++ {
		v = fmt.Sprintf("0%s", v)
	}
	return v
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (r *MongoSequenceRepository) CreateSequence(prefix string, padding int) error {
	var val int = 1

	seq := bson.D{
		"_id": prefix,
		"pad": padding,
		"seq": val,
	}

	_, err := r.database.Collection("sequence").InsertOne(context.Background(), seq)
	return err
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (r *MongoSequenceRepository) DeleteSequence(prefix string) error {

}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (r *MongoSequenceRepository) ResetSequence(prefix string) error {

}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (r *MongoSequenceRepository) DeleteAllSequences() error {

}

/* ------------------------------------------------------------------------
 * The sequence collection has the following properties:
 * seq{_id:"PREFIX",pad: xxx, seq:xxx}
 * Because of upsert, an error can occure. Outcomes could be:
 *  - Exactly one findAndModify() would successfully insert a new document.
 * 	- Zero or more findAndModify() methods would update the newly inserted document.
 * 	- Zero or more findAndModify() methods would fail when they attempted to insert a duplicate.
 * If the method fails due to a unique index constraint violation, retry the method. Absent a
 * delete of the document, the retry should not fail.
 * --------------------------------------------------------------------- */
func (r *MongoSequenceRepository) Next(prefix string, seperator string) sequence.Sequence {
	var result bson.M
	var padlen int32
	var ok bool
	var v interface{}
	var p interface{}
	var value int32

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"_id": prefix}
	update := bson.M{
		"$inc": bson.M{"seq": 1},
	}

	err := r.database.Collection("sequence").FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)

	if err != nil {
		// need to rerun here
		err := r.database.Collection("sequence").FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
		if err != nil {
			// no we have trouble
			return sequence.Sequence(fmt.Sprintf("%s-ERROR", prefix))
		}
	}

	p, ok = result["pad"]
	if !ok {
		padlen = SEQPADLEN
	} else {
		padlen, ok = p.(int32)
		if !ok {
			padlen = SEQPADLEN
		}
	}

	v, ok = result["seq"]
	if ok {
		value, ok = v.(int32)
		if !ok {
			return sequence.Sequence(fmt.Sprintf("%s-ERROR (seq not int32)", prefix))
		}
		return sequence.Sequence(fmt.Sprintf("%s%s%s", prefix, seperator, padnum(value, padlen)))
	} else {
		return sequence.Sequence(fmt.Sprintf("%s-ERROR", prefix))
	}

}
