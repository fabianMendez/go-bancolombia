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
	Logout() error
	GetDepositsBalance() error
}

type client struct {
	httpClient *http.Client
	log        *log.Logger
	baseURL    string
	cst        string
	csrfToken  string
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
func (c *client) updateCsrfToken(doc *html.Node) {
	csrfToken := parseCsrfToken(doc)
	if csrfToken != "" {
		fmt.Println("CSRF =>", csrfToken)
		c.csrfToken = csrfToken
	}
}

func (c *client) Login(username, password string) error {
	deviceId := ""
	userlanguage := "en-US"
	deviceprint := ""
	pgid := ""
	uievent := ""

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

	c.updateCsrfToken(doc)
	loginUserForm := getElementByID(doc, "loginUserForm")
	if loginUserForm == nil {
		return fmt.Errorf("could not find login user form: %w", err)
	}

	action := getAttribute(loginUserForm, "action")

	values := url.Values{
		"username":     []string{username},
		"device_id":    []string{deviceId},
		"userlanguage": []string{userlanguage},
		"deviceprint":  []string{deviceprint},
		"pgid":         []string{pgid},
		"uievent":      []string{uievent},
	}

	resp, err = c.httpClient.PostForm(c.baseURL+action, values)
	if err != nil {
		return fmt.Errorf("could not submit loginUserForm: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error submitting loginUserForm: %s", resp.Status)
	}

	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}

	c.updateCsrfToken(doc)
	loginUserForm = getElementByID(doc, "loginUserForm")
	if loginUserForm == nil {
		return fmt.Errorf("could not find login user form: %w", err)
	}

	t1Assertion := parseT1Assertion(doc)
	keyboardSrc := parseKeyboardContent(doc)
	keyboardNode, err := html.Parse(strings.NewReader(keyboardSrc))
	if err != nil {
		return fmt.Errorf("could not parse keyboard: %w", err)
	}
	keyMap := parseKeyboardMap(keyboardNode)
	password = mapPassword(keyMap, password)
	initRngPool()
	idSs := processPassword(password, t1Assertion)
	passwordInputName := parsePasswordInputName(doc)

	action = getAttribute(loginUserForm, "action")
	values = url.Values{
		"id_ss": []string{idSs},
		// "id_ss":        []string{"m%2F9O6xc%2F74KR2OAOyIWYj%2BDEYONRbr1qxJawcOwyiW2bfvTPZGL1AKxFvg8kyoP8%2FpopdqyyLC4rZDQy0P1n18xhtSAc3aj2k%2BnEGPpcZn9Jv%2Bcmy4La%2B8adCn6678MSdd9SKUctVM9hHoD%2F4KkW7HTLgu%2Brl3caboNLYPpKzg3LDSY9NZMbTLq0NTj90cm%2BIOnrcSSVuQcs0QqNv1OVAf5De2pouOGi83tdNWHrlpQS4Rj2fjxZ1v349S0vXlj%2FvaZs%2BBoe7%2FV3tijKwPWvfUbrp8sVZLCUiHohItxaQ0uORELCBDP7FhR3xel2jUF6X7BWj%2BEo1T23Za9EqJUKVg%3D%3D"},
		"tempUserID": []string{""},
		"HIT_KEY":    []string{"0"},
		"HIT_VKEY":   []string{"0"},
		"userId":     []string{""},
		"password":   []string{passwordInputName},
		// "password":     []string{"uvdEMkTtiXlW"},
		"device_id":    []string{deviceId},
		"userlanguage": []string{userlanguage},
		"deviceprint":  []string{deviceprint},
		"pgid":         []string{pgid},
		"uievent":      []string{uievent},
	}
	resp, err = c.httpClient.PostForm(c.baseURL+action, values)
	if err != nil {
		return fmt.Errorf("could not submit loginUserForm: %w", err)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.StatusCode)
	// io.Copy(os.Stdout, resp.Body)
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "post-return")
	if err != nil {
		return fmt.Errorf("could not submit post-return form: %w", err)
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

	resp, err = c.httpClient.Get(c.baseURL + "/mua/CONTINUE_SM")
	if err != nil {
		return fmt.Errorf("could not submit request: %w", err)
	}
	defer resp.Body.Close()
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "post-return")
	if err != nil {
		return fmt.Errorf("could not submit post-return form: %w", err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}

	c.updateCsrfToken(doc)
	action = parseUrlRedirect(doc)
	action = filterUrl(action)
	tokenM := parseTokenMua(doc)
	code := parseCodeRedirect(doc)
	resp, err = c.submitFormValues(doc, resp.Request.URL.String(), "post-login", action, url.Values{
		"tokenM": []string{tokenM},
		"code":   []string{code},
	})
	if err != nil {
		return fmt.Errorf("could not submit post-login form: %w", err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}

	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "post-link-mada")
	if err != nil {
		return fmt.Errorf("could not submit post-link-mada form: %w", err)
	}
	defer resp.Body.Close()
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "post-link-index")
	if err != nil {
		return fmt.Errorf("could not submit post-link-index form: %w", err)
	}
	defer resp.Body.Close()

	//post-login
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "post-login")
	if err != nil {
		return fmt.Errorf("could not submit post-login form: %w", err)
	}
	defer resp.Body.Close()

	//invocacion
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "invocacion")
	if err != nil {
		return fmt.Errorf("could not submit invocacion form: %w", err)
	}
	defer resp.Body.Close()

	//validateUser
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "validateUser")
	if err != nil {
		return fmt.Errorf("could not submit validateUser form: %w", err)
	}
	defer resp.Body.Close()

	//loginSimulateFormID
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "loginSimulateFormID")
	if err != nil {
		return fmt.Errorf("could not submit loginSimulateFormID form: %w", err)
	}
	defer resp.Body.Close()

	// loginForm
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "loginForm")
	if err != nil {
		return fmt.Errorf("could not submit loginForm form: %w", err)
	}
	defer resp.Body.Close()

	// loginForm1
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse html document: %w", err)
	}
	c.updateCsrfToken(doc)
	resp, err = c.submitForm(doc, resp.Request.URL.String(), "loginForm1")
	if err != nil {
		return fmt.Errorf("could not submit loginForm1 form: %w", err)
	}
	defer resp.Body.Close()

	// mainPage
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse login-init document: %w", err)
	}
	c.updateCsrfToken(doc)
	action, err = c.buildAction(resp.Request.URL.String(), parseLocationReplace(doc))
	if err != nil {
		return err
	}
	resp, err = c.httpClient.Get(action)
	if err != nil {
		return fmt.Errorf("could not load mainPage: %w", err)
	}
	defer resp.Body.Close()

	// index
	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse mainPage document: %w", err)
	}

	c.updateCsrfToken(doc)
	c.cst = parseCstParam(doc)
	fmt.Println("cstParam1:", c.cst)

	action, err = c.buildAction(resp.Request.URL.String(), "index.jsp")
	if err != nil {
		return err
	}
	resp, err = c.httpClient.Get(action)
	if err != nil {
		return fmt.Errorf("could not load mainPage: %w", err)
	}
	defer resp.Body.Close()

	doc, err = html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse index document: %w", err)
	}

	c.updateCsrfToken(doc)
	c.cst = parseCstParam(doc)
	fmt.Println("cstParam2:", c.cst)

	return nil
}

func (c *client) buildAction(baseURL, action string) (string, error) {
	fmt.Println(baseURL)
	if !strings.HasPrefix(action, "http://") && !strings.HasPrefix(action, "https://") {
		if strings.HasPrefix(action, "/") {
			return c.baseURL + action, nil
		}

		parent := baseURL[:strings.LastIndex(baseURL, "/")]
		return parent + "/" + action, nil
	}

	actionurl, err := url.Parse(action)
	if err != nil {
		return "", fmt.Errorf("could not parse action url: %w", err)
	}
	return actionurl.String(), nil
}

func (c *client) submitForm(doc *html.Node, baseURL, id string) (*http.Response, error) {
	form := getElementByID(doc, id)
	action := getAttribute(form, "action")
	return c.submitFormValues(doc, baseURL, id, action, nil)
}

func (c *client) submitFormValues(doc *html.Node, baseURL, id, action string, values url.Values) (*http.Response, error) {
	form := getElementByID(doc, id)
	fields := parseFormFields(form)
	replaceValues(fields, values)

	u, err := c.buildAction(baseURL, action)
	if err != nil {
		return nil, err
	}

	return c.httpClient.PostForm(u, fields)
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
func (c *client) cstUrl(u string) string {
	param := "cst=" + c.cst
	if strings.Contains(u, "?") {
		return u + "&" + param
	}
	return u + "?" + param
}

func (c *client) postAjax(endpoint string) error {
	return c.postAjaxValues(endpoint, url.Values{})
}

func (c *client) postAjaxValues(endpoint string, values url.Values) error {
	replaceValues(values, url.Values{
		"cst":        []string{c.csrfToken},
		"CSRF_TOKEN": []string{c.csrfToken},
	})

	u := fmt.Sprintf(`%s%s`, c.baseURL, endpoint)
	resp, err := c.httpClient.PostForm(u, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// doc, err := html.Parse(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// c.updateCsrfToken(doc)

	return nil
}

func (c *client) preGetDepositsBalance() error {
	var err error

	/*
		err = c.postAjax(c.cstUrl("/cb/pages/jsp-ns/olb/ChangeNameACHDataEntry"))
		if err != nil {
			return err
		}

		err = c.postAjax(c.cstUrl("/cb/pages/jsp-ns/olb/PersonalizeProductsNameDataEntry"))
		if err != nil {
			return err
		}

		err = c.postAjax("/cb/pages/jsp-ns/olb/UsersPreferencesAction")
		if err != nil {
			return err
		}

		err = c.postAjax("/cb/pages/jsp/updateData/getUpdateDynamicData.action")
		if err != nil {
			return err
		}

		err = c.postAjax(c.cstUrl("/cb/pages/jsp-ns/olb/PreApprovedAction"))
		if err != nil {
			return err
		}
	*/

	/*
		err = c.postAjaxValues("/cb/pages/jsp/ga/GATokenGeneration.action", url.Values{
			"id_ga":    []string{"anuncio-0"},
			"timeZone": []string{"GMT-0500"},
		})
		if err != nil {
			return err
		}
	*/

	// falta agregarle el parámetro _={timestamp}
	_, err = c.httpClient.Get(c.baseURL + c.cstUrl("/cb/pages/jsp-ns/olb/InitAccountSummary?redirect=ALLACCOUNTS_HOME&type=ALLACCOUNTS_HOME"))
	if err != nil {
		return err
	}

	/*
		err = c.postAjaxValues("/cb/pages/jsp/ga/GATokenGeneration.action", url.Values{
			"id_ga":    []string{"anuncio-0"},
			"timeZone": []string{"GMT-0500"},
		})
		if err != nil {
			return err
		}
	*/

	/*
		// falta agregarle el parámetro _={timestamp}
		_, err = c.httpClient.Get(c.baseURL + c.cstUrl("/cb/web/js/account/account_grid_bancolombia.js?version=3.2.1.RC1"))
		if err != nil {
			return err
		}

		// falta agregarle el parámetro _={timestamp}
		_, err = c.httpClient.Get(c.baseURL + c.cstUrl("/cb/web/js/landing_grid.js?version=3.2.1.RC1"))
		if err != nil {
			return err
		}

		// falta agregarle el parámetro _={timestamp}
		_, err = c.httpClient.Get(c.baseURL + c.cstUrl("/cb/web/js/datePicker/datePicker.js?version=3.2.1.RC1"))
		if err != nil {
			return err
		}

		// falta agregarle el parámetro _={timestamp}
		_, err = c.httpClient.Get(c.baseURL + c.cstUrl("/cb/web/js/account/account_grid_cards_credits.js?version=3.2.1.RC1"))
		if err != nil {
			return err
		}
	*/

	// err = c.postAjaxValues("/cb/pages/jsp/account/getDepositsBalanceBancolombiaHome.action", url.Values{
	// 	"id_ga":    []string{"anuncio-0"},
	// 	"timeZone": []string{"GMT-0500"},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *client) Logout() error {
	var err error
	_, err = c.httpClient.Get(c.baseURL + c.cstUrl("/cb/pages/jsp-ns/olb/SafeExit"))
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetDepositsBalance() error {
	err := c.preGetDepositsBalance()
	if err != nil {
		return err
	}
	// cst := `UHbyzWG7x0g40Q3DmG4IQ9mEVlHIOE1cp8aXhf9Rgex8%2BnX6g%2Bq3G%2BLyxr4kpwqg`
	// csrfToken := `1804780420988843334`
	// cookie := `JSESSIONID=0RxB7Vd23fc-lHnZrKZEVGfF; PREFS="X9RLQLaAmqHW8JqmdTMTk0MFBxw="; NSC_JOr2zhh2e44kdkqd4uupeqdgxr1z1c0=4150a3cb0364cd0d3ea8dab25b436be0cee4445a7e513fa4949b4998994b0c64f398c052; __cflb=02DiuF7aX6zsQEVJrpLGtHaWFTk3VhwPwDGH9EUtDMv7v; UUID=354673ff5f6608abc152664feaab0e5b; SMCHALLENGE=SSL_CHALLENGE_DONE; SMTEXT=SUCCESS; SMSESSION=kQX5V5f0dC3PJyq6Snj1zQrp9anupBACjGhYWFnV9TNmBqmZEI/j8MHL/BLv8xPLlvplHBrIN3YXSlkKyFpj7wer0825qFF63+PgTl567Q7u31whmM2NmXGXilnbuzeg5cAvmF1uzS/4XzvEdtG8xJRbqqh9WULVinssUTJWFEWI4LfGK0hShiLMHiCbLXV4TNc9Y/IWor8iZjA2Qd33Z21qktMGKCNSNfy4lc1W/NDhqierf+7QmMS0en9YbV7ygUXAAUxt6HGjsPPCij8/AoE1s14CPYgEm7MY+7e4XJANLYs6B+o5a7xkn4e+R353mdGOG/dCQCIB0ejvjD9SkZcppQf1y3gHoSykaCpXuoh4YYKvJI5Bk9B9nUR1R4tfAUgVQqOVaBf2T5mmHAf1ydO6+lPspNzwjI/QG2kn9EPG3QRqwrLYdHLJUDKOHxs93X50tM+I7M1/HvBP7wS3RnKsDvoYZhvaQYboBTQtXTs2r2FrIudtHBAEfdHe6VdPS2ECUf38FF/cXtMVHGcPb4+MWTQ8urisrYTKGf9alDP94VuFUrjqDPqd6DRwMOvT7XE8QFmCqeXF7S55Ii/f0//lYoeN9ivNQAFkwyxLSD0qMDH0QIcRBP80vDFehWwxu5/Jnnl//AtvsvRZn+9X7brcVtH8M93nw3wFgSqJm88SsU8aYusaiGkd1jk7CymJ8PKiFwHOYujeQU9pYsyoRMW+JHe2nVCLChWxQvGQ89aL89pz0CRKFc6spdDPP3X0Dslw9W9dCsCUuVORZ0YGNdIG0dg0QoALrWQkTskWMmXtvJgX9octZjn0tkf6yydymErZjUCwxtJCAaq/ygCi2Ifnd6rxs8Rza6IlG63NWJSZ+Xy4b+okNIaF5L2CCx8ytzBgIdb0LAzFWS2dE+giZJ6nX77dpNDFRwMn0uB/y2ez4XIbXzeHj+3bHxZk87y2Y01MvMmLPTbsPFzSj2bv1LTuHY1XPV2W; T1_OLBP_COOKIE=""`

	u := fmt.Sprintf(`%s/cb/pages/jsp/account/getDepositsBalanceBancolombiaHome.action?cst=%s`, c.baseURL, c.cst)
	bodyStr := fmt.Sprintf(`type=DEPOSITS&CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
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
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/cb/pages/jsp/home/index.jsp")
	// req.Header.Set("Cookie", cookie)

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

func mapPassword(keymap map[string]string, password string) string {
	newPassword := ""

	for _, c := range password {
		if val, found := keymap[string(c)]; found {
			newPassword += val
		}
	}

	return newPassword
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

func replaceValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = vs
	}
}

func filterUrl(u string) string {
	const qst = "?"
	const semCol = ";"
	if strings.Contains(u, semCol) {
		u2 := u[:strings.Index(u, semCol)]
		if strings.Contains(u, qst) {
			u = u2 + u[strings.Index(u, qst):]
		} else {
			u = u2
		}
	}
	return u
}
