package hydrator

import (
	"testing"

	"golang.org/x/net/context"

	etcd "github.com/coreos/etcd/client"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type EtcdSuite struct {
	client etcd.Client
	kAPI   etcd.KeysAPI
}

var _ = Suite(&EtcdSuite{})

const testFolder = "testing"

var fixtures = map[string]string{
	"/testing/string":        "foo",
	"/testing/struct.string": "qux",
	"/testing/bool":          "true",
	"/testing/float32":       "42.24",
	"/testing/integer":       "42",
}

func (s *EtcdSuite) SetUpSuite(c *C) {
	cfg := etcd.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: etcd.DefaultTransport,
	}

	var err error
	s.client, err = etcd.New(cfg)
	if err != nil {
		panic(err)
	}

	s.kAPI = etcd.NewKeysAPI(s.client)

	for key, value := range fixtures {
		if _, err := s.kAPI.Set(context.Background(), key, value, nil); err != nil {
			panic(err)
		}
	}
}

func (s *EtcdSuite) TearDownSuite(c *C) {
	recursive := &etcd.DeleteOptions{
		Recursive: true,
	}

	if _, err := s.kAPI.Delete(context.Background(), testFolder, recursive); err != nil {
		panic(err)
	}
}

func (s *EtcdSuite) TestBasic(c *C) {
	foo := &Example{}

	r := NewHydrator(s.client)
	r.Folder = testFolder
	r.Hydrate(foo)

	s.assertTypes(c, foo)
}

func (s *EtcdSuite) assertTypes(c *C, value *Example) {
	c.Assert(value.String, Equals, "foo")
	c.Assert(value.Aliased, Equals, "foo")
	c.Assert(value.Bytes, DeepEquals, []byte{'f', 'o', 'o'})
	c.Assert(value.Struct.String, Equals, "qux")
	c.Assert(value.Integer, Equals, 42)
	c.Assert(value.Integer8, Equals, int8(42))
	c.Assert(value.Integer16, Equals, int16(42))
	c.Assert(value.Integer32, Equals, int32(42))
	c.Assert(value.Integer64, Equals, int64(42))
	c.Assert(value.UInteger, Equals, uint(42))
	c.Assert(value.UInteger8, Equals, uint8(42))
	c.Assert(value.UInteger16, Equals, uint16(42))
	c.Assert(value.UInteger32, Equals, uint32(42))
	c.Assert(value.UInteger64, Equals, uint64(42))
	c.Assert(value.Float32, Equals, float32(42.24))
	c.Assert(value.Float64, Equals, 42.24)
	c.Assert(value.Bool, Equals, true)
}

type Example struct {
	String     string
	Bytes      []byte `etcd:"string"`
	Aliased    string `etcd:"string"`
	Bool       bool
	Integer    int
	Integer8   int8   `etcd:"integer"`
	Integer16  int16  `etcd:"integer"`
	Integer32  int32  `etcd:"integer"`
	Integer64  int64  `etcd:"integer"`
	UInteger   uint   `etcd:"integer"`
	UInteger8  uint8  `etcd:"integer"`
	UInteger16 uint16 `etcd:"integer"`
	UInteger32 uint32 `etcd:"integer"`
	UInteger64 uint64 `etcd:"integer"`
	Float32    float32
	Float64    float64 `etcd:"float32"`

	Struct struct {
		String string
	}
}
