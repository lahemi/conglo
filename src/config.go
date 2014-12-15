// TODO: a proper external config
package main

import "os"

func getConfig() map[string]map[string]string {
	return map[string]map[string]string{
		"filecan": map[string]string{
			"base":       "/filecan",
			"uploadPath": basePath + "filecanUploads/",
			"flimit":     "327680",
			"uploadHTML": htmlPath + "filecan.html",
		},
		"fserve": map[string]string{
			"dir": os.Getenv("HOME") + "/FSERVE/",
		},
		"pastecan": map[string]string{
			"base":       "/pastecan",
			"pastePath":  basePath + "pastes/",
			"pasteIndex": htmlPath + "pastecan.html",
			"viewHTML":   htmlPath + "pasteview.html",
		},
		"datacan": map[string]string{
			"base":         "/datacan",
			"DBfile":       basePath + "musicks.db",
			"datacanIndex": htmlPath + "datacan.html",
		},
	}
}
