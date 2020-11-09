// Package sitemap provides tools for creating XML sitemaps
// and sitemap indexes and writing them to io.Writer (such as
// http.ResponseWriter).
//
// Please see http://www.sitemaps.org/ for description of sitemap contents.
package sitemap

import (
	"encoding/xml"
	"github.com/snabb/diagio"
	"io"
	"time"
)

// ChangeFreq specifies change frequency of a sitemap entry. It is just a string.
type ChangeFreq string

// Feel free to use these constants for ChangeFreq (or you can just supply
// a string directly).
const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

type GalleryLocation struct {
	Title string `xml:"title,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Video struct {
	ThumbnailLoc string           `xml:"video:thumbnail_loc,omitempty"`
	Title        string           `xml:"video:title,omitempty"`
	Description  string           `xml:"video:description,omitempty"`
	ContentLoc   string           `xml:"video:content_loc,omitempty"`
	PlayerLoc    string           `xml:"video:player_loc,omitempty"`
	AllowEmbed   string           `xml:"video:allow_embed,omitempty"`
	Duration     int64            `xml:"video:duration,omitempty"`
	GalleryLoc   *GalleryLocation `xml:"gallery_loc,omitempty"`
}

// URL entry in sitemap or sitemap index. LastMod is a pointer
// to time.Time because omitempty does not work otherwise. Loc is the
// only mandatory item. ChangeFreq and Priority must be left empty when
// using with a sitemap index.
type URL struct {
	Loc        string     `xml:"loc"`
	LastMod    *time.Time `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"`
	Priority   float32    `xml:"priority,omitempty"`
	Video      *Video     `xml:"video:video,omitempty"`
}

// Sitemap represents a complete sitemap which can be marshaled to XML.
// New instances must be created with New() in order to set the xmlns
// attribute correctly. Minify can be set to make the output less human
// readable.
type Sitemap struct {
	XMLName    xml.Name `xml:"urlset"`
	Xmlns      string   `xml:"xmlns,attr"`
	XmlnsVideo string   `xml:"xmlns:video,attr"`

	URLs []*URL `xml:"url"`

	Minify bool `xml:"-"`
}

// New returns a new Sitemap.
func New() *Sitemap {
	return &Sitemap{
		Xmlns:      "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsVideo: "http://www.google.com/schemas/sitemap-video/1.1",
		URLs:       make([]*URL, 0),
	}
}

// Add adds an URL to a Sitemap.
func (s *Sitemap) Add(u *URL) {
	s.URLs = append(s.URLs, u)
}

// WriteTo writes XML encoded sitemap to given io.Writer.
// Implements io.WriterTo.
func (s *Sitemap) WriteTo(w io.Writer) (n int64, err error) {
	cw := diagio.NewCounterWriter(w)

	_, err = cw.Write([]byte(xml.Header))
	if err != nil {
		return cw.Count(), err
	}
	en := xml.NewEncoder(cw)
	if !s.Minify {
		en.Indent("", "  ")
	}
	err = en.Encode(s)
	cw.Write([]byte{'\n'})
	return cw.Count(), err
}

var _ io.WriterTo = (*Sitemap)(nil)

// ReadFrom reads and parses an XML encoded sitemap from io.Reader.
// Implements io.ReaderFrom.
func (s *Sitemap) ReadFrom(r io.Reader) (n int64, err error) {
	de := xml.NewDecoder(r)
	err = de.Decode(s)
	return de.InputOffset(), err
}

var _ io.ReaderFrom = (*Sitemap)(nil)
