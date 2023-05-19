package gqlscalar

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectId primitive.ObjectID

func (x ObjectId) Hex() string {
	return primitive.ObjectID(x).Hex()
}

func (x ObjectId) Timestamp() time.Time {
	return primitive.ObjectID(x).Timestamp()
}

func (x ObjectId) String() string {
	return primitive.ObjectID(x).String()
}

func (x ObjectId) IsZero() bool {
	return primitive.ObjectID(x).IsZero()
}

func (x ObjectId) MarshalText() ([]byte, error) {
	return primitive.ObjectID(x).MarshalText()
}

func (x *ObjectId) UnmarshalText(b []byte) error {
	oid, err := primitive.ObjectIDFromHex(string(b))
	if err != nil {
		return err
	}
	*x = ObjectId(oid)
	return nil
}

func (x ObjectId) MarshalJSON() ([]byte, error) {
	t := primitive.ObjectID(x)
	return t.MarshalJSON()
}

func (x *ObjectId) UnmarshalJSON(b []byte) error {
	t := primitive.NewObjectID()
	if err := t.UnmarshalJSON(b); err != nil {
		return err
	}
	*x = ObjectId(t)
	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (x *ObjectId) UnmarshalGQL(v any) error {
	t, ok := v.(string)
	if !ok {
		return fmt.Errorf("ObjectId must be a string")
	}

	parsed, err := primitive.ObjectIDFromHex(t)
	if err != nil {
		return err
	}
	*x = ObjectId(parsed)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (x ObjectId) MarshalGQL(w io.Writer) {
	w.Write([]byte(fmt.Sprintf(`"%s"`, (primitive.ObjectID)(x).Hex())))
}

func MongoRegistry() *bsoncodec.Registry {
	rb := bson.NewRegistryBuilder()
	t := reflect.TypeOf(ObjectId{})

	rb.RegisterTypeEncoder(t, bsoncodec.ValueEncoderFunc(func(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
		if !val.IsValid() || val.Type() != t {
			return bsoncodec.ValueEncoderError{Name: "ObjectIDEncodeValue", Types: []reflect.Type{t}, Received: val}
		}
		return vw.WriteObjectID(primitive.ObjectID(val.Interface().(ObjectId)))
	}))
	rb.RegisterTypeDecoder(t, bsoncodec.ValueDecoderFunc(func(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
		if !val.CanSet() || val.Type() != t {
			return bsoncodec.ValueDecoderError{Name: "ObjectIDDecodeValue", Types: []reflect.Type{t}, Received: val}
		}

		var oid primitive.ObjectID
		var err error
		switch vrType := vr.Type(); vrType {
		case bsontype.ObjectID:
			oid, err = vr.ReadObjectID()
			if err != nil {
				return err
			}
		case bsontype.String:
			str, err := vr.ReadString()
			if err != nil {
				return err
			}
			if oid, err = primitive.ObjectIDFromHex(str); err == nil {
				break
			}
			if len(str) != 12 {
				return fmt.Errorf("an ObjectID string must be exactly 12 bytes long (got %v)", len(str))
			}
			byteArr := []byte(str)
			copy(oid[:], byteArr)
		case bsontype.Null:
			if err = vr.ReadNull(); err != nil {
				return err
			}
		case bsontype.Undefined:
			if err = vr.ReadUndefined(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot decode %v into an ObjectID", vrType)
		}

		val.Set(reflect.ValueOf(ObjectId(oid)))
		return nil
	}))

	return rb.Build()
}
