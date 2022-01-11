package kindle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	_ "embed"

	"github.com/aki237/nscjar"
	"github.com/bmaupin/go-epub"
	"github.com/gocolly/colly"
	"github.com/olegsu/go-tools/pkg/logger"
	"github.com/spf13/cobra"
	mail "github.com/xhit/go-simple-mail/v2"
)

type sendCmdFlagsOptions struct {
	cookies                  string
	language                 string
	titleSelector            string
	kindleEMail              string
	originEmail              string
	originEmailPassword      string
	mainPageSelector         string
	contentPageTitleSelector string
	contentPageTextSelector  string
}

var sendCmdFlags = sendCmdFlagsOptions{}
var sendCmd = &cobra.Command{
	Use:  "send",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lgr := logger.New()
		lgr.Info("starting", "args", args)
		link := args[0]
		cookies := []*http.Cookie{}
		if sendCmdFlags.cookies != "" {
			c, err := loadCookies(sendCmdFlags.cookies, lgr)
			if err != nil {
				panic(err)
			}
			cookies = c
		}

		title, links, err := scrapeMainPage(sendCmdFlags.titleSelector, sendCmdFlags.mainPageSelector, link, cookies)
		if err != nil {
			panic(err)
		}

		if title == "" {
			panic("no title")
		}

		lgr.Info("parsing content", "title", title)

		ebook := epub.NewEpub(title)
		for _, u := range links {
			lgr.Info("Visiting", "url", u)
			title, content, err := scrapeContentPage(sendCmdFlags.contentPageTitleSelector, sendCmdFlags.contentPageTextSelector, u, cookies)
			if err != nil {
				lgr.Info("failed to scrape page", "url", u, "error", err.Error())
				continue
			}
			ebook.AddSection(content, title, "", "")
		}
		ebook.SetLang(sendCmdFlags.language)
		out := &bytes.Buffer{}
		ebook.WriteTo(out)

		if err := os.WriteFile(title+".epub", out.Bytes(), os.ModePerm); err != nil {
			lgr.Info("failed to save local epub version")
		}

		if sendCmdFlags.kindleEMail != "" {
			lgr.Info("Sending email", "kindle", sendCmdFlags.kindleEMail)
			if err := sendEmail(sendCmdFlags.kindleEMail, sendCmdFlags.originEmail, sendCmdFlags.originEmailPassword, title, out.Bytes()); err != nil {
				panic(err)
			}
		}

		lgr.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&sendCmdFlags.cookies, "cookies", "c", "", "Cookie file")
	sendCmd.Flags().StringVarP(&sendCmdFlags.language, "language", "l", "eng", "Book main language")
	sendCmd.Flags().StringVarP(&sendCmdFlags.titleSelector, "title-selector", "t", "", "Book Title")
	sendCmd.Flags().StringVarP(&sendCmdFlags.kindleEMail, "kindle-email", "k", "", "Kindle email address")
	sendCmd.Flags().StringVarP(&sendCmdFlags.originEmail, "origin-email", "e", "", "Origin email address")
	sendCmd.Flags().StringVarP(&sendCmdFlags.originEmailPassword, "origin-email-password", "p", "", "Origin email password")
	sendCmd.Flags().StringVarP(&sendCmdFlags.mainPageSelector, "href-selector", "a", "a[href]", "HTML selector to be used in the main page")
	sendCmd.Flags().StringVar(&sendCmdFlags.contentPageTitleSelector, "content-title-selector", "h1", "HTML selector to be used every content page to find the title")
	sendCmd.Flags().StringVar(&sendCmdFlags.contentPageTextSelector, "content-selector", "p", "HTML selector to be used every content page to find the articles content")
}

func sendEmail(kindle string, origin string, password string, bookName string, bookData []byte) error {
	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("From Me <%s>", origin))
	email.AddTo(kindle)

	email.SetBody(mail.TextPlain, "")
	email.AddAttachmentData(bookData, bookName+".txt", "text/plain")
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = origin
	server.Password = password
	server.Encryption = mail.EncryptionTLS
	smtp, err := server.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to smtp server: %w", err)
	}
	smtp.SendTimeout = time.Minute

	if err := email.Send(smtp); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func scrapeMainPage(titleSelector string, selector string, link string, cookies []*http.Cookie) (string, []string, error) {
	title := ""
	res := []string{}
	c := colly.NewCollector()
	if err := c.SetCookies(link, cookies); err != nil {
		return "", nil, err
	}

	c.OnHTML(titleSelector, func(h *colly.HTMLElement) {
		if h.Text != "" {
			title = h.Text
			return
		}
	})

	c.OnHTML(selector, func(h *colly.HTMLElement) {
		link := h.Attr("href")
		if link == "" {
			return
		}
		res = append(res, link)
	})
	err := c.Visit(link)
	return strings.ReplaceAll(strings.TrimSpace(title), "\n", ""), res, err
}

func scrapeContentPage(titleSelector string, contentSelector string, link string, cookies []*http.Cookie) (string, string, error) {
	title := ""
	content := strings.Builder{}
	c := colly.NewCollector()
	if err := c.SetCookies(link, cookies); err != nil {
		panic(err)
	}

	c.OnHTML(titleSelector, func(h *colly.HTMLElement) {
		title = h.Text
		content.WriteString(fmt.Sprintf("<h1> %s </h1>", h.Text))
	})

	c.OnHTML(contentSelector, func(h *colly.HTMLElement) {
		content.WriteString(fmt.Sprintf("<%s> %s </%s>", h.Name, h.Text, h.Name))
	})
	err := c.Visit(link)
	return title, content.String(), err
}

func loadCookies(path string, lgr *logger.Logger) ([]*http.Cookie, error) {
	lgr.Info("loading cookie file")
	f, err := ioutil.ReadFile(sendCmdFlags.cookies)
	if err != nil {
		return nil, fmt.Errorf("failed to load cookie file: %w", err)
	}
	jar := nscjar.Parser{}
	c, err := jar.Unmarshal(bytes.NewReader(f))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cookie file: %w", err)
	}
	return c, nil
}
