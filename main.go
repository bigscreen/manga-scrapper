package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bigscreen/manga-scrapper/scrapper/mangasail"
	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("Hello world")
	chromeCtx, cancel := initChromeDpContext()
	defer cancel()
	start := time.Now()
	fetchMangasailContent(chromeCtx)
	fmt.Printf("Took: %f seconds\n", time.Since(start).Seconds())
}

func initChromeDpContext() (context.Context, context.CancelFunc) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug usage
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Chrome/73.0.3683.103`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf))
	chromeCtx, cancel = context.WithTimeout(chromeCtx, 20*time.Second)
	return chromeCtx, cancel
}

func fetchMangasailContent(chromeCtx context.Context) {
	fetchMangasailHomeContent(chromeCtx)
	fmt.Println("---------------------")
	fetchMangasailMangaDetailsContent(chromeCtx)
	fmt.Println("---------------------")
	fetchMangasailChapterDetailsContent(chromeCtx)
}

func fetchMangasailHomeContent(chromeCtx context.Context) {
	homeScrapper := mangasail.NewHomePageScrapper(chromeCtx)
	content, err := homeScrapper.GetContent()
	if err != nil {
		fmt.Println("fetchMangasailHomeContent, error:", err)
		return
	}
	fmt.Println("fetchMangasailHomeContent, result:", content.ToString())
}

func fetchMangasailMangaDetailsContent(chromeCtx context.Context) {
	detailsScrapper := mangasail.NewMangaDetailsPageScrapper(chromeCtx)
	content, err := detailsScrapper.GetContent("/content/checkmate-manga")
	if err != nil {
		fmt.Println("fetchMangasailMangaDetailsContent, error:", err)
		return
	}
	fmt.Println("fetchMangasailMangaDetailsContent, result:", content.ToString())
}

func fetchMangasailChapterDetailsContent(chromeCtx context.Context) {
	detailsScrapper := mangasail.NewChapterDetailsPageScrapper(chromeCtx)
	content, err := detailsScrapper.GetContent("/content/checkmate-35")
	if err != nil {
		fmt.Println("fetchMangasailChapterDetailsContent, error:", err)
		return
	}
	fmt.Println("fetchMangasailChapterDetailsContent, result:", content.ToString())
}
