package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestParseT1Assertion(t *testing.T) {
	f, err := os.Open("./password.html")
	require.NoError(t, err)
	defer f.Close()

	doc, err := html.Parse(f)
	require.NoError(t, err)

	actual := parseT1Assertion(doc)
	assert.Equal(t, "ozzutFJrd0LBW45xy1kS", actual)
}

func TestParseKeyboardContent(t *testing.T) {
	f, err := os.Open("./password.html")
	require.NoError(t, err)
	defer f.Close()

	doc, err := html.Parse(f)
	require.NoError(t, err)

	actual := parseKeyboardContent(doc)
	assert.Equal(t, `  <table class='keyboard' border='0' cellspacing='0' cellpadding='0' align='left' valign='top'>  <tr>    <td width='0' height='37' ></td>    <td></td>  </tr>  <tr>    <td height='0' width='2'></td>    <td colspan='2'>      <table align='left' valign='top' cellspacing='0' cellpadding='0' class='bg_button'>        <tr align='left'>                  <td valign='top' align='left'> <table class='bg_button' id='_KEYBRD' valign='top' >  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"LONk\");'>  <div border='0' id ='taQpCXWNIycL9' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>9</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"UMjS\");'>  <div border='0' id ='taQpCXWNIycL1' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>1</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"lKFE\");'>  <div border='0' id ='taQpCXWNIycL0' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>0</div></td></tr>  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"WTQa\");'>  <div border='0' id ='taQpCXWNIycL7' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>7</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"ZJDn\");'>  <div border='0' id ='taQpCXWNIycL8' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>8</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"igA6\");'>  <div border='0' id ='taQpCXWNIycL5' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>5</div></td></tr>  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"YXfV\");'>  <div border='0' id ='taQpCXWNIycL6' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>6</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"974C\");'>  <div border='0' id ='taQpCXWNIycL2' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>2</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"5bG8\");'>  <div border='0' id ='taQpCXWNIycL4' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>4</div></td></tr>  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"eHd3\");'>  <div border='0' id ='taQpCXWNIycL3' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>3</div></td><td colspan='2' onclick='clearKeys();' class='bg_buttonSmall'><div id='clearKey' border='0' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>Borrar</div></td></tr></table><table class='bg_button' id='_CONSTRAST' valign='top' cellspacing='0'>  <tr><td><img width='90' height='34' border='0' src='/mua/images/kb/Contraste" + contrastLevel + ".gif?v=4.1.1.RC2_1622257216208' name='constrastImg' id='constrastImg' usemap='#numericKeyboardMap' > <map name='numericKeyboardMap' id='numericKeyboardMap'><area shape='circle' class='cursorContrast' coords='10,30,15' onmouseover=setHandCursor(document.constrastImg) onclick='changeContrastLevel(1)' onmouseout='setDefaultCursor(document.constrastImg)'><area shape='circle' class='cursorContrast' coords='50,30,15' onmouseover=setHandCursor(document.constrastImg) onclick='changeContrastLevel(2)' onmouseout='setDefaultCursor(document.constrastImg)'><area shape='circle' class='cursorContrast' coords='90,30,15' onmouseover=setHandCursor(document.constrastImg) onclick='changeContrastLevel(3)' onmouseout='setDefaultCursor(document.constrastImg)'></map></td></tr></table></td>        </tr>      </table>    </td>  </tr><tr>    <td height='17'></td>    <td colspan='2'></td>  </tr> </table>`, actual)
}

func TestParseKeyboardMap(t *testing.T) {
	src := `  <table class='keyboard' border='0' cellspacing='0' cellpadding='0' align='left' valign='top'>  <tr>    <td width='0' height='37' ></td>    <td></td>  </tr>  <tr>    <td height='0' width='2'></td>    <td colspan='2'>      <table align='left' valign='top' cellspacing='0' cellpadding='0' class='bg_button'>        <tr align='left'>                  <td valign='top' align='left'> <table class='bg_button' id='_KEYBRD' valign='top' >  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"LONk\");'>  <div border='0' id ='taQpCXWNIycL9' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>9</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"UMjS\");'>  <div border='0' id ='taQpCXWNIycL1' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>1</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"lKFE\");'>  <div border='0' id ='taQpCXWNIycL0' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>0</div></td></tr>  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"WTQa\");'>  <div border='0' id ='taQpCXWNIycL7' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>7</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"ZJDn\");'>  <div border='0' id ='taQpCXWNIycL8' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>8</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"igA6\");'>  <div border='0' id ='taQpCXWNIycL5' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>5</div></td></tr>  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"YXfV\");'>  <div border='0' id ='taQpCXWNIycL6' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>6</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"974C\");'>  <div border='0' id ='taQpCXWNIycL2' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>2</div></td><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"5bG8\");'>  <div border='0' id ='taQpCXWNIycL4' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>4</div></td></tr>  <tr><td class='bg_buttonSmall'  align='center' style='cursor:default' onMouseOver='cwavtLjsmzYX();' onmouseout='changeToOrigKeyboard();' onclick='gyxdNaciylcN(\"eHd3\");'>  <div border='0' id ='taQpCXWNIycL3' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>3</div></td><td colspan='2' onclick='clearKeys();' class='bg_buttonSmall'><div id='clearKey' border='0' valign='center' align='center' onfocus='this.blur();' class='colorContrast + contrastLevel + '>Borrar</div></td></tr></table><table class='bg_button' id='_CONSTRAST' valign='top' cellspacing='0'>  <tr><td><img width='90' height='34' border='0' src='/mua/images/kb/Contraste" + contrastLevel + ".gif?v=4.1.1.RC2_1622257216208' name='constrastImg' id='constrastImg' usemap='#numericKeyboardMap' > <map name='numericKeyboardMap' id='numericKeyboardMap'><area shape='circle' class='cursorContrast' coords='10,30,15' onmouseover=setHandCursor(document.constrastImg) onclick='changeContrastLevel(1)' onmouseout='setDefaultCursor(document.constrastImg)'><area shape='circle' class='cursorContrast' coords='50,30,15' onmouseover=setHandCursor(document.constrastImg) onclick='changeContrastLevel(2)' onmouseout='setDefaultCursor(document.constrastImg)'><area shape='circle' class='cursorContrast' coords='90,30,15' onmouseover=setHandCursor(document.constrastImg) onclick='changeContrastLevel(3)' onmouseout='setDefaultCursor(document.constrastImg)'></map></td></tr></table></td>        </tr>      </table>    </td>  </tr><tr>    <td height='17'></td>    <td colspan='2'></td>  </tr> </table>`
	doc, err := html.Parse(strings.NewReader(src))
	require.NoError(t, err)

	actual := parseKeyboardMap(doc)
	assert.Equal(t, map[string]string{
		"0": "lKFE",
		"1": "UMjS",
		"2": "974C",
		"3": "eHd3",
		"4": "5bG8",
		"5": "igA6",
		"6": "YXfV",
		"7": "WTQa",
		"8": "ZJDn",
		"9": "LONk",
	}, actual)
}

func TestParsePasswordInputName(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./password.html",
			expected: "uvdEMkTtiXlW",
		},
		{
			filename: "./testdata/response2.html",
			expected: "yhgsuYXvJQii",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {

			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parsePasswordInputName(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseUrlRedirect(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response4.html",
			expected: "/cb/pages/jsp-ns/login-mada.jsp",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseUrlRedirect(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseTokenMua(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response4.html",
			expected: "gt9Evelcvl9yQeulJn8MGyvH3Sr6FcKNnmtAwEyW/CU0DE9O96EqxgxtTJtjzr0EIQ2IFD37maPMhVNw5Lh2Z4glHUe0MR+lXTtTCUj8dN/uJCgggsnJSfPyuYbPs45/GtnSwIKlI8k9EWk4C4/ZKRBE+3Ny5KcWFHZJEvVjS8bC67HObBGzks7E5Iv7lqg2pzQxhbjFzG6dcDqTCZQ8GgA6ZUf/HTef75ya8JVOGee2dON3CJIzs5Y2fQ++DKKykVdY22wOUo2Vt6ohGH0A6vZhsnML/xtk5pKkEqWnXjwdEQ56wwfT+rUdKNiLzjKEC4ZpSVTheppdIL2nY4A0tCqnZaRn9LS7/FTTRPa1t8U1NfnSbFjii61ZowCMKbluRdIc9FSBT0pieYSi18h6CFTLieVVG8ly6mRiVN9GKm5nKGh1YOYRZrHRXIV4EkQAHbu9bP7QuP9Pyp1X6+3RLd2AEukvHa99xEpl35ossCnlLIzNeXzIkCr7iJqhV7eIpWOCSC+HSfuWD0O/l4RA0OAsNi+yjuK33ssbZLbYQIDWWXVqdMNyoi58kp4IGuDsHlCnPxqu4Czrf1LRD03S0UX16GDRraQNy6tgvEv/piucQVNL3CqPptnytkkP45g8fnMxoew6wMpIaaFSCCwP1mVTcPnyP8sy6MzqqHDbjZTH6nX2p+qwh922uf/0zD4LVwXTaaLQkGLQRZsBN0PpaB5xPYMpFKd9WWVnviFqkNKsymRsW09KiTWNJ6lIie90LvE5j+ZDDEKBJi/sZ1kO8Lmr1fDTjhEmhjEgSa2XBQL4Puvt2NTDwNm5tHz1lJLQWMZAw2vbq7kDuiQ02S5PlLEJgCxk6OMPQQm2pnOqxxj62UqPDLHACg6ZO6enqTX4bS4UCR/xZaiHcNmvMlZ50YINqVAk31claSfl4KZRnKGfQIUYfOHY3Wm66bHkJlQrood1ql/7MoRidkdVqILQdBh2aA8W+fK9jo2fvtCHyESkc3JBi8cZi0AyYJeJTq0n5ifFj2bebQOSUPHh2pOXU8P9OxsMG9SoapO5YSUdkVoyaMhqzVRCSasfFijV4NToAXgUu+6mvtHh/AnQ1I1ImclMMA4oiQdBXNL/GYpSxynfa+yArzRJWs24tlVI6wfgwLGV2x7aZSuWpyuiousQi2wsCjCoHpp45bApIf11Y8l9U+OkkrpdUqHeUODi5SOfJnEmyQKEWQPhmB8P/Aueu91SNuifOh16AayBjxWGDblIEnFJgto+TrTYZIEtcSgW0RChmZEyPyT/m8PaQPxLrXsRWl4LlUg5uGpSdDNTVZO7dQm8JtJ27kLXI2nqV6kIk+K1pC6Tr0yzrXzsskz5hjWi/oru55H3sz3GD3zcFxkRc9kN//POpl1g9ABwg3DcwRxkMtVrmAhh8ycA2TM+UndHC8Zing+6RJgt37ZsZ3onbgGlq6p1J3l/wMDgj/xebnJt83hH/VF9wow8Yt+NIUf4YbM=",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseTokenMua(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseCodeRedirect(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response4.html",
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseCodeRedirect(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseLocationReplace(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response5.html",
			expected: "/cb/pages/jsp/home/mainPage.jsp",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseLocationReplace(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseCstParam(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response6.html",
			expected: "A%2BCfdsj9tJ2aCxHoQFFGAmYGRYWBTPX2hsWlYdZ0AM8OqnP4UG0TUMeCOiE3PXz3",
		},
		{
			filename: "./testdata/response7.html",
			expected: "9YZU648yLoJNNc%2FHBwjNZrZLdNsV%2BgK6CKYl92s4XPnBEYKfMgRpiEDw3T09%2BmCO",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseCstParam(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseCsrfToken(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response8.html",
			expected: "1992276304611140186",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseCsrfToken(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseTokenValue(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response9.html",
			expected: "TOKEN_VALUE",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseTokenValue(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseTokenMada(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response10.html",
			expected: "UJi8r94NHGDTXFcs2WmJpQLDEByYuEV8DQc%2BiCVVtA67hHHq3LvjHrXQ10Kbf%2BpCoXvp38hMDK4s15ipbjSuPaMjrXwwIrfR90GU2sDFsZA%2FumLNtmYIzYTDRZaIZAT%2BuUZBdVUY5wQ8S%2BOo2Zm110rNC5u3%2BqKn6Bhfhrlma6d6tsQbmDgyRc6pDXYQompvTVzGv7bAKLTlYAcvpp%2Fr6mCscMNzjTVF2wy4r6n6oN6QMXGptaVqS38oPivaf0Fa4v%2BKLBL8xPRhNhUPpLPCRHK%2Bh0AJibGZYLsnzupM7A0839SVVfY3uhmpk43Fm0Cnf%2Fnjc4juN8OEi30GlJwyY5jnzcWMQogsZITHmXH2MSq074FvM3sKQCU6Tuku24AC93Hq58dF0KkI0mXvM3Ojo5n%2FdrrHN8mrp8HBAtlGUiZK5rdPN%2B5F3SOAHuC2Nd73WWwaxUkbvIMsCQ9NL4A70QzvYaw1YoSEnvt2nppXG9b87lK95e39KsVg9yukCQ0sb%2BF2xJ9%2FnKrsbPKwgX0vZMPDVZaFTt0Wx5MzRuHZ04PnhU3c0jyyzC9IW0Aq1Z6fs6eCRvgW4%2FUOL70RnUJ8XAqdaintqQdVBoSeQKTnehoK6wGmxwowoES2kPJWVbDklMMR6%2BoS2MtmPEpiI5Eh66WIMli1erYiEYoJ1FCTwnsXWgBgoGS9v2Bq1%2B0i5%2Ffb9ETkbUCFx%2Fy3AGxtE%2FRaMZaCVJKwpnZS10kZLRtHerR2t%2FF2c7rHC94BaBeqP0rJSPBEc4WO8Bp%2FaHQQVCc8TLLhrx3uDo9y8j2M0loW7DmB82mSvPfjlEuPfSiV%2B2RCFfPhQcs0ynWBtoQU6f2WlW%2Br8srkwH37VxMC%2FE8fALPF8uWalSVMGjvXxgTy%2FFSOkGYz2noLNYUNpY8Ws1NbQR7jq9X8pQmsFoXa8iLiJBYqKwqjGsu3fK2V3jGKAtyzIyHcFDgh9xULYTPmV6ZbGNJ426AolrZE7XlMv120Rlh4DQ2BLFHtu09tMjORSPbrb77qA3YZeaiXGzLSVDIxXcPCiyBE69DVfKh2Mj05Gr%2B2sRBaJ%2FjebC3ubPTepmadvtOaLmyk8LM4rJztN34oo9V0qQvsLni1KnnzUaZxStqZPbyNKiHnWRlYg0s6ywQEzDm0Xv1IGa1Q67K2Y5lMdSWLPWggyBN25NWu0pf6kPXM8ooz%2BAOYgEvsvaZTiv39HCS%2BsFoKg16X3T0Sb4pQPmBIyneyH4LI6LIYgyokrFoprDb8BDjL6z45rfELODx8BPajcjUmTwKmmdGsmMzG9i%2BQW%2FRZ%2FtrZMAwUU%2B%2F51kT%2FI0ltp9aBkz10fwcH0dxglclEw39dj4zX1qLfvCtMaljxvVIFZoapJ%2FGLwUuVo%2FicTuJ6388acsnU8ETl7oU8LuxB%2FxrE4%2BmTEHj3hvsfD%2BsE6eZ1gPVcJbUQJ2nFkjIs0X8A%2FnCfM2podC%2FIbji1Lh74PBIEAkniLEehFnD25r1RzcjIK3YSqi3%2FLryjYRgHAirVdIsvpxJCFeNAh7jt5a6tAfiWvDlJHoxFtO%2BQVtEabDdcoMRat4hQ%2Fp%2Fu%2BSRsTDWL%2FM3GH5DOANXxV6PXmnrAZ2uYthnPsGkjChkeJLAi94etS6dh1sGZ52A6JdONxg2lR7MHIh4kO36ozraXnhWmYX1PX%2BiPG1V21v1QCrbof%2BVXLTTTtLFtw14MDmJ542GMO5pISFl8PHkKiJQJzvkRvTH7mmpRL976%2FfTHmnzTqK%2F5ZQGx6KAD0BFaMZIejTton53FzPEaJc3Rv1hPtQueV%2FA9bB6ydg8z%2FjYeeR4wgvTsC25A%2BhNpBHXSFyZgTUYDnXRv1NUDWn5I1kXLuPa6b8OUHItA4H5emgzJBCfFWgXp0nVGnxByDm%2BuflOMlauSiQm4UhhjVe8wKJH3noHbRSHhK9TEQ%2B90a7TawKpDk6SYBKhJcA09HFd7iS8Pzjp%2Ff2IKYqHsDh6C5oftjrJjp2fvwFYvUzK1kzwLXAEVFkTpBzburK2gbg0N1W1bGDi132ABPq2ahubb7LV3tkP3h6FCr0JuZN0b4KPUEevGNHIl4i%2FvTlvGS8RHKze7LvjfOOICWPbwI6TPcqYZxgkJKQ%3D%3D",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseTokenMada(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseJQuerySummary(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{
			filename: "./testdata/response11.html",
			expected: "Cuenta No existe en Depositos CODIGO:BC 917",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f, err := os.Open(tt.filename)
			require.NoError(t, err)
			defer f.Close()

			doc, err := html.Parse(f)
			require.NoError(t, err)

			actual := parseJQuerySummary(doc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
