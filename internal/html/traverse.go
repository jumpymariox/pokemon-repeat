package selector

import (
	"golang.org/x/net/html"
)

// traverse node attribute
func TraverseNodeAttr(node *html.Node, attrArray []string, nodeTag string, nodeAttr string) []string {
	if node.Type == html.ElementNode && node.Data == nodeTag {
		for _, img := range node.Attr {
			if img.Key == nodeAttr {
				// attrArray[len(attrArray)] = img.Val
				attrArray = append(attrArray, img.Val)
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		attrArray = TraverseNodeAttr(child, attrArray, nodeTag, nodeAttr)
	}

	return attrArray
}
