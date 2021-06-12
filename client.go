package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Client interface {
	Login(username, password string) error
	GetDepositsBalance() error
}

type client struct {
	httpClient *http.Client
	log        *log.Logger
	baseURL    string
}

func NewClient() (Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient.Jar = jar

	return &client{
		httpClient: httpClient,
		log:        log.Default(),
		baseURL:    `https://sucursalpersonas.transaccionesbancolombia.com`,
	}, nil
}

func (c *client) Login(username, password string) error {
	u := fmt.Sprintf(`%s/mua/initAuthProcess`, c.baseURL)
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return fmt.Errorf("could not init auth process: %w", err)
	}

	// c.log.Println("defering body close")
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}

	loginUserForm := getElementByID(doc, "loginUserForm")
	if loginUserForm == nil {
		return fmt.Errorf("could not find login user form: %w", err)
	}

	action := getAttribute(loginUserForm, "action")

	values := url.Values{
		"username":     []string{username},
		"device_id":    []string{""},
		"userlanguage": []string{"en-US"},
		"deviceprint":  []string{""},
		"pgid":         []string{""},
		"uievent":      []string{""},
	}

	resp, err = c.httpClient.PostForm(c.baseURL+action, values)
	if err != nil {
		return fmt.Errorf("could not submit loginUserForm: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error submitting loginUserForm: %s", resp.Status)
	}

	/*doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}

	loginUserForm = getElementByID(doc, "loginUserForm")
	if loginUserForm == nil {
		return fmt.Errorf("could not find login user form: %w", err)
	}

	action = getAttribute(loginUserForm, "action")
	values = url.Values{
		"id_ss":        []string{"m%2F9O6xc%2F74KR2OAOyIWYj%2BDEYONRbr1qxJawcOwyiW2bfvTPZGL1AKxFvg8kyoP8%2FpopdqyyLC4rZDQy0P1n18xhtSAc3aj2k%2BnEGPpcZn9Jv%2Bcmy4La%2B8adCn6678MSdd9SKUctVM9hHoD%2F4KkW7HTLgu%2Brl3caboNLYPpKzg3LDSY9NZMbTLq0NTj90cm%2BIOnrcSSVuQcs0QqNv1OVAf5De2pouOGi83tdNWHrlpQS4Rj2fjxZ1v349S0vXlj%2FvaZs%2BBoe7%2FV3tijKwPWvfUbrp8sVZLCUiHohItxaQ0uORELCBDP7FhR3xel2jUF6X7BWj%2BEo1T23Za9EqJUKVg%3D%3D"},
		"tempUserID":   []string{""},
		"HIT_KEY":      []string{"0"},
		"HIT_VKEY":     []string{"0"},
		"userId":       []string{""},
		"password":     []string{"uvdEMkTtiXlW"},
		"device_id":    []string{""},
		"userlanguage": []string{"en-US"},
		"deviceprint":  []string{""},
		"pgid":         []string{""},
		"uievent":      []string{""},
	}
	resp, err = c.httpClient.PostForm(c.baseURL+action, values)
	if err != nil {
		return fmt.Errorf("could not submit loginUserForm: %w", err)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.StatusCode)
	*/
	io.Copy(os.Stdout, resp.Body)

	return nil
}

func escape(s string) string {
	escaped := url.QueryEscape(s)
	escaped = strings.ReplaceAll(escaped, `~`, `%7E`)
	escaped = strings.ReplaceAll(escaped, `!`, `%21`)
	escaped = strings.ReplaceAll(escaped, `*`, `%2A`)
	escaped = strings.ReplaceAll(escaped, `(`, `%28`)
	escaped = strings.ReplaceAll(escaped, `)`, `%29`)
	escaped = strings.ReplaceAll(escaped, `'`, `%27`)
	escaped = strings.ReplaceAll(escaped, `-`, `%2D`)
	escaped = strings.ReplaceAll(escaped, `_`, `%5F`)
	escaped = strings.ReplaceAll(escaped, `.`, `%2E`)
	return escaped
}

func encodeDevicePrint(devicePrint string) string {
	return escape(escape(devicePrint))
}

func decodeDevicePrint(devicePrint string) (string, error) {
	a, err := url.QueryUnescape(devicePrint)
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(a)
}

func (c *client) GetDepositsBalance() error {
	cst := `UHbyzWG7x0g40Q3DmG4IQ9mEVlHIOE1cp8aXhf9Rgex8%2BnX6g%2Bq3G%2BLyxr4kpwqg`
	csrfToken := `1804780420988843334`
	cookie := `JSESSIONID=0RxB7Vd23fc-lHnZrKZEVGfF; PREFS="X9RLQLaAmqHW8JqmdTMTk0MFBxw="; NSC_JOr2zhh2e44kdkqd4uupeqdgxr1z1c0=4150a3cb0364cd0d3ea8dab25b436be0cee4445a7e513fa4949b4998994b0c64f398c052; __cflb=02DiuF7aX6zsQEVJrpLGtHaWFTk3VhwPwDGH9EUtDMv7v; UUID=354673ff5f6608abc152664feaab0e5b; SMCHALLENGE=SSL_CHALLENGE_DONE; SMTEXT=SUCCESS; SMSESSION=kQX5V5f0dC3PJyq6Snj1zQrp9anupBACjGhYWFnV9TNmBqmZEI/j8MHL/BLv8xPLlvplHBrIN3YXSlkKyFpj7wer0825qFF63+PgTl567Q7u31whmM2NmXGXilnbuzeg5cAvmF1uzS/4XzvEdtG8xJRbqqh9WULVinssUTJWFEWI4LfGK0hShiLMHiCbLXV4TNc9Y/IWor8iZjA2Qd33Z21qktMGKCNSNfy4lc1W/NDhqierf+7QmMS0en9YbV7ygUXAAUxt6HGjsPPCij8/AoE1s14CPYgEm7MY+7e4XJANLYs6B+o5a7xkn4e+R353mdGOG/dCQCIB0ejvjD9SkZcppQf1y3gHoSykaCpXuoh4YYKvJI5Bk9B9nUR1R4tfAUgVQqOVaBf2T5mmHAf1ydO6+lPspNzwjI/QG2kn9EPG3QRqwrLYdHLJUDKOHxs93X50tM+I7M1/HvBP7wS3RnKsDvoYZhvaQYboBTQtXTs2r2FrIudtHBAEfdHe6VdPS2ECUf38FF/cXtMVHGcPb4+MWTQ8urisrYTKGf9alDP94VuFUrjqDPqd6DRwMOvT7XE8QFmCqeXF7S55Ii/f0//lYoeN9ivNQAFkwyxLSD0qMDH0QIcRBP80vDFehWwxu5/Jnnl//AtvsvRZn+9X7brcVtH8M93nw3wFgSqJm88SsU8aYusaiGkd1jk7CymJ8PKiFwHOYujeQU9pYsyoRMW+JHe2nVCLChWxQvGQ89aL89pz0CRKFc6spdDPP3X0Dslw9W9dCsCUuVORZ0YGNdIG0dg0QoALrWQkTskWMmXtvJgX9octZjn0tkf6yydymErZjUCwxtJCAaq/ygCi2Ifnd6rxs8Rza6IlG63NWJSZ+Xy4b+okNIaF5L2CCx8ytzBgIdb0LAzFWS2dE+giZJ6nX77dpNDFRwMn0uB/y2ez4XIbXzeHj+3bHxZk87y2Y01MvMmLPTbsPFzSj2bv1LTuHY1XPV2W; T1_OLBP_COOKIE=""`

	u := fmt.Sprintf(`https://sucursalpersonas.transaccionesbancolombia.com/cb/pages/jsp/account/getDepositsBalanceBancolombiaHome.action?cst=%s`, cst)
	bodyStr := fmt.Sprintf(`type=DEPOSITS&CSRF_TOKEN=%s&cst=%s`, csrfToken, cst)
	body := strings.NewReader(bodyStr)

	req, err := http.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0")
	req.Header.Set("Accept", "*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://sucursalpersonas.transaccionesbancolombia.com")
	req.Header.Set("Referer", "https://sucursalpersonas.transaccionesbancolombia.com/cb/pages/jsp/home/index.jsp")
	req.Header.Set("Cookie", cookie)

	c.log.Println("sending request")
	resp, err := c.httpClient.Do(req)
	c.log.Println("request sent")
	if err != nil {
		return err
	}

	c.log.Println("defering body close")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %w", err)
	}

	c.log.Println("reading request")
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(content))
	return nil
}
