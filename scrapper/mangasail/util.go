package mangasail

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/config"
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
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent("Chrome/77.0.3830.0"),
		chromedp.WindowSize(1920, 1080),
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
