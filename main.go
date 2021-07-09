package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bigscreen/manga-scrapper/scrapper/mangasail"
	"github.com/bigscreen/manga-scrapper/service"
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
	mService := service.NewMangsailService(service.MangasailServiceParams{
		HomeScrapper:           mangasail.NewHomePageScrapper(chromeCtx),
		MangaDetailsScrapper:   mangasail.NewMangaDetailsPageScrapper(chromeCtx),
		ChapterDetailsScrapper: mangasail.NewChapterDetailsPageScrapper(chromeCtx),
	})
	fetchMangasailHomeContent(mService)
	fmt.Println("---------------------")
	fetchMangasailMangaDetailsContent(mService)
	fmt.Println("---------------------")
	fetchMangasailChapterDetailsContent(mService)
}

func fetchMangasailHomeContent(svc service.FetchService) {
	content, err := svc.GetHomeCards()
	if err != nil {
		fmt.Println("fetchMangasailHomeContent, error:", err)
		return
	}
	fmt.Println("fetchMangasailHomeContent, result:", content.String())
}

func fetchMangasailMangaDetailsContent(svc service.FetchService) {
	content, err := svc.GetMangaDetails("checkmate-manga")
	if err != nil {
		fmt.Println("fetchMangasailMangaDetailsContent, error:", err)
		return
	}
	fmt.Println("fetchMangasailMangaDetailsContent, result:", content.String())
}

func fetchMangasailChapterDetailsContent(svc service.FetchService) {
	content, err := svc.GetChapterDetails("checkmate-35")
	if err != nil {
		fmt.Println("fetchMangasailChapterDetailsContent, error:", err)
		return
	}
	fmt.Println("fetchMangasailChapterDetailsContent, result:", content.String())
}
