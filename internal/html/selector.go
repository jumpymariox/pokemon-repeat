package selector

import (
	"golang.org/x/net/html"
)

// traverse img
func TraverseNodeAttr(node *html.Node, attrArray []string, nodeTag string, nodeAttr string) []string {
	if node.Type == html.ElementNode && node.Data == "img" {
		for _, img := range node.Attr {
			if img.Key == "src" {
				attrArray = append(attrArray, img.Val)
				// fmt.Println(img.Val)
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		attrArray = TraverseNodeAttr(child, attrArray, nodeTag, nodeAttr)
	}

	return attrArray
}
