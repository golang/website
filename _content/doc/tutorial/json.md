<!--{
  "Title": "Tutorial: Working with JSON",
  "Breadcrumb": true
}-->

JSON (JavaScript Object Notation) is a simple data interchange format.
Syntactically it resembles the objects and lists of JavaScript.
It is commonly used for communication with networked API services,
but it is used in many other places, too.
Its home page, [json.org](http://json.org),
provides a wonderfully clear and concise definition of the standard.

With the [encoding/json/v2 package](/pkg/encoding/json/v2) it's a snap to read
and write JSON data from your Go programs.
This package provides a cleaner API and better defaults than the older
[encoding/json](/pkg/encoding/json) package.

## Encoding

To encode JSON data we use the [`Marshal`](/pkg/encoding/json/v2#Marshal) function.

```
func Marshal(in any, opts ...Options) (out []byte, err error)
```

Given the Go data structure, `Message`,

```
type Message struct {
	Name string
	Body string
	Time time.Time
}
```

and an instance of `Message`

```
m := Message{"Alice", "Hello", time.Date(2011, 1, 25, 0, 0, 0, 0, time.UTC)}
```

we can marshal a JSON-encoded version of m using `json.Marshal`:

```
b, err := json.Marshal(m)
```

If all is well, `err` will be `nil` and `b` will be a `[]byte` containing this JSON data:

```
b == []byte(`{"Name":"Alice","Body":"Hello","Time":"2011-01-25T00:00:00Z"}`)
```

Only data structures that can be represented as valid JSON will be encoded:

  - Pointers will be encoded as the values they point to (or `null` if the pointer is `nil`).

  - JSON objects only support strings as keys.
    Thus, to encode a Go map, it must have a key that encodes as a JSON string,
    such as `map[string]T` (where `T` is any Go type supported by the json
    package).

  - Channel, complex, and function types cannot be encoded.

  - Cyclic data structures are not supported.

  - [`Marshal`](/pkg/encoding/json/v2#Marshal) documents the full set of encoding semantics.

The json package only accesses the exported fields of struct types (those
that begin with an uppercase letter).
Therefore only the exported fields of a struct will be present in the JSON output.

## Decoding

To decode JSON data we use the [`Unmarshal`](/pkg/encoding/json/v2#Unmarshal) function.

```
func Unmarshal(in []byte, out any, opts ...Options) (err error)
```

We must first create a place where the decoded data will be stored

```
var m Message
```

and call `json.Unmarshal`, passing it a `[]byte` of JSON data and a pointer to `m`

```
err := json.Unmarshal(b, &m)
```

If `b` contains valid JSON that fits in `m`,
after the call `err` will be `nil` and the data from `b` will have been
stored in the struct `m`,
as if by an assignment like:

```
m = Message{
	Name: "Alice",
	Body: "Hello",
	Time: time.Date(2011, 1, 25, 0, 0, 0, 0, time.UTC),
}
```

How does `Unmarshal` identify the fields in which to store the decoded data?
For a given JSON key `"Foo"`,
`Unmarshal` will look through the destination struct's fields to find (in
order of preference):

  - An exported field with a tag of `json:"Foo"` (see the [Go spec](/ref/spec#Struct_types)
    for more on struct tags), or

  - An exported field named `Foo`.

What happens when the structure of the JSON data doesn't exactly match the Go type?

```
b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
var m Message
err := json.Unmarshal(b, &m)
```

`Unmarshal` will decode only the fields that it can find in the destination type.
In this case, only the `Name` field of `m` will be populated,
and the `Food` field will be ignored.
This behavior is particularly useful when you wish to pick only a few specific
fields out of a large JSON blob.
It also means that any unexported fields in the destination struct will
be unaffected by `Unmarshal`.

But what if you don't know the structure of your JSON data beforehand?

## Generic JSON with `any`

The `any` type describes any Go type.
`any` is defined as `interface{}` (empty interface), which describes an interface with zero methods.
Every Go type implements at least zero methods and therefore satisfies the empty interface.

`any` serves as a general container type:

	var a any
	a = "a string"
	a = 2011
	a = 2.777

A type assertion accesses the underlying concrete type:

	r := a.(float64)
	fmt.Println("the circle's area", math.Pi*r*r)

Or, if the underlying type is unknown, a type switch determines the type:

	switch v := a.(type) {
	case int:
	    fmt.Println("twice a is", v*2)
	case float64:
	    fmt.Println("the reciprocal of a is", 1/v)
	case string:
	    h := len(v) / 2
	    fmt.Println("a swapped by halves is", v[h:]+v[:h])
	default:
	    // a isn't one of the types above
	}

The json package uses `map[string]any` and
`[]any` values to store arbitrary JSON objects and arrays;
it will happily unmarshal any valid JSON blob into a plain
`any` value.  The default concrete Go types are:

  - `bool` for JSON booleans,

  - `float64` for JSON numbers,

  - `string` for JSON strings, and

  - `nil` for JSON null.

## Decoding arbitrary data

Consider this JSON data, stored in the variable `b`:

```
b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
```

Without knowing this data's structure, we can decode it into an `any` value with `Unmarshal`:

```
var f any
err := json.Unmarshal(b, &f)
```

At this point the Go value in `f` would be a map whose keys are strings
and whose values are themselves stored as `any` values:

```
f = map[string]any{
	"Name": "Wednesday",
	"Age":  6,
	"Parents": []any{
		"Gomez",
		"Morticia",
	},
}
```

To access this data we can use a type assertion to access `f`'s underlying `map[string]any`:

```
m := f.(map[string]any)
```

We can then iterate through the map with a range statement and use a type
switch to access its values as their concrete types:

```
for k, v := range m {
	switch vv := v.(type) {
	case string:
		fmt.Println(k, "is string", vv)
	case float64:
		fmt.Println(k, "is float64", vv)
	case []any:
		fmt.Println(k, "is an array:")
		for i, u := range vv {
			fmt.Println(i, u)
		}
	default:
		fmt.Println(k, "is of a type I don't know how to handle")
	}
}
```

In this way you can work with unknown JSON data while still enjoying the benefits of type safety.

## Reference Types

Let's define a Go type to contain the data from the previous example:

```
type FamilyMember struct {
	Name    string
	Age     int
	Parents []string
}

var m FamilyMember
err := json.Unmarshal(b, &m)
```

Unmarshaling that data into a `FamilyMember` value works as expected,
but if we look closely we can see a remarkable thing has happened.
With the var statement we allocated a `FamilyMember` struct,
and then provided a pointer to that value to `Unmarshal`,
but at that time the `Parents` field was a `nil` slice value.
To populate the `Parents` field, `Unmarshal` allocated a new slice behind the scenes.
This is typical of how `Unmarshal` works with the supported reference types
(pointers, slices, and maps).

Consider unmarshaling into this data structure:

```
type Foo struct {
	Bar *Bar
}
```

If there were a `Bar` field in the JSON object,
`Unmarshal` would allocate a new `Bar` and populate it.
If not, `Bar` would be left as a `nil` pointer.

From this a useful pattern arises: if you have an application that receives
a few distinct message types,
you might define "receiver" structure like

```
type IncomingMessage struct {
	Cmd *Command
	Msg *Message
}
```

and the sending party can populate the `Cmd` field and/or the `Msg` field
of the top-level JSON object,
depending on the type of message they want to communicate.
`Unmarshal`, when decoding the JSON into an `IncomingMessage` struct,
will only allocate the data structures present in the JSON data.
To know which messages to process, the programmer need simply test that
either `Cmd` or `Msg` is not `nil`.

## Streaming Marshal and Unmarshal

The [`io.Reader`](/pkg/io#Reader) and [`io.Writer`](/pkg/io#Writer) interfaces are ubiquitous in Go, providing streaming access to resources such as HTTP connections, WebSockets, or files.
The `MarshalWrite` and `UnmarshalRead` functions allow marshaling and unmarshaling directly to and from these streams without requiring an intermediate `[]byte` holding the entire message.

```
func MarshalWrite(out io.Writer, in any, opts ...Options) (err error)
func UnmarshalRead(in io.Reader, out any, opts ...Options) (err error)
```

For example, we can write marshal output directly to standard out:

```
err := json.MarshalWrite(os.Stdout, m)
```

## Custom Marshal and Unmarshal

Sometimes the default marshal behavior does not make sense for your type.

For example, suppose we have a type describing a software version.

```
type Version struct {
	Major, Minor, Patch int64
}

v := Version{1, 2, 3}
b, err := json.Marshal(v)
```

This marshals as `{"Major":1,"Minor":2,"Patch":3}`, as expected based on the struct definition, but a version string like "1.2.3" would likely be a better representation for JSON consumers.

By implementing the [`MarshalerTo`](/pkg/encoding/json/v2#MarshalerTo) interface, we can provide a custom JSON representation.
[`UnmarshalerFrom`](/pkg/encoding/json/v2#UnmarshalerFrom) provides unmarshaling of the custom JSON representation.

```
type MarshalerTo interface {
	MarshalJSONTo(*jsontext.Encoder) error
}

type UnmarshalerFrom interface {
	UnmarshalJSONFrom(*jsontext.Decoder) error
}
```

We can implement these method to convert to and from string representations.

```
func (v Version) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
}
```

With custom marshaling, the `Version{1, 2, 3}` value now marshals as `"1.2.3"`.

```
func (v *Version) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k != jsontext.KindString {
		// Value must be a string.
		return &json.SemanticError{JSONKind: k}
	}

	var s string
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return err
	}

	_, err := fmt.Sscanf(s, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	return err
}
```

With custom unmarshaling, `"1.2.3"` now unmarshals as `Version{1, 2, 3}`.
This simple example requires the exact `"major.minor.patch"` version format, but a more complex `UnmarshalJSONFrom` could choose to add flexibility, such as making the minor and patch versions optional.

## References

For more information see the [encoding/json/v2](/pkg/encoding/json/v2) and [encoding/json/jsontext](/pkg/encoding/json/jsontext) package documentation.
