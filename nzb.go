// Easy NZB Parsing with support for non UTF-8 Charsets
package nzb

import (
    "bytes"
    charset "code.google.com/p/go-charset/charset"
    _ "code.google.com/p/go-charset/data"
    "encoding/xml"
    "io"
    "os"
)

// a slice of NzbFiles extended to allow sorting
type NzbFileSlice []*NzbFile

func (s NzbFileSlice) Len() int           { return len(s) }
func (s NzbFileSlice) Less(i, j int) bool { return s[i].Part < s[j].Part }
func (s NzbFileSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Nzb struct {
    Meta  map[string]string
    Files NzbFileSlice
}

// Parse in a NZB from file
func FromFile(filepath string) (*Nzb, error) {
    fdata, err := os.Open(filepath)
    if err != nil {
        return new(Nzb), err
    }
    return New(fdata)
}

// Parse NZB Data from string input
func FromString(data string) (*Nzb, error) {
    return New(bytes.NewBufferString(data))
}

// Parse NZB Data
func New(buf io.Reader) (*Nzb, error) {
    xnzb := new(xNzb)
    decoder := xml.NewDecoder(buf)
    decoder.CharsetReader = charset.NewReader
    err := decoder.Decode(xnzb)
    if err != nil {
        return new(Nzb), err
    }
    // convert to nicer format
    nzb := new(Nzb)
    // convert metadata
    nzb.Meta = make(map[string]string)
    for _, md := range xnzb.Metadata {
        nzb.Meta[md.Type] = md.Value
    }
    // copy files into (sortable) NzbFileSlice
    nzb.Files = make(NzbFileSlice, 0)
    for i, _ := range xnzb.File {
        nzb.Files = append(nzb.Files, &xnzb.File[i])
    }
    return nzb, nil
}

// used only for unmarshalling xml
type xNzb struct {
    XMLName  xml.Name   `xml:"nzb"`
    Metadata []xNzbMeta `xml:"head>meta"`
    File     []NzbFile  `xml:"file"` // xml:tag name doesn't work?
}

// used only in unmarshalling xml
type xNzbMeta struct {
    Type  string `xml:"type,attr"`
    Value string `xml:",innerxml"`
}

type NzbFile struct {
    Groups   []string     `xml:"groups>group"`
    Segments []NzbSegment `xml:"segments>segment"`
    Poster   string       `xml:"poster,attr"`
    Date     int          `xml:"date,attr"`
    Subject  string       `xml:"subject,attr"`
    Part     int
}

type NzbSegment struct {
    XMLName xml.Name `xml:"segment"`
    Bytes   int      `xml:"bytes,attr"`
    Number  int      `xml:"number,attr"`
    Id      string   `xml:",innerxml"`
}
