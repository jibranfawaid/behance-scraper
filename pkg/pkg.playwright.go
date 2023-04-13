package pkg

import (
	"context"
	"errors"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"math/rand"
	e "optimus/internal/errors"
	"optimus/internal/utilities"
	"runtime"
	"strings"
	"time"
)

func NewPlaywright() (*playwright.Playwright, error) {
	return playwright.Run()
}

const (
	CHROME  = "chrome"
	FIREFOX = "firefox"
)

type PlaywrightPackage interface {
	GeneratePlaywrightBrowser(ctx context.Context) error
	CloseBrowser(ctx context.Context)
}

type PlaywrightConfiguration struct {
	BrowserName string
	IsHeadless  bool
}

type PlaywrightBrowser struct {
	Browser playwright.Browser
	Page    playwright.Page

	pw *playwright.Playwright

	config PlaywrightConfiguration
}

func NewPlaywrightPackage(pw *playwright.Playwright, cfg PlaywrightConfiguration) *PlaywrightBrowser {
	return &PlaywrightBrowser{
		pw:     pw,
		config: cfg,
	}
}

// List of usable user agents
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_7; rv:107.0) Gecko/20010101 Firefox/107.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:107.0) Gecko/20100101 Firefox/107.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0/8mqTxTuL-47/8mqTxTuL-47",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
}

func (p *PlaywrightBrowser) GeneratePlaywrightBrowser(ctx context.Context) error {
	log.WithContext(ctx).Info("Launching browser")

	var bt playwright.BrowserType

	switch strings.ToLower(p.config.BrowserName) {
	case CHROME:
		bt = p.pw.Chromium
	case FIREFOX:
		bt = p.pw.Firefox
	default:
		log.WithContext(ctx).Error("Unable to find browser type")
		return errors.New(e.GeneralError)
	}

	browser, err := bt.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: utilities.NewBoolean(p.config.IsHeadless),
		Timeout:  utilities.NewFloat(90 * 1000), // 90 seconds
		Args: []string{
			"--disable-setuid-sandbox",
			"--disable-dev-shm-usage",
			"--disable-accelerated-2d-canvas",
			"--no-first-run",
			"--no-zygote",
			// "--single-process",
			"--disable-web-security",
			"--disable-gpu",
			"--no-sandbox",
		},
	})
	if err != nil {
		log.WithContext(ctx).Error("Unable to start browser")
		p.CloseBrowser(ctx)

		return errors.New(e.GeneralError)
	}

	// Rotating user agents
	ua := &userAgents[rand.Intn(len(userAgents))]

	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		UserAgent:       ua,
		AcceptDownloads: utilities.NewBoolean(true),
	})
	if err != nil {
		log.WithContext(ctx).Error("Unable to create page")
		p.CloseBrowser(ctx)

		return errors.New(e.GeneralError)
	}

	// Browser lifetime
	time.AfterFunc(10*time.Minute, func() {
		if p.Browser.IsConnected() {
			p.CloseBrowser(ctx)
		}
	})

	p.Browser = browser
	p.Page = page

	log.WithContext(ctx).Info("Browser has successfully launched")

	return nil
}

func (p *PlaywrightBrowser) CloseBrowser(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.WithContext(ctx).Error("Recovered from panic while closing browser: ", err)
		}
	}()

	log.WithContext(ctx).Info("Closing browser")

	for _, bc := range p.Browser.Contexts() {
		for _, Page := range bc.Pages() {
			_ = Page.Close()
		}
		_ = bc.Close()
	}

	_ = p.Browser.Close()
	runtime.GC()

	log.WithContext(ctx).Info("Successfully close the browser")
}
