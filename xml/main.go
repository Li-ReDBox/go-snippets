// https://play.golang.org/p/d9BkGclp-1 provided by a SO question
// This is similar to the example CustomMarshalXML
// The purpose is to has a generic way to traverse an XML

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
)

// This is a nested struct
type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

/*
Unmarshaler is the interface implemented by objects that can unmarshal an XML element description of themselves.

UnmarshalXML decodes a single XML element beginning with the given start element. If it returns an error,
the outer call to Unmarshal stops and returns that error. UnmarshalXML must consume exactly one XML element.
One common implementation strategy is to unmarshal into a separate value with a layout matching the expected
XML using d.DecodeElement, and then to copy the data from that value into the receiver. Another common
strategy is to use d.Token to process the XML object one token at a time. UnmarshalXML may not use d.RawToken.

If Unmarshal encounters a field type that implements the Unmarshaler interface,
Unmarshal calls its UnmarshalXML method to produce the value from the XML element.
Otherwise, if the value implements encoding.TextUnmarshaler, Unmarshal calls that
value's UnmarshalText method.
*/
/* Get the attributes of the current node */
func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	// Without this new type, main panics!
	// runtime: goroutine stack exceeds 1000000000-byte limit
	// fatal error: stack overflow
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func describe(n Node) {
	fmt.Printf("Current node: %s\n", n)
	if n.XMLName.Local == "p" {
		fmt.Println(string(n.Content))
		fmt.Println(n.Attrs)
	}
}

func walk(nodes []Node, f func(Node)) {
	fmt.Printf("In walk func, the length of current nodes = %d\n", len(nodes))
	for _, n := range nodes {
		f(n)
		walk(n.Nodes, f)
	}
}

func main() {
	data := `
	<content>
		<p class="foo">this is content area</p>
		<animal>
			<p>This id dog</p>
			<dog>
			   <p>tommy</p>
			</dog>
		</animal>
		<birds>
			<p class="bar">this is birds</p>
			<p>this is birds</p>
		</birds>
		<animal>
			<p>this is animals</p>
		</animal>
	</content>`

	buf := bytes.NewBuffer([]byte(data))
	dec := xml.NewDecoder(buf)

	var n Node
	err := dec.Decode(&n)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Decode is done, the raw result is:")
	fmt.Printf("%s\n", n)

	// Need to convert first to get the so called length

	walk([]Node{n}, describe)
}
