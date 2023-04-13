package scraper

import (
	"context"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"net/url"
	"optimus/internal/models"
	"optimus/pkg"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl = "https://www.behance.net"
)

type BehanceScraper interface {
	Search(ctx context.Context, search string) ([]*models.SearchModel, error)
}

type behanceScraper struct {
	pw *playwright.Playwright
}

func NewBehanceScraper(pw *playwright.Playwright) *behanceScraper {
	return &behanceScraper{
		pw: pw,
	}
}

func (sc *behanceScraper) Search(ctx context.Context, search string) ([]*models.SearchModel, error) {
	ctx, sp := otel.Tracer("").Start(ctx, "Search")
	defer sp.End()

	pwBrowser := pkg.NewPlaywrightPackage(sc.pw, pkg.PlaywrightConfiguration{
		BrowserName: "firefox",
		IsHeadless:  true,
	})
	defer pwBrowser.CloseBrowser(ctx)

	err := pwBrowser.GeneratePlaywrightBrowser(ctx)
	if err != nil {
		return []*models.SearchModel{}, err
	}

	params := url.Values{}
	params.Set("tracking_source", "typeahead_search_direct")
	params.Set("search", search)

	urlObj, err := url.Parse(baseUrl)
	if err != nil {
		return []*models.SearchModel{}, err
	}

	urlObj.RawQuery = params.Encode()
	log.WithContext(ctx).Info("Start scraping behance")
	_, err = pwBrowser.Page.Goto(urlObj.String())
	if err != nil {
		return []*models.SearchModel{}, err
	}
	time.Sleep(time.Second * 2)

	// Set timeout to 30 seconds
	timeout := time.After(30 * time.Second)

	// Scroll until end of page
	for {
		var endOfPage playwright.ElementHandle
		select {
		// If the timeout expires, exit the loop
		case <-timeout:
			log.WithContext(ctx).Warn("Unable to reach the end of page due to timeout")
			break
		// Otherwise, keep scrolling
		default:
			// Force remove google login popup
			pwBrowser.Page.Evaluate(`document.querySelector("iframe").remove()`)

			pwBrowser.Page.Keyboard().Press("End")
			endOfPage, _ = pwBrowser.Page.QuerySelector(".blockerEndingShort")
			if endOfPage == nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			break
		}

		if endOfPage != nil {
			break
		}
	}

	contentBox, err := pwBrowser.Page.QuerySelectorAll("div[class*=\"Projects-grid\"]")
	if err != nil {
		return []*models.SearchModel{}, err
	}

	var searchResult []*models.SearchModel
	for _, content := range contentBox {
		var itemBox []playwright.ElementHandle
		if class, _ := content.GetAttribute("class"); strings.Contains(class, "ContentGrid") {
			itemBox, err = content.QuerySelectorAll("li")
			if err != nil {
				return []*models.SearchModel{}, err
			}
		} else {
			itemBox, err = content.QuerySelectorAll("> *")
			if err != nil {
				return []*models.SearchModel{}, err
			}
		}

		for _, item := range itemBox {
			projectUrlElement, _ := item.QuerySelector("a[title*=\"Link to project\"]")
			projectUrl, _ := projectUrlElement.GetAttribute("href")

			imageUrlElement, _ := item.QuerySelector("img")
			imageUrl, _ := imageUrlElement.GetAttribute("src")

			ownerElement, _ := item.QuerySelector("[class*=\"TitleOwner\"]")
			titleElement, _ := ownerElement.QuerySelector("a")
			title, _ := titleElement.InnerText()
			authorElement, _ := ownerElement.QuerySelector("div")
			author, _ := authorElement.InnerText()

			statsElement, _ := item.QuerySelector("[class*=\"ProjectCover-stats\"]")
			stats, _ := statsElement.QuerySelectorAll("span")
			likesElement, _ := stats[0].GetAttribute("title")
			viewsElement, _ := stats[1].GetAttribute("title")

			likes, _ := strconv.Atoi(strings.ReplaceAll(likesElement, ",", ""))
			views, _ := strconv.Atoi(strings.ReplaceAll(viewsElement, ",", ""))

			searchResult = append(searchResult, &models.SearchModel{
				ProjectUrl: projectUrl,
				ImageUrl:   imageUrl,
				Title:      strings.TrimSpace(title),
				Author:     author,
				TotalLikes: likes,
				TotalViews: views,
			})
		}
	}
	log.WithContext(ctx).Info("Successfully retrieved ", len(searchResult), " results")

	return searchResult, nil
}
