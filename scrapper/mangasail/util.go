package mangasail

import (
	"context"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func getCrawledHtmlDocument(ctx context.Context, url string, waitJSPathSel string, wantedJSPathSel string) (*goquery.Document, error) {
	var htmlContent string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(waitJSPathSel, chromedp.ByJSPath),
		chromedp.OuterHTML(wantedJSPathSel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		return nil, err
	}
	if htmlContent == "" {
		return nil, errors.New("empty html content")
	}

	return goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
}
