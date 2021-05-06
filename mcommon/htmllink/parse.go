package htmllink

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linknodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func linknodes(n *html.Node) []*html.Node {
	// if n.Type == html.ElementNode && n.Data == "a" {

	// 	return []*html.Node{n}
	// }
	// if n.Type == html.ElementNode && n.Data == "li" {

	// 	return []*html.Node{n}
	// }
	if n.Type == html.ElementNode && n.Data == "strong" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linknodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
		}
	}
	ret.Text = findText(n)
	return ret
}
func findText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += findText(c)
	}
	ret = strings.Join(strings.Fields(ret), " ")
	return ret
}
