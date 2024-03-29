package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fabianMendez/htmldom"
	"golang.org/x/net/html"
	"golang.org/x/term"
)

type Client interface {
	Login(username, password string) error
	Logout() error
	GetDepositsBalance() (DepositsBalance, error)
	GetSavingsDetail(id string, page int) ([]SavingsDetail, error)
	AccountEnrroll() error
}

type client struct {
	httpClient       *http.Client
	log              *log.Logger
	baseURL          string
	aditionalHeaders map[string]string
	refererURL       string
	cst              string
	csrfToken        string
	rsa              *rsa
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

	// nullf, err := os.Open("/dev/null")
	// if err != nil {
	// 	return nil, err
	// }
	// l := log.New(nullf, "", log.LstdFlags)
	l := log.Default()

	r := new(rsa)
	r.init()

	baseURL := `https://sucursalpersonas.transaccionesbancolombia.com`
	return &client{
		httpClient: httpClient,
		log:        l,
		baseURL:    baseURL,
		aditionalHeaders: map[string]string{
			"User-Agent":       "Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
			"X-Requested-With": "XMLHttpRequest",
			"Referer":          baseURL + "/cb/pages/jsp/home/index.jsp",
			"Accept-Language":  "en-US,es-CO;q=0.7,en;q=0.3",
			// req.Header.Set("Accept", "*")
			// req.Header.Set("Accept-Language", "en-US,en;q=0.5")
			// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			// req.Header.Set("Origin", c.baseURL)
			// req.Header.Set("Referer", c.baseURL+"/cb/pages/jsp/home/index.jsp")
		},
		rsa: r,
	}, nil
}

func (c *client) updateCsrfToken(doc *html.Node) {
	csrfToken := parseCsrfToken(doc)
	if csrfToken != "" && csrfToken != c.csrfToken {
		c.log.Println("CSRF =>", csrfToken)
		c.csrfToken = csrfToken
	}
}

func (c *client) updateCstParam(doc *html.Node) {
	cstParam := parseCstParam(doc)
	if cstParam != "" && cstParam != c.cst {
		c.log.Println("CST =>", cstParam)
		c.cst = cstParam
	}
}

func (c *client) get(u string) (*http.Response, error) {
	return c.request(http.MethodGet, u, nil)
}

func (c *client) postForm(u string, data url.Values) (*http.Response, error) {
	return c.requestWithHeaders(http.MethodPost, u, strings.NewReader(data.Encode()), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
}

func (c *client) request(method, u string, body io.Reader) (*http.Response, error) {
	return c.requestWithHeaders(method, u, body, nil)
}

func (c *client) requestWithHeaders(method, u string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	for headerKey, headerValue := range c.aditionalHeaders {
		req.Header.Add(headerKey, headerValue)
	}

	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	c.log.Println("sending request")
	resp, err := c.httpClient.Do(req)
	c.log.Println("request sent")
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	refererURL := resp.Request.URL.String()
	if refererURL != "" {
		c.refererURL = refererURL
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			c.log.Printf("retry after %s seconds\n", retryAfter)
			errorMessage, err := io.ReadAll(resp.Body)
			if err == nil {
				return nil, errors.New(string(errorMessage))
			}
		}
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	return resp, nil
}

func (c *client) requestJSON(method, u string, body io.Reader, v interface{}) error {
	return c.requestJSONWithHeaders(method, u, body, v, nil)
}

func (c *client) requestJSONWithHeaders(method, u string, body io.Reader, v interface{}, headers map[string]string) error {
	resp, err := c.requestWithHeaders(method, u, body, headers)
	if err != nil {
		return err
	}
	c.log.Println("defering body close")
	defer resp.Body.Close()

	c.log.Println("reading reponse")
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("could not decode response: %w", err)
	}

	return nil
}

func (c *client) loadHTML(resp *http.Response, err error) (*html.Node, error) {
	if err != nil {
		return nil, err
	}
	c.log.Println("defering body close")
	defer resp.Body.Close()

	c.log.Println("reading reponse")
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse html: %w", err)
	}

	c.updateCsrfToken(doc)
	c.updateCstParam(doc)

	return doc, nil
}

func (c *client) Login(username, password string) error {
	deviceId := ""
	userlanguage := "en-US"
	deviceprint := ""
	pgid := ""
	uievent := ""

	u := fmt.Sprintf(`%s/mua/initAuthProcess`, c.baseURL)
	doc, err := c.loadHTML(c.get(u))
	if err != nil {
		return fmt.Errorf("could not init auth process: %w", err)
	}

	loginUserForm := getElementByID(doc, "loginUserForm")
	if loginUserForm == nil {
		html.Render(os.Stderr, doc)
		return fmt.Errorf("could not find login user form")
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

	doc, err = c.loadHTML(c.postForm(c.baseURL+action, values))
	if err != nil {
		return fmt.Errorf("could not submit loginUserForm: %w", err)
	}

	openTop := getElementByID(doc, "openTop")
	if openTop != nil {
		u := fmt.Sprintf(`%s/mua/initAuthProcess`, c.baseURL)
		doc, err = c.loadHTML(c.get(u))
		if err == nil {
			summary := getElementByID(doc, "summary")
			errorText := strings.TrimSpace(getInnerText(summary))
			if errorText != "" {
				return fmt.Errorf(errorText)
			}
		}
		return fmt.Errorf("invalid username")
	}

	loginUserForm = getElementByID(doc, "loginUserForm")
	if loginUserForm == nil {
		html.Render(os.Stderr, doc)
		return fmt.Errorf("could not find login user form")
	}

	t1Assertion := parseT1Assertion(doc)
	keyboardSrc := parseKeyboardContent(doc)
	keyboardNode, err := html.Parse(strings.NewReader(keyboardSrc))
	if err != nil {
		return fmt.Errorf("could not parse keyboard: %w", err)
	}
	keyMap := parseKeyboardMap(keyboardNode)
	password = mapPassword(keyMap, password)
	idSs := c.rsa.processPassword(password, t1Assertion)
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
		// g-recaptcha-response: ""
	}
	doc, err = c.loadHTML(c.postForm(c.baseURL+action, values))
	if err != nil {
		return fmt.Errorf("could not submit loginUserForm: %w", err)
	}

	// 2021/06/26 09:11:35 could not login: Usuario o clave inválida. Por favor intente de nuevo.
	// could not login: captcha_error_server
	// https://sucursalpersonas.transaccionesbancolombia.com/mua/view
	// https://sucursalpersonas.transaccionesbancolombia.com/mua/CLOSE_ALL?scis=Xj%2FFadZlcYEo%2BlD0R%2FW6WDJy0aqcUoDMJv3VkuJtl1Q%3D
	// log.Println(c.refererURL)
	openTop = getElementByID(doc, "openTop")
	if openTop != nil {
		u := fmt.Sprintf(`%s/mua/initAuthProcess`, c.baseURL)
		doc, err = c.loadHTML(c.get(u))
		if err == nil {
			summary := getElementByID(doc, "summary")
			errorText := strings.TrimSpace(getInnerText(summary))
			if errorText != "" {
				return fmt.Errorf(errorText)
			}
		}
		return fmt.Errorf("invalid credentials")
	}

	for {
		locationReplace := parseLocationReplace(doc)
		// if locationReplace != "" {
		// 	fmt.Println(locationReplace)
		// }

		if strings.HasSuffix(locationReplace, "mainPage.jsp") {
			action, err = c.buildAction(locationReplace)
			if err != nil {
				return err
			}
			fmt.Println(action)

			_, err = c.loadHTML(c.get(action))
			if err != nil {
				return fmt.Errorf("could not load mainPage: %w", err)
			}
			break
		}

		form := htmldom.GetElementByTag(doc, "form")
		if form == nil {
			html.Render(os.Stderr, doc)
			break
		}

		formID := htmldom.GetAttribute(form, "id")
		// fmt.Println("formID", formID)
		if formID == "post-login" {
			tokenM := parseTokenMua(doc)
			code := parseCodeRedirect(doc)

			if tokenM != "" && code != "" {
				action = parseUrlRedirect(doc)
				action = filterUrl(action)
				// fmt.Println("post login form:", action)

				doc, err = c.loadHTML(c.submitFormValues(doc, "post-login", action, url.Values{
					"tokenM": []string{tokenM},
					"code":   []string{code},
				}))
				if err != nil {
					return fmt.Errorf("could not submit post-login form: %w", err)
				}
				continue
			}
		}

		doc, err = c.loadHTML(c.submitForm(doc, formID))
		if err != nil {
			return fmt.Errorf("could not submit %s form: %w", formID, err)
		}
	}

	action, err = c.buildAction("index.jsp")
	if err != nil {
		return err
	}

	_, err = c.loadHTML(c.get(action))
	if err != nil {
		return fmt.Errorf("could not load mainPage: %w", err)
	}

	err = c.postLogin()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) buildAction(action string) (string, error) {
	c.log.Println(c.refererURL)
	if !strings.HasPrefix(action, "http://") && !strings.HasPrefix(action, "https://") {
		if strings.HasPrefix(action, "/") {
			return c.baseURL + action, nil
		}

		parent := c.refererURL[:strings.LastIndex(c.refererURL, "/")]
		return parent + "/" + action, nil
	}

	actionurl, err := url.Parse(action)
	if err != nil {
		return "", fmt.Errorf("could not parse action url: %w", err)
	}
	return actionurl.String(), nil
}

func (c *client) submitForm(doc *html.Node, id string) (*http.Response, error) {
	form := htmldom.GetElementByID(doc, id)
	if form == nil {
		html.Render(os.Stderr, doc)
		/*
			<html><head></head><body>
			<form name="openTop" id="openTop" action="/mua/initAuthProcess" method="GET" target="_top">

			</form>


			<script type="text/javascript">
			    var form = document.getElementById("openTop");
			    form.submit();
			</script>


			</body></html>
		*/
		return nil, fmt.Errorf("form not found: %s", id)
	}
	action := getAttribute(form, "action")
	return c.submitFormValues(doc, id, action, nil)
}

func (c *client) submitFormValues(doc *html.Node, id, action string, values url.Values) (*http.Response, error) {
	form := getElementByID(doc, id)
	fields := parseFormFields(form)
	replaceValues(fields, values)

	u, err := c.buildAction(action)
	if err != nil {
		return nil, err
	}

	return c.postForm(u, fields)
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

/*
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
*/

func (c *client) preGetDepositsBalance() error {
	resp, err := c.request(http.MethodGet, c.baseURL+c.cstUrl("/cb/pages/jsp-ns/olb/InitAccountSummary?redirect=ALLACCOUNTS_HOME&type=ALLACCOUNTS_HOME"), nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil

}

func (c *client) Logout() error {
	resp, err := c.request(http.MethodGet, c.baseURL+c.cstUrl("/cb/pages/jsp-ns/olb/SafeExit"), nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}

type DepositsBalance struct {
	AccType          string `json:"accType"`
	AvailableBalance string `json:"availableBalance"`
	Currency         string `json:"currency"`
	Description      string `json:"description"`
	ID               string `json:"id"`
	NickName         string `json:"nickName"`
	Number           string `json:"number"`
	ProductName      string `json:"productName"`
	Type             string `json:"type"`
}

func (c *client) GetDepositsBalance() (DepositsBalance, error) {
	err := c.preGetDepositsBalance()
	if err != nil {
		return DepositsBalance{}, err
	}

	u := fmt.Sprintf(`%s/cb/pages/jsp/account/getDepositsBalanceBancolombiaHome.action?cst=%s`, c.baseURL, c.cst)
	bodyStr := fmt.Sprintf(`type=DEPOSITS&CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)

	var response struct {
		JSON                   string            `json:"JSON"`
		BanHidenProduct        bool              `json:"banHidenProduct"`
		ExistFiduciaria        bool              `json:"existFiduciaria"`
		ExistVirtualInvestment bool              `json:"existVirtualInvestment"`
		GridModel              []DepositsBalance `json:"gridModel"`
		Page                   int64             `json:"page"`
		Records                int64             `json:"records"`
		Rows                   int64             `json:"rows"`
		Sord                   string            `json:"sord"`
		Total                  int64             `json:"total"`
	}

	err = c.requestJSONWithHeaders(http.MethodPost, u, body, &response, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if err != nil {
		return DepositsBalance{}, err
	}

	return response.GridModel[0], nil
}

func (c *client) preGetSavingsDetail(id string, step int) error {
	var u string
	if step == 1 {
		u = c.baseURL + c.cstUrl(fmt.Sprintf("/cb/pages/jsp-ns/olb/ACCTARGETQuery?entity=MOVCA&fwviejoId=%s&operation=MOVCA&clean=true", id))
	} else {
		u = c.baseURL + c.cstUrl(fmt.Sprintf("/cb/pages/jsp-ns/olb/AccountDetailAsset?&step=%d&open=Y", step))
	}

	resp, err := c.request(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	return nil
}

type SavingsDetail struct {
	Amount      float64 `json:"amount"`
	BranchID    string  `json:"branchId"`
	Date        fecha   `json:"date"`
	Description string  `json:"description"`
	OptionalRef string  `json:"optionalRef"`
}

type fecha struct {
	time.Time
}

func (f *fecha) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	f.Time, err = time.Parse(`"2006-01-02T15:04:05"`, string(data))
	return err
}

func getND() string {
	now := time.Now()
	return strconv.FormatInt(now.UnixNano(), 10)[:13]
}

func (c *client) GetSavingsDetail(id string, page int) ([]SavingsDetail, error) {
	err := c.preGetSavingsDetail(id, page)
	if err != nil {
		return nil, err
	}

	nd := getND()

	u := fmt.Sprintf(`%s/cb/pages/jsp/account/getSavingsDetailAction.action?CSRF_TOKEN=%s`, c.baseURL, c.csrfToken)
	bodyStr := fmt.Sprintf(`_search=false&nd=%s&rows=-1&page=1&sidx=date&sord=desc`, nd)
	body := strings.NewReader(bodyStr)

	var response struct {
		JSON      string          `json:"JSON"`
		GridModel []SavingsDetail `json:"gridModel"`
		Page      int64           `json:"page"`
		Records   int64           `json:"records"`
		Rows      int64           `json:"rows"`
		Sidx      string          `json:"sidx"`
		Sord      string          `json:"sord"`
		Total     int64           `json:"total"`
	}

	err = c.requestJSON(http.MethodPost, u, body, &response)
	if err != nil {
		return nil, err
	}

	return response.GridModel, nil
}

func encodeUnescaped(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	for k, vs := range v {
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func (c *client) AccountEnrroll() error {
	var err error

	// err = c.getPaymentDetail(CreditCard)
	// if err != nil {
	// 	return err
	// }

	// err = c.getPaymentDetail(Credits)
	// if err != nil {
	// 	return err
	// }

	err = c.messageCenterLoadConfigurationCIC()
	if err != nil {
		return err
	}

	err = c.messageCenterLanding()
	if err != nil {
		return err
	}

	err = c.getMessagesForceCIC()
	if err != nil {
		return err
	}

	urlMenu := "/cb/pages/jsp/account/enrrollProduct_invoke_from_menu.jsp"

	u := fmt.Sprintf(`%s/cb/pages/jsp/security/invokeSecondPass.jsp?cst=%s`, c.baseURL, c.cst)
	bodyValues := url.Values{}
	bodyValues.Add("urlreturn", urlMenu)
	bodyValues.Add("menu", "TRANSFER")
	bodyValues.Add("sub", "TIC")
	bodyValues.Add("wizard", "N")
	bodyValues.Add("CSRF_TOKEN", c.csrfToken)
	bodyValues.Add("cst", c.cst)
	body := strings.NewReader(encodeUnescaped(bodyValues))

	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	u = fmt.Sprintf(`%s/cb/pages/jsp-ns/olb/GetOtpServicesInfo?cst=%s`, c.baseURL, c.cst)
	body = strings.NewReader(fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst))

	doc, err := c.loadHTML(c.requestWithHeaders(http.MethodPost, u, body, headers))
	if err != nil {
		return err
	}

	err = c.invokeActionError("")
	if err != nil {
		return err
	}

	err = c.messageCenterMenu()
	if err != nil {
		return err
	}

	err = c.getMessagesForceCIC()
	if err != nil {
		return err
	}

	u = fmt.Sprintf(`%s/mua/initAuth?cst=%s`, c.baseURL, c.cst)
	token := parseTokenValue(doc)
	bodyValues = url.Values{"TOKEN": []string{token}}
	body = strings.NewReader(bodyValues.Encode())
	doc, err = c.loadHTML(c.requestWithHeaders(http.MethodPost, u, body, headers))
	if err != nil {
		return err
	}

	for {
		fmt.Print("OTP: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal("OTP required")
		}
		fmt.Println()
		otp := string(bytePassword)

		// autenticationOdaForm := getElementByID(doc, "autenticationOdaForm")
		// if autenticationOdaForm == nil {
		// 	html.Render(os.Stderr, doc)
		// 	return fmt.Errorf("could not find login user form")
		// }

		t1Assertion := parseT1Assertion(doc)
		// keyboardSrc := parseKeyboardContent(doc)

		// keyboardNode, err := html.Parse(strings.NewReader(keyboardSrc))
		// if err != nil {
		// 	return fmt.Errorf("could not parse keyboard: %w", err)
		// }
		// keyMap := parseKeyboardMap(keyboardNode)
		// fmt.Println(otp)
		// otp = mapPassword(keyMap, otp)
		// fmt.Println(otp)
		idSs := c.rsa.processPassword(otp, t1Assertion)
		fmt.Println(idSs)
		action := parseAjaxRequestURL(doc)
		// action := getAttribute(autenticationOdaForm, "action")
		// fmt.Println(action)
		// errorTexto
		bodyValues = url.Values{
			"softoken":   []string{idSs},
			"tempUserID": []string{""},
			"userID":     []string{""},
		}
		doc, err = c.loadHTML(c.postForm(c.baseURL+action, bodyValues))
		if err != nil {
			return fmt.Errorf("could not submit autenticationOdaForm: %w", err)
		}

		token = parseTokenValue(doc)
		if token != "" {
			break
		}

		summary := getElementByID(doc, "summary")
		errorText := strings.TrimSpace(getInnerText(summary))
		if errorText != "" {
			fmt.Fprintln(os.Stderr, errorText)
			continue
		}

		return fmt.Errorf("unexpected error")
	}

	body = strings.NewReader(`TOKEN=` + url.QueryEscape(token))
	u = fmt.Sprintf("%s/cb/pages/jsp/security/sessionInvoke.jsp?CSRF_TOKEN=%s&menu=TRANSFER&sub=TIC&urlMenu=%s&security_domain=OTP&cst=%s",
		c.baseURL, c.csrfToken, "/cb/pages/jsp/security/sessionInvoke.jsp", c.cst)

	doc, err = c.loadHTML(c.requestWithHeaders(http.MethodPost, u, body, headers))
	if err != nil {
		return err
	}
	tokenMada := parseTokenMada(doc)
	// tokenMada, err := url.QueryUnescape(parseTokenMada(doc))
	// if err != nil {
	// 	return fmt.Errorf("could not unescape mada token: %w", err)
	// }
	u = fmt.Sprintf(`%s/cb/pages/jsp-ns/validateSecondPasswordFromMada.action?cst=%s`, c.baseURL, c.cst)
	bodyValues = url.Values{}
	// bodyValues.Set("tokenMada", tokenMada)
	bodyValues.Set("menu", "TRANSFER")
	bodyValues.Set("sub", "TIC")
	// bodyValues.Set("urlMenu", `/cb/pages/jsp/account/enrrollProduct_invoke_from_menu.jsp`)
	bodyValues.Set("flowid", "null")
	bodyValues.Set("security_domain", "OTP")
	bodyValues.Set("CSRF_TOKEN", c.csrfToken)
	// bodyValues.Set("cst", c.cst)
	body = strings.NewReader(`tokenMada=` + tokenMada + `&` + bodyValues.Encode() +
		`&urlMenu=` + urlMenu + `&cst=` + c.cst)

	resp, err = c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	u = fmt.Sprintf(`%s/cb/pages/jsp/account/enrrollProduct_invoke_from_menu.jsp`, c.baseURL)
	bodyValues = url.Values{}
	bodyValues.Set("flowid", "null")
	bodyValues.Set("security_domain", "OTP")
	// bodyValues.Set("CSRF_TOKEN", c.csrfToken)
	// bodyValues.Set("cst", c.cst)
	body = strings.NewReader(bodyValues.Encode() + fmt.Sprintf(`&CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst))
	resp, err = c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	u = fmt.Sprintf(`%s/cb/pages/jsp/account/EnrollmentGetListBanksACHAction.action?cst=%s`, c.baseURL, c.cst)
	body = strings.NewReader(fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst))
	doc, err = c.loadHTML(c.requestWithHeaders(http.MethodPost, u, body, headers))
	if err != nil {
		return err
	}

	_ = html.Render(os.Stderr, doc)

	// -----------
	if productTypes, err := c.GetProductTypes(); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("productTypes: %v\n", productTypes)
	}

	if productTypes, err := c.GetProductTypes(); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("productTypes: %v\n", productTypes)
	}

	if bankGroups, err := c.GetBankGroups(); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("bankGroups: %v\n", bankGroups)
	}

	if documentTypes, err := c.GetDocumentTypes(); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("documentTypes: %v\n", documentTypes)
	}

	if banks, err := c.GetBanks(); err != nil {
		log.Println(err)
	} else {
		fmt.Printf("banks: %v\n", banks)
	}

	// -----------

	/*
		 curl -i -s -k -X  'POST' -H 'Content-Length: 103'   -H 'Content-Type: application/x-www-form-urlencoded'  -H 'Cookie: JSESSIONID=PfMywXIH0jEfMB6+BJ1lHsG0; NSC_JOr2zhh2e44kdkqd4uupeqdgxr1z1c0=28d4a3da96b41259453c630389e9f27b968597ccb6f7ebc43f41db64e6ee8f12d04574f6; __cflb=02DiuF7aX6zsQEVJrpNgCqfZ7XAJa8kSz4UWq2yrdtctt'  \
			--data-binary $'CSRF_TOKEN=4098441916498417831&cst=H86r%2BYTjlgt2XGAC2yNdSSxWuSNa8lGiNXgw8aPBqRpKYQjd%2FwQLzG3aXJnrbwGy' \
			'https://sucursalpersonas.transaccionesbancolombia.com/cb/pages/jsp-ns/olb/EnrollACHAccountConfirmation'
	*/
	u = fmt.Sprintf(`%s/cb/pages/jsp/account/VerifyEnrollment.action?cst=%s`, c.baseURL, c.cst)
	bodyValues = url.Values{}
	bodyValues.Set("bankToId", "0007")
	bodyValues.Set("accountType1", "0")
	bodyValues.Set("accountType2", "7")
	bodyValues.Set("fiduciariaTypes", "0")
	bodyValues.Set("accountNumber", "99999999999")
	bodyValues.Set("accountNickname", "NickName")
	bodyValues.Set("documentType", "CC")
	bodyValues.Set("documentNumber", "111111111111111")
	bodyValues.Set("accountType", "7")
	bodyValues.Set("productName", "")

	// bodyValues.Set("devicePrint", "")
	// bodyValues.Set("CSRF_TOKEN", c.csrfToken)
	// bodyValues.Set("cst", c.cst)
	body = strings.NewReader(bodyValues.Encode() + fmt.Sprintf(`&CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst))
	doc, err = c.loadHTML(c.requestWithHeaders(http.MethodPost, u, body, headers))
	if err != nil {
		return err
	}

	if summary := parseJQuerySummary(doc); summary != "" {
		fmt.Println(summary)
	}

	// _ = html.Render(os.Stderr, doc)
	// c.log.Println("waiting 3 seconds")
	// <-time.After(time.Second * 3)
	// c.log.Println("continuing...")

	u = fmt.Sprintf(`%s/cb/pages/jsp-ns/olb/EnrollACHAccountConfirmation?cst=%s`, c.baseURL, c.cst)
	body = strings.NewReader(fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst))
	h := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/x-www-form-urlencoded",
		// "Origin":          "https://sucursalpersonas.transaccionesbancolombia.com",
		// "Referer":         "https://sucursalpersonas.transaccionesbancolombia.com/cb/pages/jsp/home/index.jsp",
	}
	doc, err = c.loadHTML(c.requestWithHeaders(http.MethodPost, u, body, h))
	if err != nil {
		return err
	}

	_ = html.Render(os.Stderr, doc)

	if summary := parseJQuerySummary(doc); summary != "" {
		fmt.Println(summary)
	}

	return nil
}

// mapPassword encodes the actual password using the keymap
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

func (c *client) GetBankGroups() (GetBankGroupsResponse, error) {
	var resp GetBankGroupsResponse

	u := c.baseURL + `/cb/pages/jsp/account/GetBankGroups.action`
	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)

	return resp, err
}

func (c *client) GetProductTypes() (GetProductTypesResponse, error) {
	var resp GetProductTypesResponse

	u := c.baseURL + `/cb/pages/jsp/account/GetProductTypes.action`
	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)

	return resp, err
}

func (c *client) GetDocumentTypes() (GetDocumentTypesResponse, error) {
	var resp GetDocumentTypesResponse

	u := c.baseURL + `/cb/pages/jsp/account/GetDocumentTypes.action`
	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)

	return resp, err
}

func (c *client) GetBanks() (GetBanksResponse, error) {
	var resp GetBanksResponse

	u := c.baseURL + `/cb/pages/jsp/account/GetBanks.action`
	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)

	return resp, err
}

func (c *client) invokeActionError(code string) error {
	var resp map[string]interface{}
	u := c.baseURL + `/cb/pages/jsp/errors/ConfigFunctionalityErrorCode.action?cst=` + c.cst

	bodyStr := fmt.Sprintf(`err=%s&CSRF_TOKEN=%s&cst=%s`,
		url.QueryEscape(code), c.csrfToken, c.cst)

	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)
	c.log.Printf("resp: %v\n", resp)
	return err
}

func (c *client) messageCenterLoadConfigurationCIC() error {
	var resp map[string]interface{}
	u := c.baseURL + `/cb/pages/jsp/messagecenter/MessageCenterLoadConfigurationCIC.action?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)
	c.log.Printf("resp: %v\n", resp)
	return err
}

func (c *client) messageCenterLanding() error {
	u := c.baseURL + `/cb/pages/jsp/messagecenter/MessageCenterLanding.jsp?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	return nil
}

func (c *client) messageCenterMenu() error {
	u := c.baseURL + `/cb/pages/jsp/messagecenter/MessageCenterMenu.jsp?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	return nil
}

func (c *client) getMessagesForceCIC() error {
	var resp map[string]interface{}
	u := fmt.Sprintf(`%s/cb/pages/jsp/messagecenter/GetMessagesForceCIC.action?CSRF_TOKEN=%s&cst=%s`,
		c.baseURL, c.csrfToken, c.cst)

	nd := getND()
	bodyStr := fmt.Sprintf(`_search=false&nd=%s&rows=0&page=1&sidx=receptionDate&sord=desc`, nd)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)
	if err != nil {
		return err
	}

	c.log.Printf("resp: %v\n", resp)

	return nil
}

func (c *client) gATokenGeneration(idGA string) error {
	var resp map[string]interface{}
	u := c.baseURL + `/cb/pages/jsp/ga/GATokenGeneration.action`

	bodyStr := fmt.Sprintf(`id_ga=%s&timeZone=GMT-0500&cst=%s&CSRF_TOKEN=%s`, idGA, c.cst, c.csrfToken)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)
	if err != nil {
		return err
	}

	c.log.Printf("resp: %v\n", resp)

	return nil
}

func (c *client) getUpdateDynamicData() error {
	var resp map[string]interface{}
	u := c.baseURL + `/cb/pages/jsp/updateData/getUpdateDynamicData.action`

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)
	c.log.Printf("resp: %v\n", resp)
	return err
}

func (c *client) preApprovedAction() error {
	var resp map[string]interface{}
	u := c.baseURL + `/cb/pages/jsp-ns/olb/PreApprovedAction?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	err := c.requestJSONWithHeaders(http.MethodPost, u, body, &resp, headers)
	c.log.Printf("resp: %v\n", resp)
	return err
}

func (c *client) usersPreferencesAction() error {
	u := c.baseURL + `/cb/pages/jsp-ns/olb/UsersPreferencesAction`

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

func (c *client) userSecondPass() error {
	u := c.baseURL + `/cb/pages/jsp/notification/userSecondPass.action?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

func (c *client) personalizeProductsNameDataEntry() error {
	u := c.baseURL + `/cb/pages/jsp-ns/olb/PersonalizeProductsNameDataEntry?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

func (c *client) changeNameACHDataEntry() error {
	u := c.baseURL + `/cb/pages/jsp-ns/olb/ChangeNameACHDataEntry?cst=` + c.cst

	bodyStr := fmt.Sprintf(`CSRF_TOKEN=%s&cst=%s`, c.csrfToken, c.cst)
	body := strings.NewReader(bodyStr)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp, err := c.requestWithHeaders(http.MethodPost, u, body, headers)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

type PaymentType string

const (
	CreditCard PaymentType = "CREDIT_CARD"
	Credits    PaymentType = "CREDITS"
)

func (c *client) getPaymentDetail(paymentType PaymentType) error {
	var resp map[string]interface{}
	u := fmt.Sprintf(`%s/cb/pages/jsp/account/GetPaymentDetailAction.action?type=%s&cst=%s`, c.baseURL, paymentType, c.cst)

	headers := map[string]string{"Accept": "*/*"}
	err := c.requestJSONWithHeaders(http.MethodGet, u, nil, &resp, headers)
	c.log.Printf("resp: %v\n", resp)
	return err
}

func (c *client) postLogin() error {
	var err error
	err = c.changeNameACHDataEntry()
	if err != nil {
		return err
	}

	err = c.personalizeProductsNameDataEntry()
	if err != nil {
		return err
	}

	err = c.userSecondPass()
	if err != nil {
		return err
	}

	err = c.usersPreferencesAction()
	if err != nil {
		return err
	}

	err = c.preApprovedAction()
	if err != nil {
		return err
	}

	err = c.getUpdateDynamicData()
	if err != nil {
		return err
	}

	err = c.gATokenGeneration("anuncio-0")
	if err != nil {
		return err
	}
	err = c.gATokenGeneration("anuncio-1")
	if err != nil {
		return err
	}

	return nil
}
