package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDevicePrint(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "base case",
			input:  `{"Version":"3.5.1_4","Browser":{"userAgent":"mozilla/5.0 (x11; linux x86_64; rv:90.0) gecko/20100101 firefox/90.0","appVersion":"5.0 (X11)","platform":"Linux x86_64","appMinorVersion":"","cpuClass":"","browserLanguage":"","browserName":"Firefox","browserVersion":"90.0","browserMajor":"90","browserEngineName":"Gecko","browserEngineVersion":"90.0","osName":"Linux","browserOS":"Linux","osVersion":"x86_64","deviceVendor":"","deviceModel":"","deviceType":"","cpuArchitecture":"amd64","isPrivateMode":"0"},"General":{"language":"en-US","syslang":"","userlang":"","deviceMemory":"","hardwareConcurrency":"8","resolution":"1760x990","colorDepth":"24","screenWidth":"1760","screenHeight":"990","availableHeight":"990","availableResolution":"1760x990","screenAvailableWdth":"1760","timeZone":"-5","timezoneOffset":"300","sessionStorage":"0","cookieEnabled":"1","localStorage":"1","indexedDb":"0","cpuClass":"","openDatabase":"1","navigatorPlatform":"Linux x86_64","vendorWebGL":"1","rendererVideo":"AMD Radeon(TM) Vega 11 Graphics (RAVEN)","software":"","javaEnabled":"0","allSoftware":"","appName":"Netscape","appCodeName":"Mozilla","onLine":"true","opsProfile":"","userProfile":"","screenBufferDepth":"","screendDeviceXDPI":"","screenDeviceYDPI":"","screenLogicalXDPI":"","screenLogicalYPDI":"","screenFontSmoothingEnabled":"","screenUpdateInterval":"","pingIn":"","pingEx":""},"Personalization":{"numberPlugins":"0","numberFonts":"14"},"Alterations":{"adblock":"0","hasLiedLanguages":"0","hasLiedResolution":"0","hasLiedOs":"0","hasLiedBrowser":"0","touchSupport":"0"},"Network":{"publicIp":"","localIp":""},"Site":{"host":"sucursalpersonas.transaccionesbancolombia.com","hostName":"sucursalpersonas.transaccionesbancolombia.com","href":"https://sucursalpersonas.transaccionesbancolombia.com/mua/USER?scis=SyMIhw3R7Bqwr4rJqB%2FUccExc76YT63v9g9%2FKsEB%2FgA%3D","origin":"","pathname":"/mua/USER","port":"","protocol":"https:"},"Identifiers":{"cookie":"354673ff5f6608abc152664feaab0e5b","localStorageValue":"2d4709e2199f2c952bdeab044d7362e0","hash":"131B0C35E2D7EDEF.B592CFEF0A3AFD9D.44"}}`,
			output: `%257B%2522Version%2522%253A%25223%252E5%252E1%255F4%2522%252C%2522Browser%2522%253A%257B%2522userAgent%2522%253A%2522mozilla%252F5%252E0%2520%2528x11%253B%2520linux%2520x86%255F64%253B%2520rv%253A90%252E0%2529%2520gecko%252F20100101%2520firefox%252F90%252E0%2522%252C%2522appVersion%2522%253A%25225%252E0%2520%2528X11%2529%2522%252C%2522platform%2522%253A%2522Linux%2520x86%255F64%2522%252C%2522appMinorVersion%2522%253A%2522%2522%252C%2522cpuClass%2522%253A%2522%2522%252C%2522browserLanguage%2522%253A%2522%2522%252C%2522browserName%2522%253A%2522Firefox%2522%252C%2522browserVersion%2522%253A%252290%252E0%2522%252C%2522browserMajor%2522%253A%252290%2522%252C%2522browserEngineName%2522%253A%2522Gecko%2522%252C%2522browserEngineVersion%2522%253A%252290%252E0%2522%252C%2522osName%2522%253A%2522Linux%2522%252C%2522browserOS%2522%253A%2522Linux%2522%252C%2522osVersion%2522%253A%2522x86%255F64%2522%252C%2522deviceVendor%2522%253A%2522%2522%252C%2522deviceModel%2522%253A%2522%2522%252C%2522deviceType%2522%253A%2522%2522%252C%2522cpuArchitecture%2522%253A%2522amd64%2522%252C%2522isPrivateMode%2522%253A%25220%2522%257D%252C%2522General%2522%253A%257B%2522language%2522%253A%2522en%252DUS%2522%252C%2522syslang%2522%253A%2522%2522%252C%2522userlang%2522%253A%2522%2522%252C%2522deviceMemory%2522%253A%2522%2522%252C%2522hardwareConcurrency%2522%253A%25228%2522%252C%2522resolution%2522%253A%25221760x990%2522%252C%2522colorDepth%2522%253A%252224%2522%252C%2522screenWidth%2522%253A%25221760%2522%252C%2522screenHeight%2522%253A%2522990%2522%252C%2522availableHeight%2522%253A%2522990%2522%252C%2522availableResolution%2522%253A%25221760x990%2522%252C%2522screenAvailableWdth%2522%253A%25221760%2522%252C%2522timeZone%2522%253A%2522%252D5%2522%252C%2522timezoneOffset%2522%253A%2522300%2522%252C%2522sessionStorage%2522%253A%25220%2522%252C%2522cookieEnabled%2522%253A%25221%2522%252C%2522localStorage%2522%253A%25221%2522%252C%2522indexedDb%2522%253A%25220%2522%252C%2522cpuClass%2522%253A%2522%2522%252C%2522openDatabase%2522%253A%25221%2522%252C%2522navigatorPlatform%2522%253A%2522Linux%2520x86%255F64%2522%252C%2522vendorWebGL%2522%253A%25221%2522%252C%2522rendererVideo%2522%253A%2522AMD%2520Radeon%2528TM%2529%2520Vega%252011%2520Graphics%2520%2528RAVEN%2529%2522%252C%2522software%2522%253A%2522%2522%252C%2522javaEnabled%2522%253A%25220%2522%252C%2522allSoftware%2522%253A%2522%2522%252C%2522appName%2522%253A%2522Netscape%2522%252C%2522appCodeName%2522%253A%2522Mozilla%2522%252C%2522onLine%2522%253A%2522true%2522%252C%2522opsProfile%2522%253A%2522%2522%252C%2522userProfile%2522%253A%2522%2522%252C%2522screenBufferDepth%2522%253A%2522%2522%252C%2522screendDeviceXDPI%2522%253A%2522%2522%252C%2522screenDeviceYDPI%2522%253A%2522%2522%252C%2522screenLogicalXDPI%2522%253A%2522%2522%252C%2522screenLogicalYPDI%2522%253A%2522%2522%252C%2522screenFontSmoothingEnabled%2522%253A%2522%2522%252C%2522screenUpdateInterval%2522%253A%2522%2522%252C%2522pingIn%2522%253A%2522%2522%252C%2522pingEx%2522%253A%2522%2522%257D%252C%2522Personalization%2522%253A%257B%2522numberPlugins%2522%253A%25220%2522%252C%2522numberFonts%2522%253A%252214%2522%257D%252C%2522Alterations%2522%253A%257B%2522adblock%2522%253A%25220%2522%252C%2522hasLiedLanguages%2522%253A%25220%2522%252C%2522hasLiedResolution%2522%253A%25220%2522%252C%2522hasLiedOs%2522%253A%25220%2522%252C%2522hasLiedBrowser%2522%253A%25220%2522%252C%2522touchSupport%2522%253A%25220%2522%257D%252C%2522Network%2522%253A%257B%2522publicIp%2522%253A%2522%2522%252C%2522localIp%2522%253A%2522%2522%257D%252C%2522Site%2522%253A%257B%2522host%2522%253A%2522sucursalpersonas%252Etransaccionesbancolombia%252Ecom%2522%252C%2522hostName%2522%253A%2522sucursalpersonas%252Etransaccionesbancolombia%252Ecom%2522%252C%2522href%2522%253A%2522https%253A%252F%252Fsucursalpersonas%252Etransaccionesbancolombia%252Ecom%252Fmua%252FUSER%253Fscis%253DSyMIhw3R7Bqwr4rJqB%25252FUccExc76YT63v9g9%25252FKsEB%25252FgA%25253D%2522%252C%2522origin%2522%253A%2522%2522%252C%2522pathname%2522%253A%2522%252Fmua%252FUSER%2522%252C%2522port%2522%253A%2522%2522%252C%2522protocol%2522%253A%2522https%253A%2522%257D%252C%2522Identifiers%2522%253A%257B%2522cookie%2522%253A%2522354673ff5f6608abc152664feaab0e5b%2522%252C%2522localStorageValue%2522%253A%25222d4709e2199f2c952bdeab044d7362e0%2522%252C%2522hash%2522%253A%2522131B0C35E2D7EDEF%252EB592CFEF0A3AFD9D%252E44%2522%257D%257D`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := encodeDevicePrint(tt.input)
			// assert.Equal(t, tt.output, encoded)

			decoded, err := decodeDevicePrint(encoded)
			require.NoError(t, err)
			assert.Equal(t, tt.input, decoded)
		})
	}
}

func TestDecodeDevicePrint(t *testing.T) {
	t.SkipNow()

	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "base case",
			input:  `%257B%2522Version%2522%253A%25223%252E5%252E1%255F4%2522%252C%2522Browser%2522%253A%257B%2522userAgent%2522%253A%2522mozilla%252F5%252E0%2520%2528x11%253B%2520linux%2520x86%255F64%253B%2520rv%253A90%252E0%2529%2520gecko%252F20100101%2520firefox%252F90%252E0%2522%252C%2522appVersion%2522%253A%25225%252E0%2520%2528X11%2529%2522%252C%2522platform%2522%253A%2522Linux%2520x86%255F64%2522%252C%2522appMinorVersion%2522%253A%2522%2522%252C%2522cpuClass%2522%253A%2522%2522%252C%2522browserLanguage%2522%253A%2522%2522%252C%2522browserName%2522%253A%2522Firefox%2522%252C%2522browserVersion%2522%253A%252290%252E0%2522%252C%2522browserMajor%2522%253A%252290%2522%252C%2522browserEngineName%2522%253A%2522Gecko%2522%252C%2522browserEngineVersion%2522%253A%252290%252E0%2522%252C%2522osName%2522%253A%2522Linux%2522%252C%2522browserOS%2522%253A%2522Linux%2522%252C%2522osVersion%2522%253A%2522x86%255F64%2522%252C%2522deviceVendor%2522%253A%2522%2522%252C%2522deviceModel%2522%253A%2522%2522%252C%2522deviceType%2522%253A%2522%2522%252C%2522cpuArchitecture%2522%253A%2522amd64%2522%252C%2522isPrivateMode%2522%253A%25220%2522%257D%252C%2522General%2522%253A%257B%2522language%2522%253A%2522en%252DUS%2522%252C%2522syslang%2522%253A%2522%2522%252C%2522userlang%2522%253A%2522%2522%252C%2522deviceMemory%2522%253A%2522%2522%252C%2522hardwareConcurrency%2522%253A%25228%2522%252C%2522resolution%2522%253A%25221760x990%2522%252C%2522colorDepth%2522%253A%252224%2522%252C%2522screenWidth%2522%253A%25221760%2522%252C%2522screenHeight%2522%253A%2522990%2522%252C%2522availableHeight%2522%253A%2522990%2522%252C%2522availableResolution%2522%253A%25221760x990%2522%252C%2522screenAvailableWdth%2522%253A%25221760%2522%252C%2522timeZone%2522%253A%2522%252D5%2522%252C%2522timezoneOffset%2522%253A%2522300%2522%252C%2522sessionStorage%2522%253A%25220%2522%252C%2522cookieEnabled%2522%253A%25221%2522%252C%2522localStorage%2522%253A%25221%2522%252C%2522indexedDb%2522%253A%25220%2522%252C%2522cpuClass%2522%253A%2522%2522%252C%2522openDatabase%2522%253A%25221%2522%252C%2522navigatorPlatform%2522%253A%2522Linux%2520x86%255F64%2522%252C%2522vendorWebGL%2522%253A%25221%2522%252C%2522rendererVideo%2522%253A%2522AMD%2520Radeon%2528TM%2529%2520Vega%252011%2520Graphics%2520%2528RAVEN%2529%2522%252C%2522software%2522%253A%2522%2522%252C%2522javaEnabled%2522%253A%25220%2522%252C%2522allSoftware%2522%253A%2522%2522%252C%2522appName%2522%253A%2522Netscape%2522%252C%2522appCodeName%2522%253A%2522Mozilla%2522%252C%2522onLine%2522%253A%2522true%2522%252C%2522opsProfile%2522%253A%2522%2522%252C%2522userProfile%2522%253A%2522%2522%252C%2522screenBufferDepth%2522%253A%2522%2522%252C%2522screendDeviceXDPI%2522%253A%2522%2522%252C%2522screenDeviceYDPI%2522%253A%2522%2522%252C%2522screenLogicalXDPI%2522%253A%2522%2522%252C%2522screenLogicalYPDI%2522%253A%2522%2522%252C%2522screenFontSmoothingEnabled%2522%253A%2522%2522%252C%2522screenUpdateInterval%2522%253A%2522%2522%252C%2522pingIn%2522%253A%2522%2522%252C%2522pingEx%2522%253A%2522%2522%257D%252C%2522Personalization%2522%253A%257B%2522numberPlugins%2522%253A%25220%2522%252C%2522numberFonts%2522%253A%252214%2522%257D%252C%2522Alterations%2522%253A%257B%2522adblock%2522%253A%25220%2522%252C%2522hasLiedLanguages%2522%253A%25220%2522%252C%2522hasLiedResolution%2522%253A%25220%2522%252C%2522hasLiedOs%2522%253A%25220%2522%252C%2522hasLiedBrowser%2522%253A%25220%2522%252C%2522touchSupport%2522%253A%25220%2522%257D%252C%2522Network%2522%253A%257B%2522publicIp%2522%253A%2522%2522%252C%2522localIp%2522%253A%2522%2522%257D%252C%2522Site%2522%253A%257B%2522host%2522%253A%2522sucursalpersonas%252Etransaccionesbancolombia%252Ecom%2522%252C%2522hostName%2522%253A%2522sucursalpersonas%252Etransaccionesbancolombia%252Ecom%2522%252C%2522href%2522%253A%2522https%253A%252F%252Fsucursalpersonas%252Etransaccionesbancolombia%252Ecom%252Fmua%252FUSER%253Fscis%253DSyMIhw3R7Bqwr4rJqB%25252FUccExc76YT63v9g9%25252FKsEB%25252FgA%25253D%2522%252C%2522origin%2522%253A%2522%2522%252C%2522pathname%2522%253A%2522%252Fmua%252FUSER%2522%252C%2522port%2522%253A%2522%2522%252C%2522protocol%2522%253A%2522https%253A%2522%257D%252C%2522Identifiers%2522%253A%257B%2522cookie%2522%253A%2522354673ff5f6608abc152664feaab0e5b%2522%252C%2522localStorageValue%2522%253A%25222d4709e2199f2c952bdeab044d7362e0%2522%252C%2522hash%2522%253A%2522131B0C35E2D7EDEF%252EB592CFEF0A3AFD9D%252E44%2522%257D%257D`,
			output: `{"Version":"3.5.1_4","Browser":{"userAgent":"mozilla/5.0 (x11; linux x86_64; rv:90.0) gecko/20100101 firefox/90.0","appVersion":"5.0 (X11)","platform":"Linux x86_64","appMinorVersion":"","cpuClass":"","browserLanguage":"","browserName":"Firefox","browserVersion":"90.0","browserMajor":"90","browserEngineName":"Gecko","browserEngineVersion":"90.0","osName":"Linux","browserOS":"Linux","osVersion":"x86_64","deviceVendor":"","deviceModel":"","deviceType":"","cpuArchitecture":"amd64","isPrivateMode":"0"},"General":{"language":"en-US","syslang":"","userlang":"","deviceMemory":"","hardwareConcurrency":"8","resolution":"1760x990","colorDepth":"24","screenWidth":"1760","screenHeight":"990","availableHeight":"990","availableResolution":"1760x990","screenAvailableWdth":"1760","timeZone":"-5","timezoneOffset":"300","sessionStorage":"0","cookieEnabled":"1","localStorage":"1","indexedDb":"0","cpuClass":"","openDatabase":"1","navigatorPlatform":"Linux x86_64","vendorWebGL":"1","rendererVideo":"AMD Radeon(TM) Vega 11 Graphics (RAVEN)","software":"","javaEnabled":"0","allSoftware":"","appName":"Netscape","appCodeName":"Mozilla","onLine":"true","opsProfile":"","userProfile":"","screenBufferDepth":"","screendDeviceXDPI":"","screenDeviceYDPI":"","screenLogicalXDPI":"","screenLogicalYPDI":"","screenFontSmoothingEnabled":"","screenUpdateInterval":"","pingIn":"","pingEx":""},"Personalization":{"numberPlugins":"0","numberFonts":"14"},"Alterations":{"adblock":"0","hasLiedLanguages":"0","hasLiedResolution":"0","hasLiedOs":"0","hasLiedBrowser":"0","touchSupport":"0"},"Network":{"publicIp":"","localIp":""},"Site":{"host":"sucursalpersonas.transaccionesbancolombia.com","hostName":"sucursalpersonas.transaccionesbancolombia.com","href":"https://sucursalpersonas.transaccionesbancolombia.com/mua/USER?scis=SyMIhw3R7Bqwr4rJqB%2FUccExc76YT63v9g9%2FKsEB%2FgA%3D","origin":"","pathname":"/mua/USER","port":"","protocol":"https:"},"Identifiers":{"cookie":"354673ff5f6608abc152664feaab0e5b","localStorageValue":"2d4709e2199f2c952bdeab044d7362e0","hash":"131B0C35E2D7EDEF.B592CFEF0A3AFD9D.44"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := decodeDevicePrint(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.output, decoded)

			encoded := encodeDevicePrint(decoded)
			assert.Equal(t, tt.input, encoded)
		})
	}
}

// 1EDE161732CBEFA
// 123.456.78.9

func TestMapPassword(t *testing.T) {
	actual := mapPassword(map[string]string{
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
	}, "1234")
	assert.Equal(t, "UMjS974CeHd35bG8", actual)
}
