package mangasail

import (
	"context"
	"errors"
	"github.com/bigscreen/manga-scrapper/config"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/common"
	"github.com/chromedp/chromedp"
)

func getCrawledHtmlDocument(url string, waitJSPathSel string, wantedJSPathSel string) (*goquery.Document, error) {
	ctx, cancel := getChromeDpContext()
	defer cancel()

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

func getChromeDpContext() (context.Context, context.CancelFunc) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug usage
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Chrome/73.0.3683.103`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf))
	chromeCtx, cancel = context.WithTimeout(chromeCtx, config.ChromeDPTimeout())
	return chromeCtx, cancel
}

func buildPageURL(path string) string {
	return common.MangasailBaseURL + path
}

func getIdFromPath(path string) string {
	return strings.TrimPrefix(path, common.MangasailPrefixPath)
}
